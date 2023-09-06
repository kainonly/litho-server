package imessages

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) Send(ctx context.Context, method string, path string, data interface{}) (result interface{}, err error) {
	var req *http.Request
	var body io.Reader
	if data != nil {
		var b []byte
		if b, err = sonic.Marshal(data); err != nil {
			return
		}
		body = bytes.NewBuffer(b)
	}
	url := fmt.Sprintf("%s/%s", x.V.EmqxHost, path)
	if req, err = http.NewRequest(method, url, body); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(x.V.EmqxApiKey, x.V.EmqxSecretKey)
	req.WithContext(ctx)

	client := &http.Client{Timeout: time.Second * 5}
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	if err = decoder.
		NewStreamDecoder(resp.Body).
		Decode(&result); err != nil {
		return
	}
	return
}

func (x *Service) GetNodes(ctx context.Context) (r interface{}, err error) {
	if r, err = x.Send(ctx, "GET", "nodes", nil); err != nil {
		return
	}
	return
}

func (x *Service) GetMetrics(ctx context.Context, id primitive.ObjectID) (rs []interface{}, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}

	for _, pid := range data.Projects {
		var r interface{}
		if r, err = x.Send(ctx, "GET",
			fmt.Sprintf("mqtt/topic_metrics/%s%%2f%s", data.Topic, pid.Hex()), nil); err != nil {
			return
		}
		rs = append(rs, r)
	}

	return
}

func (x *Service) CreateMetrics(ctx context.Context, id primitive.ObjectID) (rs []interface{}, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		var r interface{}
		if r, err = x.Send(ctx, "POST",
			"mqtt/topic_metrics", M{"topic": fmt.Sprintf(`%s/%s`, data.Topic, pid.Hex())}); err != nil {
			return
		}
		rs = append(rs, r)
	}
	return
}

func (x *Service) DeleteMetrics(ctx context.Context, id primitive.ObjectID) (rs []interface{}, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		var r interface{}
		if r, err = x.Send(ctx, "DELETE",
			fmt.Sprintf("mqtt/topic_metrics/%s%%2f%s", data.Topic, pid.Hex()), nil); err != nil {
			return
		}
		rs = append(rs, r)
	}
	return
}

func (x *Service) Publish(ctx context.Context, dto PublishDto) (r interface{}, err error) {
	var payload string
	if payload, err = sonic.MarshalString(dto.Payload); err != nil {
		return
	}
	if r, err = x.Send(ctx, "POST",
		"publish", M{
			"topic":   dto.Topic,
			"payload": payload,
		}); err != nil {
		return
	}
	return
}

func (x *Service) Event() (err error) {
	subj := x.V.NameX(".", "events", "imessages")
	queue := x.V.Name("events", "imessages")
	if _, err = x.JetStream.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}
		switch dto.Action {
		case "create":
			id, _ := primitive.ObjectIDFromHex(dto.Result.(M)["InsertedID"].(string))
			if _, err = x.CreateMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		case "update_by_id":
			id, _ := primitive.ObjectIDFromHex(dto.Id)
			if _, err = x.DeleteMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			if _, err = x.CreateMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
