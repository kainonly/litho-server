package index

import (
	"context"
	"litho-api/common"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GetHistoryDto struct {
	Page int64 `header:"x-page" vd:"omitempty,min=0"`
}

func (x *Controller) GetHistory(ctx context.Context, c *app.RequestContext) {
	var dto GetHistoryDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, result, err := x.IndexX.GetHistory(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, result)
}

type GetHistoryResult struct {
	ID        bson.ObjectID `json:"id" bson:"_id" `
	Ts        time.Time     `json:"ts"`
	Ip        string        `json:"ip"`
	PrivateIp bool          `json:"private_ip" bson:"private_ip"`
	Country   string        `json:"country"`
	Region    string        `json:"region"`
	City      string        `json:"city"`
	Org       string        `json:"org"`
	Agent     string        `json:"agent"`
}

func (x *Service) GetHistory(ctx context.Context, user *common.IAMUser, dto GetHistoryDto) (count int64, results []GetHistoryResult, err error) {
	uid, _ := strconv.Atoi(user.ID)
	filter := bson.M{
		`meta.plat`: `system`,
		`meta.uid`:  uid,
	}
	if count, err = x.MDb.Collection(`login_audits`).
		CountDocuments(ctx, filter); err != nil {
		return
	}
	var cursor *mongo.Cursor
	opt := options.Find().
		SetSort(bson.D{{Key: "ts", Value: -1}}).
		SetLimit(20).
		SetSkip(dto.Page * 20)
	if cursor, err = x.MDb.Collection(`login_audits`).
		Find(ctx, filter, opt); err != nil {
		return
	}
	defer cursor.Close(ctx)
	results = make([]GetHistoryResult, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return
	}
	return
}
