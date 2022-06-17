package schedules

import (
	"api/common"
	"api/model"
	"context"
	"github.com/weplanx/schedule/client"
	scheduleCommon "github.com/weplanx/schedule/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
	Client *client.Schedule
}

// List 调度服务已存在的标识
func (x *Service) List() (keys []string, err error) {
	return x.Client.List()
}

// Get 获取指定服务配置与运行状态
func (x *Service) Get(id string) ([]scheduleCommon.Job, error) {
	return x.Client.Get(id)
}

// Sync 同步服务
func (x *Service) Sync(ctx context.Context, id primitive.ObjectID) (err error) {
	var data model.Schedule
	if err = x.Db.Collection("schedules").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&data); err != nil {
		return
	}
	var jobs []scheduleCommon.Job
	for _, v := range data.Jobs {
		jobs = append(jobs, scheduleCommon.Job{
			Mode:   v.Mode,
			Spec:   v.Spec,
			Option: v.Option,
		})
	}
	if err = x.Client.Set(id.Hex(), jobs...); err != nil {
		return
	}
	return
}

// Delete 删除服务
func (x *Service) Delete(id string) (err error) {
	return x.Client.Remove(id)
}
