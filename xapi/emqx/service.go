package emqx

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	transfer "github.com/weplanx/collector/client"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) Auth(ctx context.Context, dto AuthDto) (err error) {
	var data model.Project
	if err = x.Db.Collection("projects").
		FindOne(ctx, bson.M{"_id": dto.Identity}).
		Decode(&data); err != nil {
		return
	}
	p := passport.New(
		passport.SetIssuer(x.V.Namespace),
		passport.SetKey(fmt.Sprintf(`%s:%s`, data.SecretId, data.SecretKey)),
	)
	if _, err = p.Verify(dto.Token); err != nil {
		return
	}
	return
}

func (x *Service) Acl(ctx context.Context, dto AclDto) (err error) {
	identity, deny := dto.Identity.Hex(), true
	topic := strings.Split(dto.Topic, "/")
	msg := fmt.Sprintf(`The user [%s] is not authorized for this topic [%s]`,
		identity, dto.Topic)
	if !(len(topic) >= 2 && topic[1] == identity) {
		return errors.NewPublic(msg)
	}
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"topic": topic[0]}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		if pid.Hex() == identity {
			deny = false
			break
		}
	}
	if deny {
		return errors.NewPublic(msg)
	}
	return
}

func (x *Service) Bridge(ctx context.Context, dto BridgeDto) (err error) {
	return x.Transfer.Publish(ctx, "logset_imessages", transfer.Payload{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"metadata": map[string]interface{}{
				"client": dto.Client,
				"topic":  dto.Topic,
			},
			"payload": dto.Payload,
		},
	})
}
