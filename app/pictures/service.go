package pictures

import (
	"api/common"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindLabels(ctx context.Context) (values []interface{}, err error) {
	if values, err = x.Db.Collection("pictures").
		Distinct(ctx, "labels", bson.M{"status": true}); err != nil {
		return
	}
	return
}

func (x *Service) BulkDelete(ctx context.Context, oids []primitive.ObjectID) (interface{}, error) {
	return x.Db.Collection("pictures").DeleteMany(ctx, bson.M{
		"_id": bson.M{"$in": oids},
	})
}

func (x *Service) ImageInfo(ctx context.Context, url string) (result map[string]interface{}, err error) {
	var response *cos.Response
	if response, err = x.Cos.CI.Get(ctx, url, "imageInfo", nil); err != nil {
		return
	}
	if err = jsoniter.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}
	return
}
