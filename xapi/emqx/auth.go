package emqx

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthDto struct {
	Identity string `json:"identity" vd:"required"`
	Token    string `json:"token" vd:"required"`
}

func (x *Controller) Auth(ctx context.Context, c *app.RequestContext) {
	var dto AuthDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EmqxService.Auth(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
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
