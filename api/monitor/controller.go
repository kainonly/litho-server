package monitor

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	MonitorX *Service
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
		data, err = x.MonitorX.GetMongoAvailableConnections(ctx, dto.Dates)
		break
	case "mongo_open_connections":
		data, err = x.MonitorX.GetMongoOpenConnections(ctx, dto.Dates)
		break
	case "mongo_commands_per_second":
		data, err = x.MonitorX.GetMongoCommandsPerSecond(ctx, dto.Dates)
		break
	case "mongo_query_operations":
		data, err = x.MonitorX.GetMongoQueryOperations(ctx, dto.Dates)
		break
	case "mongo_document_operations":
		data, err = x.MonitorX.GetMongoDocumentOperations(ctx, dto.Dates)
		break
	case "mongo_flushes":
		data, err = x.MonitorX.GetMongoFlushes(ctx, dto.Dates)
		break
	case "mongo_network_io":
		data, err = x.MonitorX.GetMongoNetworkIO(ctx, dto.Dates)
		break
	case "redis_mem":
		data, err = x.MonitorX.GetRedisMem(ctx, dto.Dates)
		break
	case "redis_cpu":
		data, err = x.MonitorX.GetRedisCpu(ctx, dto.Dates)
		break
	case "redis_ops_per_sec":
		data, err = x.MonitorX.GetRedisOpsPerSec(ctx, dto.Dates)
		break
	case "redis_evi_exp_keys":
		data, err = x.MonitorX.GetRedisEviExpKeys(ctx, dto.Dates)
		break
	case "redis_collections_rate":
		data, err = x.MonitorX.GetRedisCollectionsRate(ctx, dto.Dates)
		break
	case "redis_hit_rate":
		data, err = x.MonitorX.GetRedisHitRate(ctx, dto.Dates)
		break
	case "redis_network_io":
		data, err = x.MonitorX.GetRedisNetworkIO(ctx, dto.Dates)
		break
	case "nats_cpu":
		data, err = x.MonitorX.GetNatsCpu(ctx, dto.Dates)
		break
	case "nats_mem":
		data, err = x.MonitorX.GetNatsMem(ctx, dto.Dates)
		break
	case "nats_connections":
		data, err = x.MonitorX.GetNatsConnections(ctx, dto.Dates)
		break
	case "nats_subscriptions":
		data, err = x.MonitorX.GetNatsSubscriptions(ctx, dto.Dates)
		break
	case "nats_slow_consumers":
		data, err = x.MonitorX.GetNatsSlowConsumers(ctx, dto.Dates)
		break
	case "nats_msg_io":
		data, err = x.MonitorX.GetNatsMsgIO(ctx, dto.Dates)
		break
	case "nats_bytes_io":
		data, err = x.MonitorX.GetNatsBytesIO(ctx, dto.Dates)
		break
	}

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}
