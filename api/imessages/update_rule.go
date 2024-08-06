package imessages

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateRuleDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) UpdateRule(ctx context.Context, c *app.RequestContext) {
	var dto UpdateRuleDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesX.UpdateRule(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

func (x *Service) UpdateRule(ctx context.Context, id primitive.ObjectID) (r M, err error) {
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}

	if data.Rule != "" {
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
