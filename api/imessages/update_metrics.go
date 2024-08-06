package imessages

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateMetricsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) UpdateMetrics(ctx context.Context, c *app.RequestContext) {
	var dto UpdateMetricsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesX.UpdateMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

func (x *Service) UpdateMetrics(ctx context.Context, id primitive.ObjectID) (rs []interface{}, err error) {
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
