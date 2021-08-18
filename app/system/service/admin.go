package service

import (
	"context"
	"lab-api/model"
	"log"
)

type Admin struct {
	*Dependency
	Key string
}

func NewAdmin(d Dependency) *Admin {
	return &Admin{
		Dependency: &d,
		Key:        d.Config.RedisKey("admin"),
	}
}

func (x *Admin) FindByUsername(username string) (data model.Admin, err error) {
	if err = x.Db.
		Where("username = ?", username).
		Where("status = ?", true).
		First(&data).Error; err != nil {
		return
	}
	return
}

func (x *Admin) RefreshCache(ctx context.Context) (err error) {
	var data []map[string]interface{}
	if err = x.Db.Model(&model.Admin{}).
		Where("status = ?", true).
		Order("sort").
		Find(&data).Error; err != nil {
		return
	}
	log.Println(data)
	//var value []byte
	//if value, err = jsoniter.Marshal(&data); err != nil {
	//	return
	//}
	//if err = x.Redis.Set(ctx, x.Key, value, 0).Err(); err != nil {
	//	return
	//}
	return
}

func (x *Admin) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
