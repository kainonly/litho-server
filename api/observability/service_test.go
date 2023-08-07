package observability_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetQPS(t *testing.T) {
	data, err := x.ObservabilityService.GetQpsRate(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetErrorRate(t *testing.T) {
	data, err := x.ObservabilityService.GetErrorRate(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetP99(t *testing.T) {
	data, err := x.ObservabilityService.GetP99(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoAvailableConnections(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoAvailableConnections(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoOpenConnections(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoOpenConnections(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoCommandsPerSecond(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoCommandsPerSecond(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoQueryOperations(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoQueryOperations(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoDocumentOperations(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoDocumentOperations(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoFlushes(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoFlushes(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMongoNetworkIO(t *testing.T) {
	data, err := x.ObservabilityService.GetMongoNetworkIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisMem(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisMem(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCpu(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisCpu(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisOpsPerSec(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisOpsPerSec(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisEviExpKeys(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisEviExpKeys(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisCollectionsRate(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisCollectionsRate(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisHitRate(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisHitRate(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetRedisNetworkIO(t *testing.T) {
	data, err := x.ObservabilityService.GetRedisNetworkIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsCpu(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsCpu(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMem(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsMem(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsConnections(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsConnections(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSubscriptions(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsSubscriptions(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsSlowConsumers(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsSlowConsumers(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsMsgIO(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsMsgIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetNatsBytesIO(t *testing.T) {
	data, err := x.ObservabilityService.GetNatsBytesIO(context.TODO(), "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemHeapSys(t *testing.T) {
	// 从操作系统获得的堆内存
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.heap_sys", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemHeapReleased(t *testing.T) {
	// 已交还给操作系统的堆内存
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.heap_released", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemHeapAlloc(t *testing.T) {
	// 分配的堆对象的字节数
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.heap_alloc", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemHeapInuse(t *testing.T) {
	// 已使用的堆内存
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.heap_inuse", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemHeapIdle(t *testing.T) {
	// 空闲（未使用）的堆内存
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.heap_idle", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemHeapObjects(t *testing.T) {
	// 已分配的堆对象数量
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.heap_objects", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemLiveObjects(t *testing.T) {
	// 活动对象的数量
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "mem.live_objects", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetGoroutines(t *testing.T) {
	data, err := x.ObservabilityService.GetRuntime(context.TODO(), "goroutines", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetCgCount(t *testing.T) {
	// GC 累计
	data, err := x.ObservabilityService.GetRuntimeLast(context.TODO(), "process.runtime.go.gc.count", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetMemLookups(t *testing.T) {
	// 指针查询
	data, err := x.ObservabilityService.GetRuntimeLast(context.TODO(), "process.runtime.go.mem.lookups", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetCgoCalls(t *testing.T) {
	// CGO
	data, err := x.ObservabilityService.GetRuntimeLast(context.TODO(), "process.runtime.go.cgo.calls", "")
	assert.NoError(t, err)
	t.Log(data)
}

func TestService_GetUptime(t *testing.T) {
	// CGO
	data, err := x.ObservabilityService.GetRuntimeLast(context.TODO(), "runtime.uptime", "")
	assert.NoError(t, err)
	t.Log(data)
}
