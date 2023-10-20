package imessages

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/imroc/req/v3"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) R(ctx context.Context) *req.Request {
	return common.HttpClient(x.V.EmqxHost).
		SetCommonBasicAuth(x.V.EmqxApiKey, x.V.EmqxSecretKey).
		R().SetContext(ctx)
}

func (x *Service) GetNodes(ctx context.Context) (r interface{}, err error) {
	if _, err = x.R(ctx).
		SetSuccessResult(&r).
		SetErrorResult(&r).
		Get("nodes"); err != nil {
		return
	}
	return
}

func (x *Service) CreateRule(ctx context.Context, id primitive.ObjectID) (r M, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}

	var e string
	if _, err = x.R(ctx).
		SetBody(M{
			"name": data.Topic,
			"sql":  fmt.Sprintf(`SELECT * FROM "%s/#"`, data.Topic),
			"actions": []interface{}{
				"webhook:logset",
			},
			"enable":      true,
			"description": data.Description,
		}).
		SetSuccessResult(&r).
		SetErrorResult(&e).
		Post("rules"); err != nil {
		return
	}
	if e != "" {
		return nil, help.E("Imessage.CreateRule", e)
	}

	if _, err = x.Db.Collection("imessages").
		UpdateOne(ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"rule": r["id"]}},
		); err != nil {
		return
	}

	return
}

func (x *Service) UpdateRule(ctx context.Context, id primitive.ObjectID) (r M, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}

	if _, err = x.R(ctx).
		SetBody(M{
			"name": data.Topic,
			"sql":  fmt.Sprintf(`SELECT * FROM "%s/#"`, data.Topic),
			"actions": []interface{}{
				"webhook:logset",
			},
			"enable":      true,
			"description": data.Description,
		}).
		SetSuccessResult(&r).
		SetErrorResult(&r).
		Put(fmt.Sprintf(`rules/%s`, data.Rule)); err != nil {
		return
	}

	return
}

func (x *Service) DeleteRule(ctx context.Context, id primitive.ObjectID) (err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}

	if _, err = x.R(ctx).
		Delete(fmt.Sprintf(`rules/%s`, data.Rule)); err != nil {
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
		if _, err = x.R(ctx).
			SetSuccessResult(&r).
			SetErrorResult(&r).
			Get(fmt.Sprintf("mqtt/topic_metrics/%s%%2f%s", data.Topic, pid.Hex())); err != nil {
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
		if _, err = x.R(ctx).
			SetBody(M{"topic": fmt.Sprintf(`%s/%s`, data.Topic, pid.Hex())}).
			SetSuccessResult(&r).
			SetErrorResult(&r).
			Post("mqtt/topic_metrics"); err != nil {
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
		var r M
		if _, err = x.R(ctx).
			SetSuccessResult(&r).
			SetErrorResult(&r).
			Delete(fmt.Sprintf("mqtt/topic_metrics/%s%%2f%s", data.Topic, pid.Hex())); err != nil {
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
	if _, err = x.R(ctx).
		SetBody(M{
			"topic":   dto.Topic,
			"payload": payload,
		}).
		SetSuccessResult(&r).
		SetErrorResult(&r).
		Post("publish"); err != nil {
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
			if _, err = x.CreateRule(ctx, id); err != nil {
				hlog.Error(err)
			}
			if _, err = x.CreateMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		case "update_by_id":
			id, _ := primitive.ObjectIDFromHex(dto.Id)
			if _, err = x.UpdateRule(ctx, id); err != nil {
				hlog.Error(err)
			}
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
