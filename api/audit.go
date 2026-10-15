package api

import (
	"context"
	"litho-api/common"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuditDto struct {
	User        *common.IAMUser `bson:"-"`
	BodyRaw     []byte          `bson:"-"`
	ResponseRaw []byte          `bson:"-"`

	Ts       time.Time `bson:"ts"`
	Meta     AuditMeta `bson:"meta"`
	Path     string    `bson:"path"`
	Body     bson.M    `bson:"body"`
	Code     int       `bson:"code"`
	Response bson.M    `bson:"response"`
	Ip       string    `bson:"ip"`
	Agent    string    `bson:"agent"`
}

type AuditMeta struct {
	PLAT string `bson:"plat"`
	RES  string `bson:"res"`
	ACT  string `bson:"act"`
	RID  int    `bson:"rid"`
	UID  int    `bson:"uid"`
}

func (x *API) Audit() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if string(c.Method()) != "POST" {
			c.Next(ctx)
			return
		}
		c.Next(ctx)
		tracking := common.GetTracking(c)
		// 不存在审计跟踪
		if tracking == nil {
			return
		}
		user := common.GetIAM(c)
		now := time.Now()
		for _, id := range tracking.RIDS {
			rid, err := strconv.Atoi(id)
			if err != nil {
				continue
			}
			go x.AuditPush(AuditDto{
				User: user,
				Ts:   now,
				Meta: AuditMeta{
					PLAT: `system`,
					RES:  tracking.RES,
					ACT:  tracking.ACT,
					RID:  rid,
				},
				Path:        string(c.Path()),
				Ip:          string(c.GetHeader(x.V.Ip)),
				BodyRaw:     c.Request.Body(),
				Code:        c.Response.StatusCode(),
				ResponseRaw: c.Response.Body(),
				Agent:       string(c.GetHeader(`User-Agent`)),
			})
		}

	}
}

func (x *API) AuditPush(dto AuditDto) (err error) {
	if dto.Meta.UID, err = strconv.Atoi(dto.User.ID); err != nil {
		return
	}
	if err = sonic.Unmarshal(dto.BodyRaw, &dto.Body); err != nil {
		return
	}
	switch dto.Meta.RES {
	case "users":
		// 密码脱敏
		if _, ok := dto.Body["password"]; ok {
			dto.Body["password"] = "*"
		}
		if _, ok := dto.Body["old"]; ok {
			dto.Body["old"] = "*"
		}
		if _, ok := dto.Body["code"]; ok {
			dto.Body["code"] = "*"
		}
		break
	}
	if err = sonic.Unmarshal(dto.ResponseRaw, &dto.Response); err != nil {
		return
	}
	return x.Transfer.Send(`audits`, dto)
}
