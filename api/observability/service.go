package observability

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/weplanx/server/common"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) GetQpsRate(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`
		import "experimental/aggregate"

		data =
			from(bucket: "%s")
				|> range(start: -15m, stop: now())
				|> filter(fn: (r) => r["_measurement"] == "prometheus")
				|> filter(fn: (r) => r["service.name"] == "%s")
				|> filter(fn: (r) => r["_field"] == "http.server.duration_count")

		all =
			data
				|> aggregate.rate(every: 10s, unit: 1s, groupColumns: ["service.name"])
				|> set(key: "method", value: "ALL")
				|> fill(value: 0.0)

		method =
			data
				|> aggregate.rate(every: 10s, unit: 1s, groupColumns: ["service.name", "http.method"])
				|> rename(columns: {"http.method": "method"})
				|> fill(value: 0.0)
		
		union(tables: [method, all])
	`, x.V.Influx.Bucket, x.V.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}
	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Local().Format(time.TimeOnly),
			result.Record().Value(),
			result.Record().ValueByKey("method"),
		})
	}
	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetErrorRate(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`
		import "experimental/aggregate"

		data =
			from(bucket: "%s")
				|> range(start: -15m, stop: now())
				|> filter(fn: (r) => r["_measurement"] == "prometheus")
				|> filter(fn: (r) => r["service.name"] == "%s")
				|> filter(fn: (r) => r["_field"] == "http.server.duration_count")

		all =
			data
				|> aggregate.rate(every: 10s, unit: 1s, groupColumns: ["service.name"])
				|> set(key: "type", value: "ALL")
				|> fill(value: 0.0)

		err =
			data
				|> filter(fn: (r) => r["http.status_code"] =~ /^[4,5]/)
				|> aggregate.rate(every: 10s, unit: 1s, groupColumns: ["service.name"])
				|> set(key: "type", value: "ERR")
				|> fill(value: 0.0)
		
		union(tables: [all, err])
			|> pivot(rowKey: ["_time"], columnKey: ["type"], valueColumn: "_value")
			|> map(fn: (r) => ({r with _value: r.ERR / r.ALL}))
	`, x.V.Influx.Bucket, x.V.Namespace)
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

func (x *Service) GetGoroutines(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "prometheus")
		|> filter(fn: (r) => r["_field"] == "process.runtime.go.goroutines")
		|> filter(fn: (r) => r["service.name"] == "%s")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket, x.V.Namespace)
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

func (x *Service) GetGcCount(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "prometheus")
		|> filter(fn: (r) => r["_field"] == "process.runtime.go.gc.count")
		|> filter(fn: (r) => r["service.name"] == "%s")
	  	|> last()
	`, x.V.Influx.Bucket, x.V.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	for result.Next() {
		value = result.Record().Value()
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetLookups(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "prometheus")
		|> filter(fn: (r) => r["_field"] == "process.runtime.go.mem.lookups")
		|> filter(fn: (r) => r["service.name"] == "%s")
	  	|> last()
	`, x.V.Influx.Bucket, x.V.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	for result.Next() {
		value = result.Record().Value()
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetCgoCalls(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "prometheus")
		|> filter(fn: (r) => r["_field"] == "process.runtime.go.cgo.calls")
		|> filter(fn: (r) => r["service.name"] == "%s")
	  	|> last()
	`, x.V.Influx.Bucket, x.V.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	for result.Next() {
		value = result.Record().Value()
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoAvailableConnections(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "connections_available")
	  	|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	  	|> yield(name: "mean")
	`, x.V.Influx.Bucket)
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

func (x *Service) GetMongoOpenConnections(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "open_connections")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetMongoCommandsPerSecond(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "commands_per_sec")
  		|> derivative(unit: 10s,nonNegative: true)
		|> fill(value: float(v: 0))
	`, x.V.Influx.Bucket)
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

func (x *Service) GetMongoQueryOperations(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "commands" or r["_field"] == "deletes" or r["_field"] == "getmores" or r["_field"] == "inserts" or r["_field"] == "updates")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetMongoDocumentOperations(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "document_deleted" or r["_field"] == "document_inserted" or r["_field"] == "document_returned" or r["_field"] == "document_updated")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetMongoFlushes(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "flushes")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetMongoNetworkIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "net_in_bytes" or r["_field"] == "net_out_bytes")
  		|> derivative(unit: 10s,nonNegative: true)
		|> fill(value: float(v: 0))
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisCpu(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "redis")
		|> filter(fn: (r) => 
			r["_field"] == "used_cpu_user" or 
			r["_field"] == "used_cpu_sys" or 
			r["_field"] == "used_cpu_sys_children" or 
			r["_field"] == "used_cpu_user_children"
		)
  		|> derivative(unit: 10s, nonNegative: true)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisMem(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "redis")
		|> filter(fn: (r) => 
			r["_field"] == "used_memory" or 
			r["_field"] == "used_memory_dataset" or 
			r["_field"] == "used_memory_rss" or 
			r["_field"] == "used_memory_lua"
		)
  		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisOpsPerSec(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "instantaneous_ops_per_sec")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisEviExpKeys(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "evicted_keys" or r._field == "expired_keys")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisCollectionsRate(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "total_connections_received" or r._field == "rejected_connections")
		|> derivative(unit: 10s, nonNegative: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisHitRate(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "keyspace_hitrate")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetRedisNetworkIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "total_net_output_bytes" or r._field == "total_net_input_bytes")
		|> derivative(unit: 10s, nonNegative: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsCpu(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "cpu")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsMem(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "mem")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsConnections(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "connections")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsSubscriptions(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "subscriptions")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsSlowConsumers(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "slow_consumers")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsMsgIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "in_msgs" or r._field == "out_msgs")
		|> derivative(unit: 10s, nonNegative: false)
	`, x.V.Influx.Bucket)
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

func (x *Service) GetNatsBytesIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Flux.QueryAPI(x.V.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "in_bytes" or r._field == "out_bytes")
		|> derivative(unit: 10s, nonNegative: false)
	`, x.V.Influx.Bucket)
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
