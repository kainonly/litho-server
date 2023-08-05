package observability

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Controller struct {
	ObservabilityService *Service
}

type ExportersDto struct {
	Name string `path:"name,required"`
}

func (x *Controller) Exporters(ctx context.Context, c *app.RequestContext) {
	var dto ExportersDto
	var err error
	if err = c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	var data interface{}
	switch dto.Name {
	case "qps_rate":
		data, err = x.ObservabilityService.GetQpsRate(ctx)
		break
	case "error_rate":
		data, err = x.ObservabilityService.GetErrorRate(ctx)
		break
	case "p99":
		data, err = x.ObservabilityService.GetP99(ctx)
		break
	case "mongo_available_connections":
		data, err = x.ObservabilityService.GetMongoAvailableConnections(ctx)
		break
	case "mongo_open_connections":
		data, err = x.ObservabilityService.GetMongoOpenConnections(ctx)
		break
	case "mongo_commands_per_second":
		data, err = x.ObservabilityService.GetMongoCommandsPerSecond(ctx)
		break
	case "mongo_query_operations":
		data, err = x.ObservabilityService.GetMongoQueryOperations(ctx)
		break
	case "mongo_document_operations":
		data, err = x.ObservabilityService.GetMongoDocumentOperations(ctx)
		break
	case "mongo_flushes":
		data, err = x.ObservabilityService.GetMongoFlushes(ctx)
		break
	case "mongo_network_io":
		data, err = x.ObservabilityService.GetMongoNetworkIO(ctx)
		break
	case "redis_mem":
		data, err = x.ObservabilityService.GetRedisMem(ctx)
		break
	case "redis_cpu":
		data, err = x.ObservabilityService.GetRedisCpu(ctx)
		break
	case "redis_ops_per_sec":
		data, err = x.ObservabilityService.GetRedisOpsPerSec(ctx)
		break
	case "redis_evi_exp_keys":
		data, err = x.ObservabilityService.GetRedisEviExpKeys(ctx)
		break
	case "redis_collections_rate":
		data, err = x.ObservabilityService.GetRedisCollectionsRate(ctx)
		break
	case "redis_hit_rate":
		data, err = x.ObservabilityService.GetRedisHitRate(ctx)
		break
	case "redis_network_io":
		data, err = x.ObservabilityService.GetRedisNetworkIO(ctx)
		break
	case "nats_cpu":
		data, err = x.ObservabilityService.GetNatsCpu(ctx)
		break
	case "nats_mem":
		data, err = x.ObservabilityService.GetNatsMem(ctx)
		break
	case "nats_connections":
		data, err = x.ObservabilityService.GetNatsConnections(ctx)
		break
	case "nats_subscriptions":
		data, err = x.ObservabilityService.GetNatsSubscriptions(ctx)
		break
	case "nats_slow_consumers":
		data, err = x.ObservabilityService.GetNatsSlowConsumers(ctx)
		break
	case "nats_msg_io":
		data, err = x.ObservabilityService.GetNatsMsgIO(ctx)
		break
	case "nats_bytes_io":
		data, err = x.ObservabilityService.GetNatsBytesIO(ctx)
		break
	case "mem_heap_sys":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.heap_sys")
		break
	case "mem_heap_released":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.heap_released")
		break
	case "mem_heap_alloc":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.heap_alloc")
		break
	case "mem_heap_inuse":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.heap_inuse")
		break
	case "mem_heap_idle":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.heap_idle")
		break
	case "mem_heap_objects":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.heap_objects")
		break
	case "mem_live_objects":
		data, err = x.ObservabilityService.GetRuntime(ctx, "mem.live_objects")
		break
	case "goroutines":
		data, err = x.ObservabilityService.GetRuntime(ctx, "goroutines")
		break
	case "mem_lookups":
		data, err = x.ObservabilityService.GetRuntimeLast(context.TODO(), "process.runtime.go.mem.lookups")
		break
	case "cgo_calls":
		data, err = x.ObservabilityService.GetRuntimeLast(context.TODO(), "process.runtime.go.cgo.calls")
		break
	case "gc_count":
		data, err = x.ObservabilityService.GetRuntimeLast(context.TODO(), "process.runtime.go.gc.count")
		break
	case "uptime":
		data, err = x.ObservabilityService.GetRuntimeLast(context.TODO(), "runtime.uptime")
		break
	}

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}
