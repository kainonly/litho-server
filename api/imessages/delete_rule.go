package imessages

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteRuleDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) DeleteRule(ctx context.Context, c *app.RequestContext) {
	var dto DeleteRuleDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	if err := x.ImessagesX.DeleteRule(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
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
