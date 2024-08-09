package queues

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/api/projects"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"sync"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	QueuesX *Service
}

type Service struct {
	*common.Inject
	ProjectsX *projects.Service
}

type M = map[string]interface{}

var clients = sync.Map{}

func (x *Service) GetClient(projectId primitive.ObjectID) (client *nats.Conn, err error) {
	if i, ok := clients.Load(projectId.Hex()); ok {
		client = i.(*nats.Conn)
		return
	}
	var project model.Project
	if project, err = x.ProjectsX.Get(context.TODO(), projectId); err != nil {
		return
	}
	var seed []byte
	if seed, err = x.Cipher.Decode(project.Nats.Seed); err != nil {
		return
	}
	var kp nkeys.KeyPair
	if kp, err = nkeys.FromSeed(seed); err != nil {
		return
	}
	defer kp.Wipe()
	var pub string
	if pub, err = kp.PublicKey(); err != nil {
		return
	}
	if !nkeys.IsValidPublicUserKey(pub) {
		return nil, fmt.Errorf("nkey fail")
	}
	if client, err = nats.Connect(
		strings.Join(x.V.Nats.Hosts, ","),
		nats.MaxReconnects(-1),
		nats.Nkey(pub, func(nonce []byte) ([]byte, error) {
			sig, _ := kp.Sign(nonce)
			return sig, nil
		}),
	); err != nil {
		return
	}
	clients.Store(projectId.Hex(), client)
	return
}

func (x *Service) GetJetStream(ctx context.Context, projectId primitive.ObjectID) (js nats.JetStreamContext, err error) {
	var nc *nats.Conn
	if nc, err = x.GetClient(projectId); err != nil {
		return
	}
	if js, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
		nats.Context(ctx),
	); err != nil {
		return
	}
	return
}

func (x *Service) Destroy(ctx context.Context, js nats.JetStreamContext, name string) (err error) {
	if _, err = js.StreamInfo(name, nats.Context(ctx)); err != nil {
		if err != nats.ErrStreamNotFound {
			return
		} else {
			return nil
		}
	}
	if err = js.DeleteStream(name, nats.Context(ctx)); err != nil {
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
		case rest.ActionDelete:
			projectId, _ := primitive.ObjectIDFromHex(dto.Data.(M)["project"].(string))
			var js nats.JetStreamContext
			if js, err = x.GetJetStream(ctx, projectId); err != nil {
				hlog.Error(err)
			}
			if err = x.Destroy(ctx, js, dto.Data.(M)["_id"].(string)); err != nil {
				hlog.Error(err)
			}
			break
		case rest.ActionBulkDelete:
			data := dto.Data.([]interface{})
			projectId, _ := primitive.ObjectIDFromHex(data[0].(M)["project"].(string))
			var js nats.JetStreamContext
			if js, err = x.GetJetStream(ctx, projectId); err != nil {
				hlog.Error(err)
			}
			for _, v := range data {
				if err = x.Destroy(ctx, js, v.(M)["_id"].(string)); err != nil {
					hlog.Error(err)
				}
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
