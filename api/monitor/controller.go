package monitor

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	MonitorService *Service
}

type ExportersDto struct {
	Name  string `path:"name" vd:"required"`
	Dates string `query:"dates"`
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
	case "mongo_available_connections":
		data, err = x.MonitorService.GetMongoAvailableConnections(ctx, dto.Dates)
		break
	case "mongo_open_connections":
		data, err = x.MonitorService.GetMongoOpenConnections(ctx, dto.Dates)
		break
	case "mongo_commands_per_second":
		data, err = x.MonitorService.GetMongoCommandsPerSecond(ctx, dto.Dates)
		break
	case "mongo_query_operations":
		data, err = x.MonitorService.GetMongoQueryOperations(ctx, dto.Dates)
		break
	case "mongo_document_operations":
		data, err = x.MonitorService.GetMongoDocumentOperations(ctx, dto.Dates)
		break
	case "mongo_flushes":
		data, err = x.MonitorService.GetMongoFlushes(ctx, dto.Dates)
		break
	case "mongo_network_io":
		data, err = x.MonitorService.GetMongoNetworkIO(ctx, dto.Dates)
		break
	case "redis_mem":
		data, err = x.MonitorService.GetRedisMem(ctx, dto.Dates)
		break
	case "redis_cpu":
		data, err = x.MonitorService.GetRedisCpu(ctx, dto.Dates)
		break
	case "redis_ops_per_sec":
		data, err = x.MonitorService.GetRedisOpsPerSec(ctx, dto.Dates)
		break
	case "redis_evi_exp_keys":
		data, err = x.MonitorService.GetRedisEviExpKeys(ctx, dto.Dates)
		break
	case "redis_collections_rate":
		data, err = x.MonitorService.GetRedisCollectionsRate(ctx, dto.Dates)
		break
	case "redis_hit_rate":
		data, err = x.MonitorService.GetRedisHitRate(ctx, dto.Dates)
		break
	case "redis_network_io":
		data, err = x.MonitorService.GetRedisNetworkIO(ctx, dto.Dates)
		break
	case "nats_cpu":
		data, err = x.MonitorService.GetNatsCpu(ctx, dto.Dates)
		break
	case "nats_mem":
		data, err = x.MonitorService.GetNatsMem(ctx, dto.Dates)
		break
	case "nats_connections":
		data, err = x.MonitorService.GetNatsConnections(ctx, dto.Dates)
		break
	case "nats_subscriptions":
		data, err = x.MonitorService.GetNatsSubscriptions(ctx, dto.Dates)
		break
	case "nats_slow_consumers":
		data, err = x.MonitorService.GetNatsSlowConsumers(ctx, dto.Dates)
		break
	case "nats_msg_io":
		data, err = x.MonitorService.GetNatsMsgIO(ctx, dto.Dates)
		break
	case "nats_bytes_io":
		data, err = x.MonitorService.GetNatsBytesIO(ctx, dto.Dates)
		break
	}

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}
