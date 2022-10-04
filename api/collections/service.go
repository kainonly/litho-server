package collections

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Service struct {
	Db *mongo.Database
}

// Create 新增文档
func (x *Service) Create(ctx context.Context, name string, doc M) (_ interface{}, err error) {
	return x.Db.Collection(name).InsertOne(ctx, doc)
}

// BulkCreate 批量新增文档
func (x *Service) BulkCreate(ctx context.Context, name string, docs []interface{}) (_ interface{}, err error) {
	return x.Db.Collection(name).InsertMany(ctx, docs)
}

// Size 获取文档总数
func (x *Service) Size(ctx context.Context, name string, filter M) (_ int64, err error) {
	if len(filter) == 0 {
		return x.Db.Collection(name).EstimatedDocumentCount(ctx)
	}
	return x.Db.Collection(name).CountDocuments(ctx, filter)
}

// Find 获取匹配文档
func (x *Service) Find(ctx context.Context, name string, filter M, option *options.FindOptions) (data []M, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(name).Find(ctx, filter, option); err != nil {
		return
	}
	data = make([]M, 0)
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}

// FindOne 获取单个文档
func (x *Service) FindOne(ctx context.Context, name string, filter M, option *options.FindOneOptions) (data M, err error) {
	if err = x.Db.Collection(name).FindOne(ctx, filter, option).Decode(&data); err != nil {
		return
	}
	return
}

// Update 局部更新匹配文档
func (x *Service) Update(ctx context.Context, name string, filter M, update M) (_ interface{}, err error) {
	return x.Db.Collection(name).UpdateMany(ctx, filter, update)
}

// UpdateById 局部更新指定 ID 的文档
func (x *Service) UpdateById(ctx context.Context, name string, id primitive.ObjectID, update M) (_ interface{}, err error) {
	return x.Db.Collection(name).UpdateOne(ctx, M{"_id": id}, update)
}

// Replace 替换指定 ID 的文档
func (x *Service) Replace(ctx context.Context, name string, id primitive.ObjectID, doc M) (_ interface{}, err error) {
	return x.Db.Collection(name).ReplaceOne(ctx, M{"_id": id}, doc)
}

// Delete 删除指定 ID 的文档
func (x *Service) Delete(ctx context.Context, name string, id primitive.ObjectID) (_ interface{}, err error) {
	return x.Db.Collection(name).DeleteOne(ctx, M{"_id": id})
}

// BulkDelete 批量删除匹配文档
func (x *Service) BulkDelete(ctx context.Context, name string, filter M) (_ interface{}, err error) {
	return x.Db.Collection(name).DeleteMany(ctx, filter)
}

// Sort 排序文档
func (x *Service) Sort(ctx context.Context, name string, ids []primitive.ObjectID) (_ interface{}, err error) {
	var wms []mongo.WriteModel
	for i, id := range ids {
		update := M{
			"$set": M{
				"sort":        i,
				"update_time": time.Now(),
			},
		}

		wms = append(wms, mongo.NewUpdateOneModel().
			SetFilter(M{"_id": id}).
			SetUpdate(update),
		)
	}
	return x.Db.Collection(name).BulkWrite(ctx, wms)
}
