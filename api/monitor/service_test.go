package monitor_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestService_GetQPS(t *testing.T) {
	data, err := x.MonitorX.GetQPS(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}
