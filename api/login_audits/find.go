package login_audits

import (
	"context"
	"litho-api/common"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FindDto struct {
	common.FindDto

	Start string `query:"start" vd:"required,datetime=2006-01-02"`
	End   string `query:"end" vd:"required,datetime=2006-01-02"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	total, data, err := x.LoginAuditsX.Find(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("x-total", strconv.FormatInt(total, 10))
	c.JSON(200, data)
}

type FindResult struct {
	ID        bson.ObjectID   `json:"id" bson:"_id"`
	Ts        time.Time       `json:"ts" bson:"ts"`
	Meta      *LoginAuditMeta `json:"meta" bson:"meta"`
	Name      string          `json:"name" bson:"-"`
	Phone     string          `json:"phone" bson:"-"`
	OrgIDs    []string        `json:"org_ids" bson:"-"`
	Orgs      []string        `json:"orgs" bson:"-"`
	IP        string          `json:"ip" bson:"ip"`
	PrivateIP bool            `json:"private_ip" bson:"private_ip"`
	Country   string          `json:"country" bson:"country"`
	Region    string          `json:"region" bson:"region"`
	City      string          `json:"city" bson:"city"`
	Org       string          `json:"org" bson:"org"`
	Agent     string          `json:"agent" bson:"agent"`
}

type LoginAuditMeta struct {
	Plat   string `json:"plat" bson:"plat"`
	Mode   int    `json:"mode" bson:"mode"`
	UID    int64  `json:"-" bson:"uid"`
	UserID string `json:"user_id" bson:"-"`
}

func (x *Service) Find(ctx context.Context, dto FindDto) (int64, []*FindResult, error) {
	filter := bson.M{
		"meta.plat": "system",
		"ts": bson.M{
			"$gte": common.StartDate(dto.Start),
			"$lt":  common.EndDate(dto.End),
		},
	}
	if value, err := strconv.ParseInt(dto.Q, 10, 64); err == nil && value != 0 {
		filter["meta.uid"] = value
	}
	collection := x.MDb.Collection("login_audits")
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}
	pageSize := int64(dto.GetPageSize())
	option := options.Find().
		SetSort(bson.D{{Key: "ts", Value: -1}}).
		SetLimit(pageSize).
		SetSkip(pageSize * dto.Page)
	cursor, err := collection.Find(ctx, filter, option)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	results := make([]*FindResult, 0)
	userIDs := make([]string, 0)
	seen := make(map[string]struct{})
	for cursor.Next(ctx) {
		var data FindResult
		if err = cursor.Decode(&data); err != nil {
			return 0, nil, err
		}
		if data.Meta == nil {
			data.Meta = &LoginAuditMeta{}
		}
		data.OrgIDs = make([]string, 0)
		data.Orgs = make([]string, 0)
		data.Meta.UserID = strconv.FormatInt(data.Meta.UID, 10)
		results = append(results, &data)
		if _, ok := seen[data.Meta.UserID]; !ok {
			seen[data.Meta.UserID] = struct{}{}
			userIDs = append(userIDs, data.Meta.UserID)
		}
	}
	if err = cursor.Err(); err != nil {
		return 0, nil, err
	}
	users, err := x.GetUserM(ctx, userIDs)
	if err != nil {
		return 0, nil, err
	}
	for _, result := range results {
		if user, ok := users[result.Meta.UserID]; ok {
			result.Name = user.Name
			result.Phone = user.Phone
			result.OrgIDs = user.OrgIDs
			result.Orgs = user.Orgs
		}
	}
	return total, results, nil
}
