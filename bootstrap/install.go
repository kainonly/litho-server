package bootstrap

import (
	"api/common"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/qri-io/jsonschema"
	"github.com/weplanx/go/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type Install struct {
	Db       *mongo.Database
	Username string
	Password string
	Template string
}

func (x *Install) Basic(ctx context.Context) (err error) {
	var exists []string
	if exists, err = x.Db.ListCollectionNames(ctx, bson.M{
		"name": bson.M{"$in": bson.A{"roles", "users"}},
	}); err != nil {
		return
	}
	if len(exists) != 0 {
		return errors.New("操作不被允许, [roles] 与 [users] 集合是存在的")
	}
	// 初始化权限组
	var roles *mongo.InsertOneResult
	if roles, err = x.Db.Collection("roles").
		InsertOne(ctx,
			common.NewRole("超级管理员").
				SetDescription("系统默认设置").
				SetLabel("默认"),
		); err != nil {
		return
	}
	if _, err = x.Db.Collection("roles").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"name": 1},
				Options: options.Index().SetName("uk_name").SetUnique(true),
			},
		},
	); err != nil {
		return
	}

	// 初始化管理用户
	passwordHash, _ := helper.PasswordHash(x.Password)
	user := common.NewUser(x.Username, passwordHash).
		SetRoles([]primitive.ObjectID{roles.InsertedID.(primitive.ObjectID)})
	if _, err = x.Db.Collection("users").
		InsertOne(ctx, user); err != nil {
		return
	}
	if _, err = x.Db.Collection("users").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"username": 1},
				Options: options.Index().SetName("uk_username").SetUnique(true),
			},
		},
	); err != nil {
		return
	}

	// 初始化日志时序集合
	if err = x.Db.CreateCollection(ctx, "login_logs",
		options.CreateCollection().
			SetTimeSeriesOptions(options.TimeSeries().SetTimeField("time")),
	); err != nil {
		return
	}
	return
}

type TemplateDto struct {
	Schema   string    `json:"$schema"`
	Contents []Content `json:"contents"`
}

type Content struct {
	Parent     *primitive.ObjectID `bson:"parent" json:"-"`
	Name       string              `bson:"name" json:"name"`
	Icon       string              `bson:"icon,omitempty" json:"icon,omitempty"`
	Kind       string              `bson:"kind" json:"kind"`
	Schema     *common.Schema      `bson:"schema,omitempty" json:"schema,omitempty"`
	Sort       int64               `bson:"sort" json:"sort"`
	Status     *bool               `bson:"status" json:"status"`
	CreateTime time.Time           `bson:"create_time" json:"-"`
	UpdateTime time.Time           `bson:"update_time" json:"-"`
	Children   []Content           `bson:"-" json:"children"`
}

func (x *Install) UseTemplate(ctx context.Context) (err error) {
	var template TemplateDto
	if err = x.fetchTemplate(ctx, x.Template, &template); err != nil {
		return
	}
	if err = x.validateTemplate(ctx, template.Schema, template); err != nil {
		return
	}
	if err = x.setTemplate(ctx, nil, template.Contents); err != nil {
		return
	}
	return
}

func (x *Install) fetchTemplate(ctx context.Context, url string, v interface{}) (err error) {
	client := http.DefaultClient
	var req *http.Request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = client.Do(req.WithContext(ctx)); err != nil {
		return
	}
	if err = jsoniter.NewDecoder(resp.Body).Decode(v); err != nil {
		return
	}
	return
}

func (x *Install) validateTemplate(ctx context.Context, schema string, data interface{}) (err error) {
	client := http.DefaultClient
	var req *http.Request
	if req, err = http.NewRequest("GET", schema, nil); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = client.Do(req.WithContext(ctx)); err != nil {
		return
	}
	var jschema jsonschema.Schema
	if err = jsoniter.NewDecoder(resp.Body).Decode(&jschema); err != nil {
		return
	}
	valid := jschema.Validate(ctx, data)
	if !valid.IsValid() {
		return errors.New("验证格式不一致")
	}
	return
}

func (x *Install) setTemplate(ctx context.Context, parent *primitive.ObjectID, contents []Content) (err error) {
	var keys []int
	var data []interface{}
	for k, v := range contents {
		if len(v.Children) != 0 {
			keys = append(keys, k)
		}
		if parent != nil {
			v.Parent = parent
		}
		if v.Status == nil {
			v.Status = common.BoolToP(true)
		}
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		data = append(data, v)
	}
	var result *mongo.InsertManyResult
	if result, err = x.Db.Collection("pages").
		InsertMany(ctx, data); err != nil {
		return
	}
	for _, v := range keys {
		if err = x.setTemplate(ctx,
			common.ObjectIDToP(result.InsertedIDs[v]),
			contents[v].Children,
		); err != nil {
			return
		}
	}
	return
}
