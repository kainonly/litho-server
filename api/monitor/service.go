package monitor

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/weplanx/server/common"
	"strings"
	"time"
)

type Service struct {
	*common.Inject

	Flux influxdb2.Client
}

func (x *Service) Range(v string) string {
	if v == "" {
		return `|> range(start: -15m, stop: now())`
	}
	dates := strings.Split(v, ",")
	return fmt.Sprintf(`|> range(start: %s, stop: %s)`,
		dates[0], dates[1],
	)
}

func (x *Service) GetMongoAvailableConnections(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "connections_available")
	  	|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	  	|> yield(name: "mean")
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoOpenConnections(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "open_connections")
  		|> derivative(unit: 1m,nonNegative: true)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}
	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoCommandsPerSecond(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "commands_per_sec")
  		|> derivative(unit: 1m,nonNegative: true)
		|> fill(value: float(v: 0))
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoQueryOperations(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "commands" or r["_field"] == "deletes" or r["_field"] == "getmores" or r["_field"] == "inserts" or r["_field"] == "updates")
  		|> derivative(unit: 1m,nonNegative: true)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"commands": 0,
		"getmores": 1,
		"inserts":  2,
		"updates":  3,
		"deletes":  4,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoDocumentOperations(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "document_deleted" or r["_field"] == "document_inserted" or r["_field"] == "document_returned" or r["_field"] == "document_updated")
  		|> derivative(unit: 1m,nonNegative: true)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"document_returned": 0,
		"document_inserted": 1,
		"document_updated":  2,
		"document_deleted":  3,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoFlushes(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "flushes")
  		|> derivative(unit: 1m,nonNegative: true)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoNetworkIO(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "net_in_bytes" or r["_field"] == "net_out_bytes")
  		|> derivative(unit: 1m,nonNegative: true)
		|> fill(value: float(v: 0))
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"net_in_bytes":  0,
		"net_out_bytes": 1,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisCpu(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "redis")
		|> filter(fn: (r) => 
			r["_field"] == "used_cpu_user" or 
			r["_field"] == "used_cpu_sys" or 
			r["_field"] == "used_cpu_sys_children" or 
			r["_field"] == "used_cpu_user_children"
		)
  		|> derivative(unit: 1m, nonNegative: true)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"used_cpu_user":          0,
		"used_cpu_sys":           1,
		"used_cpu_sys_children":  2,
		"used_cpu_user_children": 3,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisMem(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r["_measurement"] == "redis")
		|> filter(fn: (r) => 
			r["_field"] == "used_memory" or 
			r["_field"] == "used_memory_dataset" or 
			r["_field"] == "used_memory_rss" or 
			r["_field"] == "used_memory_lua"
		)
  		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"used_memory":         0,
		"used_memory_dataset": 1,
		"used_memory_rss":     2,
		"used_memory_lua":     3,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisOpsPerSec(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "instantaneous_ops_per_sec")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisEviExpKeys(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "evicted_keys" or r._field == "expired_keys")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"evicted_keys": 0,
		"expired_keys": 1,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisCollectionsRate(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "total_connections_received" or r._field == "rejected_connections")
		|> derivative(unit: 1m, nonNegative: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisHitRate(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "keyspace_hitrate")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisNetworkIO(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "total_net_output_bytes" or r._field == "total_net_input_bytes")
		|> derivative(unit: 1m, nonNegative: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"total_net_input_bytes":  0,
		"total_net_output_bytes": 1,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsCpu(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "cpu")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			result.Record().ValueByKey("server"),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsMem(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "mem")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			result.Record().ValueByKey("server"),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsConnections(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "connections")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			result.Record().ValueByKey("server"),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsSubscriptions(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "subscriptions")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsSlowConsumers(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "slow_consumers")
		|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsMsgIO(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "in_msgs" or r._field == "out_msgs")
		|> derivative(unit: 1m, nonNegative: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"in_msgs":  0,
		"out_msgs": 1,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsBytesIO(ctx context.Context, dates string) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		%s
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "in_bytes" or r._field == "out_bytes")
		|> derivative(unit: 1m, nonNegative: false)
	`, x.V.Influx.Bucket, x.Range(dates))
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	index := M{
		"in_bytes":  0,
		"out_bytes": 1,
	}
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			index[result.Record().Field()],
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}
