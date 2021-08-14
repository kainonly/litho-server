package service

import (
	"context"
	"errors"
	"time"
)

var Inconsistent = errors.New("verification is inconsistent")

type Index struct {
	*Dependency
}

func NewIndex(d Dependency) *Index {
	return &Index{Dependency: &d}
}

func (x *Index) SetCode(key string, code string) error {
	return x.Redis.Set(context.Background(), "code:"+key, code, time.Minute).Err()
}

func (x *Index) CheckCode(key string, code string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := x.Redis.Exists(ctx, "code:"+key).Result()
	if err != nil {
		return
	}
	if result == 0 {
		return Inconsistent
	}
	text, err := x.Redis.Get(ctx, "code:"+key).Result()
	if err != nil {
		return
	}
	if text != code {
		return Inconsistent
	}
	return nil
}

func (x *Index) DelCode(key string) error {
	return x.Redis.Del(context.Background(), "code:"+key).Err()
}
