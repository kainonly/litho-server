package imessages

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/common/errors"
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

type HttpOption struct {
	Method string
	Path   string
	Data   interface{}
	Result interface{}
}

func (x *Service) Send(ctx context.Context, option HttpOption) (err error) {
	var req *http.Request
	var body io.Reader
	if option.Data != nil {
		var b []byte
		if b, err = sonic.Marshal(option.Data); err != nil {
			return
		}
		body = bytes.NewBuffer(b)
	}
	if req, err = http.NewRequest(option.Method, x.V.EmqxHost+option.Path, body); err != nil {
		return
	}
	req.SetBasicAuth(x.V.EmqxApiKey, x.V.EmqxSecretKey)
	req.WithContext(ctx)

	client := &http.Client{Timeout: time.Second * 5}
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	ok := map[int]bool{200: true, 201: true, 204: true}
	if !ok[resp.StatusCode] {
		err = errors.NewPublic(resp.Status)
		return
	}
	return decoder.NewStreamDecoder(resp.Body).Decode(&option.Result)
}

func (x *Service) GetNodes(ctx context.Context) (r []M, err error) {
	if err = x.Send(ctx, HttpOption{
		Method: "GET",
		Path:   "/nodes",
		Result: &r,
	}); err != nil {
		return
	}
	return
}

func (x *Service) GetMetrics(ctx context.Context, id primitive.ObjectID) (result []M, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}

	for _, pid := range data.Projects {
		var v M
		if err = x.Send(ctx, HttpOption{
			Method: "GET",
			Path:   fmt.Sprintf(`topic_metrics/%s/%s`, data.Topic, pid.Hex()),
			Result: &v,
		}); err != nil {
			return
		}
		result = append(result, v)
	}

	return
}

func (x *Service) CreateMetrics(ctx context.Context, id primitive.ObjectID) (result []M, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		var v M
		if err = x.Send(ctx, HttpOption{
			Method: "POST",
			Path:   fmt.Sprintf(`topic_metrics/%s/%s`, data.Topic, pid.Hex()),
			Result: &v,
		}); err != nil {
			return
		}
		result = append(result, v)
	}
	return
}

func (x *Service) DeleteMetrics(ctx context.Context, id primitive.ObjectID) (result []M, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		var v M
		if err = x.Send(ctx, HttpOption{
			Method: "DELETE",
			Path:   fmt.Sprintf(`topic_metrics/%s/%s`, data.Topic, pid.Hex()),
			Result: &v,
		}); err != nil {
			return
		}
		result = append(result, v)
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
