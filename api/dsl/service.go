package dsl

import (
	"context"
	"errors"
	"github.com/weplanx/server/utils/passlib"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Service struct {
	Db *mongo.Database
}

// Transform 格式转换
func (x *Service) Transform(data M, rules []string) (err error) {
	for _, rule := range rules {
		spec := strings.Split(rule, ":")
		keys, cursor := strings.Split(spec[0], "."), data
		n := len(keys) - 1
		for _, key := range keys[:n] {
			if v, ok := cursor[key].(M); ok {
				cursor = v
			}
		}
		key := keys[n]
		if cursor[key] == nil {
			continue
		}
		switch spec[1] {
		case "oid":
			// 转换为 ObjectId
			if cursor[key], err = primitive.ObjectIDFromHex(cursor[key].(string)); err != nil {
				return
			}
			break

		case "oids":
			// 转换为 ObjectId 数组
			oids := cursor[key].([]interface{})
			for i, id := range oids {
				if oids[i], err = primitive.ObjectIDFromHex(id.(string)); err != nil {
					return
				}
			}
			break
		case "date":
			// 转换为 ISODate
			if cursor[key], err = time.Parse(time.RFC1123, cursor[key].(string)); err != nil {
				return
			}
			break

		case "password":
			// 密码类型，转换为 Argon2id
			if cursor[key], err = passlib.Hash(cursor[key].(string)); err != nil {
				if errors.Is(err, passlib.ErrNotMatch) {
					//return errs.NewPublic(0, err.Error())
					return
				}
				return
			}
			break
		}
	}
	return
}

// Create 创建文档
func (x *Service) Create(ctx context.Context, model string, doc M, xdoc []string) (_ interface{}, err error) {
	if err = x.Transform(doc, xdoc); err != nil {
		return
	}
	doc["create_time"] = time.Now()
	doc["update_time"] = time.Now()
	return x.Db.Collection(model).InsertOne(ctx, doc)
}

// BulkCreate 批量创建文档
func (x *Service) BulkCreate(ctx context.Context, model string, docs []M, xdoc []string) (_ interface{}, err error) {
	data := make([]interface{}, len(docs))
	for i, doc := range docs {
		if err = x.Transform(doc, xdoc); err != nil {
			return
		}
		doc["create_time"] = time.Now()
		doc["update_time"] = time.Now()
		data[i] = doc
	}
	return x.Db.Collection(model).InsertMany(ctx, data)
}

// Size 获取文档总数
func (x *Service) Size(ctx context.Context, model string, filter M, xfilter []string) (_ int64, err error) {
	if len(filter) == 0 {
		return x.Db.Collection(model).EstimatedDocumentCount(ctx)
	}
	if err = x.Transform(filter, xfilter); err != nil {
		return
	}
	return x.Db.Collection(model).CountDocuments(ctx, filter)
}

// Find 获取匹配文档
func (x *Service) Find(ctx context.Context, model string, filter M, xfilter []string, opt FindOption) (data []M, err error) {
	if err = x.Transform(filter, xfilter); err != nil {
		return
	}

	option := options.Find().
		SetSort(M{"_id": -1}).
		SetLimit(20).
		SetSkip(0)

	sort := opt.GetSort()
	if len(sort) != 0 {
		option.SetSort(sort)
	}

	if opt.Keys != nil {
		option = option.SetProjection(opt.Keys)
	}

	if opt.Limit != 0 {
		option = option.SetLimit(opt.Limit)
	}

	if opt.Skip != 0 {
		option = option.SetSkip(opt.Skip)
	}

	if opt.Page != 0 {
		option = option.SetSkip((opt.Page - 1) * *option.Limit)
	}

	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(model).Find(ctx, filter, option); err != nil {
		return
	}
	if err = cursor.All(ctx, &data); err != nil {
		return
	}

	return
}

// FindOne 获取单个文档
func (x *Service) FindOne(ctx context.Context, model string, filter M, xfilter []string, opt FindOption) (data M, err error) {
	if err = x.Transform(filter, xfilter); err != nil {
		return
	}

	option := options.FindOne()
	if opt.Keys != nil {
		option = option.SetProjection(opt.Keys)
	}

	if err = x.Db.Collection(model).FindOne(ctx, filter, option).Decode(&data); err != nil {
		return
	}
	return
}

// FindById 获取指定 Id 的文档
func (x *Service) FindById(ctx context.Context, model string, id string, opt FindOption) (data M, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	return x.FindOne(ctx, model, M{"_id": oid}, nil, opt)
}

// Update 局部更新多个匹配文档
func (x *Service) Update(ctx context.Context, model string, filter M, xfilter []string, update M, xdoc []string) (_ interface{}, err error) {
	if err = x.Transform(filter, xfilter); err != nil {
		return
	}
	if err = x.Transform(update, xdoc); err != nil {
		return
	}
	if _, ok := update["$set"]; !ok {
		update["$set"] = M{}
	}
	update["$set"].(M)["update_time"] = time.Now()
	return x.Db.Collection(model).UpdateMany(ctx, filter, update)
}

// UpdateById 局部更新指定 Id 的文档
func (x *Service) UpdateById(ctx context.Context, model string, id string, update M, xdoc []string) (_ interface{}, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	if err = x.Transform(update, xdoc); err != nil {
		return
	}
	if _, ok := update["$set"]; !ok {
		update["$set"] = M{}
	}
	update["$set"].(M)["update_time"] = time.Now()
	return x.Db.Collection(model).UpdateOne(ctx, M{"_id": oid}, update)
}

// Replace 替换指定 Id 的文档
func (x *Service) Replace(ctx context.Context, model string, id string, doc M, xdoc []string) (_ interface{}, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	if err = x.Transform(doc, xdoc); err != nil {
		return
	}
	doc["update_time"] = time.Now()
	return x.Db.Collection(model).ReplaceOne(ctx, M{"_id": oid}, doc)
}

// Delete 删除指定 Id 的文档
func (x *Service) Delete(ctx context.Context, model string, id string) (_ interface{}, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	return x.Db.Collection(model).DeleteOne(ctx, M{"_id": oid})
}

// BulkDelete 批量删除匹配文档
func (x *Service) BulkDelete(ctx context.Context, model string, filter M, xfilter []string) (_ interface{}, err error) {
	if err = x.Transform(filter, xfilter); err != nil {
		return
	}
	return x.Db.Collection(model).DeleteMany(ctx, filter)
}

// Sort 通用排序
func (x *Service) Sort(ctx context.Context, model string, oids []primitive.ObjectID) (_ interface{}, err error) {
	var wms []mongo.WriteModel
	for i, oid := range oids {
		update := M{
			"$set": M{
				"sort":        i,
				"update_time": time.Now(),
			},
		}

		wms = append(wms, mongo.NewUpdateOneModel().
			SetFilter(M{"_id": oid}).
			SetUpdate(update),
		)
	}
	return x.Db.Collection(model).BulkWrite(ctx, wms)
}
