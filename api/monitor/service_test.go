package monitor_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetMongoAvailableConnections(t *testing.T) {
	data, err := x.MonitorService.GetMongoAvailableConnections(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoOpenConnections(t *testing.T) {
	data, err := x.MonitorService.GetMongoOpenConnections(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoCommandsPerSecond(t *testing.T) {
	data, err := x.MonitorService.GetMongoCommandsPerSecond(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoQueryOperations(t *testing.T) {
	data, err := x.MonitorService.GetMongoQueryOperations(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoDocumentOperations(t *testing.T) {
	data, err := x.MonitorService.GetMongoDocumentOperations(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoFlushes(t *testing.T) {
	data, err := x.MonitorService.GetMongoFlushes(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoNetworkIO(t *testing.T) {
	data, err := x.MonitorService.GetMongoNetworkIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisMem(t *testing.T) {
	data, err := x.MonitorService.GetRedisMem(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCpu(t *testing.T) {
	data, err := x.MonitorService.GetRedisCpu(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisOpsPerSec(t *testing.T) {
	data, err := x.MonitorService.GetRedisOpsPerSec(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisEviExpKeys(t *testing.T) {
	data, err := x.MonitorService.GetRedisEviExpKeys(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCollectionsRate(t *testing.T) {
	data, err := x.MonitorService.GetRedisCollectionsRate(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisHitRate(t *testing.T) {
	data, err := x.MonitorService.GetRedisHitRate(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisNetworkIO(t *testing.T) {
	data, err := x.MonitorService.GetRedisNetworkIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsCpu(t *testing.T) {
	data, err := x.MonitorService.GetNatsCpu(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMem(t *testing.T) {
	data, err := x.MonitorService.GetNatsMem(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsConnections(t *testing.T) {
	data, err := x.MonitorService.GetNatsConnections(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSubscriptions(t *testing.T) {
	data, err := x.MonitorService.GetNatsSubscriptions(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSlowConsumers(t *testing.T) {
	data, err := x.MonitorService.GetNatsSlowConsumers(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMsgIO(t *testing.T) {
	data, err := x.MonitorService.GetNatsMsgIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsBytesIO(t *testing.T) {
	data, err := x.MonitorService.GetNatsBytesIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}
