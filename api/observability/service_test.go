package observability_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetQPS(t *testing.T) {
	data, err := x.ObservabilityService.GetQpsRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetErrorRate(t *testing.T) {
	data, err := x.ObservabilityService.GetErrorRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetGoroutines(t *testing.T) {
	data, err := x.ObservabilityService.GetGoroutines(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetUptime(t *testing.T) {
	data, err := x.ObservabilityService.GetUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetLookups(t *testing.T) {
	data, err := x.ObservabilityService.GetLookups(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetCgCount(t *testing.T) {
	data, err := x.ObservabilityService.GetGcCount(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetCgoCalls(t *testing.T) {
	data, err := x.ObservabilityService.GetCgoCalls(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoUptime(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoAvailableConnections(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoAvailableConnections(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoOpenConnections(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoOpenConnections(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoCommandsPerSecond(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoCommandsPerSecond(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoQueryOperations(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoQueryOperations(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoDocumentOperations(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoDocumentOperations(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoFlushes(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoFlushes(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoNetworkIO(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoNetworkIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisUptime(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisMem(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisMem(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCpu(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisCpu(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisOpsPerSec(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisOpsPerSec(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisEviExpKeys(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisEviExpKeys(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCollectionsRate(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisCollectionsRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisConnectedSlaves(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisConnectedSlaves(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisHitRate(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisHitRate(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisNetworkIO(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisNetworkIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsUptime(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsUptime(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsCpu(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsCpu(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMem(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsMem(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsConnections(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsConnections(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSubscriptions(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsSubscriptions(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSlowConsumers(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsSlowConsumers(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMsgIO(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsMsgIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsBytesIO(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsBytesIO(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}
