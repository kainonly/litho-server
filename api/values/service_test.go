package values_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Load(t *testing.T) {
	// 清除默认 Object
	x.KeyValue.Delete("values")
	// 当 Object 不存在时初始并加载动态配置
	err := x.ValuesService.Load(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, x.Values.DynamicValues, dynamic_values.Default)
	// 当 Object 存在时加载动态配置
	err = x.ValuesService.Load(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, x.Values.DynamicValues, dynamic_values.Default)

}

func TestService_Sync(t *testing.T) {

}

func TestService_Get(t *testing.T) {

}

func TestService_Set(t *testing.T) {

}

func TestService_Remove(t *testing.T) {

}
