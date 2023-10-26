package queues

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) Sync(ctx context.Context, id primitive.ObjectID) (err error) {
	var data model.Queue
	if err = x.Db.Collection("queues").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&data); err != nil {
		return
	}
	if _, err = x.JetStream.StreamInfo(data.Name); err != nil {
		if err != nats.ErrStreamNotFound {
			return
		}
	}
	cfg := &nats.StreamConfig{
		Name:      data.Name,
		Subjects:  data.Subjects,
		MaxMsgs:   data.MaxMsgs,
		MaxBytes:  data.MaxBytes,
		MaxAge:    data.MaxAge,
		Retention: nats.WorkQueuePolicy,
	}
	if data.Description != "" {
		cfg.Description = data.Description
	}
	if err == nats.ErrStreamNotFound {
		if _, err = x.JetStream.AddStream(cfg, nats.Context(ctx)); err != nil {
			return
		}
	} else {
		if _, err = x.JetStream.UpdateStream(cfg, nats.Context(ctx)); err != nil {
			return
		}
	}
	return
}

func (x *Service) Destroy(ctx context.Context, ids []primitive.ObjectID) (err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("queues").Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	}); err != nil {
		return
	}
	for cursor.Next(ctx) {
		var data model.Queue
		if err = cursor.Decode(&data); err != nil {
			return
		}

		if _, err = x.JetStream.StreamInfo(data.Name); err != nil {
			if err != nats.ErrStreamNotFound {
				return
			} else {
				return nil
			}
		}
		if err = x.JetStream.DeleteStream(data.Name); err != nil {
			return
		}
	}
	return
}

func (x *Service) Info(ctx context.Context, id primitive.ObjectID) (r *nats.StreamInfo, err error) {
	var data model.Queue
	if err = x.Db.Collection("queues").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	if r, err = x.JetStream.StreamInfo(data.Name, nats.Context(ctx)); err != nil {
		return
	}
	r.Cluster = nil
	return
}

func (x *Service) Publish(ctx context.Context, dto PublishDto) (r interface{}, err error) {
	var payload []byte
	if payload, err = sonic.Marshal(dto.Payload); err != nil {
		return
	}
	if r, err = x.JetStream.Publish(dto.Subject, payload, nats.Context(ctx)); err != nil {
		return
	}
	return
}

func (x *Service) Event() (err error) {
	if _, err = x.JetStream.QueueSubscribe(`events.queues`, `EVENT_queues`, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}
		switch dto.Action {
		case rest.ActionCreate:
			id, _ := primitive.ObjectIDFromHex(dto.Result.(M)["InsertedID"].(string))
			if err = x.Sync(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		case rest.ActionUpdateById:
			id, _ := primitive.ObjectIDFromHex(dto.Id)
			if err = x.Sync(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
