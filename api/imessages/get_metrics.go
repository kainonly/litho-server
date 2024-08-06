package imessages

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetMetricsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetMetrics(ctx context.Context, c *app.RequestContext) {
	var dto GetMetricsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesX.GetMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
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
