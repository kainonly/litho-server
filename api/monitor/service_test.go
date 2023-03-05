package monitor_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetQPS(t *testing.T) {
	data, err := x.MonitorX.GetQpsRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetErrorRate(t *testing.T) {
	data, err := x.MonitorX.GetErrorRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetCgoCalls(t *testing.T) {
	data, err := x.MonitorX.GetCgoCalls(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoUptime(t *testing.T) {
	data, err := x.MonitorX.GetMongoUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoAvailableConnections(t *testing.T) {
	data, err := x.MonitorX.GetMongoAvailableConnections(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoOpenConnections(t *testing.T) {
	data, err := x.MonitorX.GetMongoOpenConnections(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoCommandsPerSecond(t *testing.T) {
	data, err := x.MonitorX.GetMongoCommandsPerSecond(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoQueryOperations(t *testing.T) {
	data, err := x.MonitorX.GetMongoQueryOperations(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoDocumentOperations(t *testing.T) {
	data, err := x.MonitorX.GetMongoDocumentOperations(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoFlushes(t *testing.T) {
	data, err := x.MonitorX.GetMongoFlushes(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoNetworkIO(t *testing.T) {
	data, err := x.MonitorX.GetMongoNetworkIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisUptime(t *testing.T) {
	data, err := x.MonitorX.GetRedisUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisMem(t *testing.T) {
	data, err := x.MonitorX.GetRedisMem(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCpu(t *testing.T) {
	data, err := x.MonitorX.GetRedisCpu(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisOpsPerSec(t *testing.T) {
	data, err := x.MonitorX.GetRedisOpsPerSec(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisEviExpKeys(t *testing.T) {
	data, err := x.MonitorX.GetRedisEviExpKeys(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCollectionsRate(t *testing.T) {
	data, err := x.MonitorX.GetRedisCollectionsRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisConnectedSlaves(t *testing.T) {
	data, err := x.MonitorX.GetRedisConnectedSlaves(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisHitRate(t *testing.T) {
	data, err := x.MonitorX.GetRedisHitRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisNetworkIO(t *testing.T) {
	data, err := x.MonitorX.GetRedisNetworkIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsUptime(t *testing.T) {
	data, err := x.MonitorX.GetNatsUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsCpu(t *testing.T) {
	data, err := x.MonitorX.GetNatsCpu(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMem(t *testing.T) {
	data, err := x.MonitorX.GetNatsMem(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsConnections(t *testing.T) {
	data, err := x.MonitorX.GetNatsConnections(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSubscriptions(t *testing.T) {
	data, err := x.MonitorX.GetNatsSubscriptions(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSlowConsumers(t *testing.T) {
	data, err := x.MonitorX.GetNatsSlowConsumers(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMsgIO(t *testing.T) {
	data, err := x.MonitorX.GetNatsMsgIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsBytesIO(t *testing.T) {
	data, err := x.MonitorX.GetNatsBytesIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}
