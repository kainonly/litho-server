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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) Auth(ctx context.Context, dto AuthDto) (err error) {
	var data model.Project
	id, _ := primitive.ObjectIDFromHex(dto.Identity)
	if err = x.Db.Collection("projects").
		FindOne(ctx, bson.M{"_id": id}).
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
	deny := true
	topic := strings.Split(dto.Topic, "/")
	msg := fmt.Sprintf(`The user [%s] is not authorized for this topic [%s]`,
		dto.Identity, dto.Topic)
	if !(len(topic) >= 2 && topic[1] == dto.Identity) {
		return errors.NewPublic(msg)
	}
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"topic": topic[0]}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		if pid.Hex() == dto.Identity {
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
