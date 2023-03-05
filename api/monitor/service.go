package monitor

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
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
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
	`, x.Values.Influx.Bucket, x.Values.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}
	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
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
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
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
	`, x.Values.Influx.Bucket, x.Values.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}
	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}
	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetCgoCalls(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "prometheus")
		|> filter(fn: (r) => r["_field"] == "process.runtime.go.cgo.calls")
		|> filter(fn: (r) => r["service.name"] == "%s")
		|> group(columns: ["service.name"], mode: "by")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
		|> yield(name: "mean")
	`, x.Values.Influx.Bucket, x.Values.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}
	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}
	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoUptime(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "uptime_ns")
	  	|> last()
  		|> map(fn: (r) => ({r with _value: r._value / 3600000000000}))
	`, x.Values.Influx.Bucket)
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
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "connections_available")
	  	|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
	  	|> yield(name: "mean")
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoOpenConnections(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "open_connections")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoCommandsPerSecond(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "commands_per_sec")
  		|> derivative(unit: 10s,nonNegative: true)
		|> fill(value: float(v: 0))
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) ToMongoQueryOperate(operate string) int {
	switch operate {
	case "commands":
		return 0
	case "getmores":
		return 1
	case "inserts":
		return 2
	case "updates":
		return 3
	case "deletes":
		return 4
	}
	return 0
}

func (x *Service) GetMongoQueryOperations(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "commands" or r["_field"] == "deletes" or r["_field"] == "getmores" or r["_field"] == "inserts" or r["_field"] == "updates")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			x.ToMongoQueryOperate(result.Record().Field()),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) ToMongoDocumentOperate(operate string) int {
	switch operate {
	case "document_returned":
		return 0
	case "document_inserted":
		return 1
	case "document_updated":
		return 2
	case "document_deleted":
		return 3
	}
	return 0
}

func (x *Service) GetMongoDocumentOperations(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
	  	|> filter(fn: (r) => r["_field"] == "document_deleted" or r["_field"] == "document_inserted" or r["_field"] == "document_returned" or r["_field"] == "document_updated")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			x.ToMongoDocumentOperate(result.Record().Field()),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetMongoFlushes(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "flushes")
  		|> derivative(unit: 10s,nonNegative: true)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) ToMongoNetworkIO(v string) int {
	switch v {
	case "net_in_bytes":
		return 0
	case "net_out_bytes":
		return 1
	}
	return 0
}

func (x *Service) GetMongoNetworkIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "mongodb")
		|> filter(fn: (r) => r["_field"] == "net_in_bytes" or r["_field"] == "net_out_bytes")
  		|> derivative(unit: 10s,nonNegative: true)
		|> fill(value: float(v: 0))
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			x.ToMongoNetworkIO(result.Record().Field()),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisUptime(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "redis")
		|> filter(fn: (r) => r._field == "uptime")
	  	|> last()
	`, x.Values.Influx.Bucket)
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

func (x *Service) ToRedisMem(v string) int {
	switch v {
	case "used_memory":
		return 0
	case "used_memory_dataset":
		return 1
	case "used_memory_rss":
		return 2
	case "used_memory_lua":
		return 3
	}
	return 0
}

func (x *Service) GetRedisMem(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
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
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			x.ToRedisMem(result.Record().Field()),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) ToRedisCpu(v string) int {
	switch v {
	case "used_cpu_user":
		return 0
	case "used_cpu_sys":
		return 1
	case "used_cpu_sys_children":
		return 2
	case "used_cpu_user_children":
		return 3
	}
	return 0
}

func (x *Service) GetRedisCpu(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "redis")
		|> filter(fn: (r) => 
			r["_field"] == "used_cpu_user" or 
			r["_field"] == "used_cpu_sys" or 
			r["_field"] == "used_cpu_sys_children" or 
			r["_field"] == "used_cpu_user_children"
		)
  		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			x.ToRedisCpu(result.Record().Field()),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisOpsPerSec(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "instantaneous_ops_per_sec")
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisEviExpKeys(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "evicted_keys" or r._field == "expired_keys")
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisCollectionsRate(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "total_connections_received" or r._field == "rejected_connections")
		|> derivative(unit: 1s, nonNegative: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetRedisConnectedSlaves(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "connected_slaves")
		|> aggregateWindow(every: 1s, fn: min)
		|> min()
	`, x.Values.Influx.Bucket)
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

func (x *Service) GetRedisHitRate(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "keyspace_hitrate")
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) ToRedisNetworkIO(v string) int {
	switch v {
	case "total_net_output_bytes":
		return 0
	case "total_net_input_bytes":
		return 1
	}
	return 0
}

func (x *Service) GetRedisNetworkIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "redis")
		|> filter(fn: (r) => r._field == "total_net_output_bytes" or r._field == "total_net_input_bytes")
		|> derivative(unit: 1s, nonNegative: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			x.ToRedisNetworkIO(result.Record().Field()),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsUptime(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "uptime")
		|> last()
	`, x.Values.Influx.Bucket)
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

func (x *Service) GetNatsCpu(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`
		import "dict"
		targets = ["http://nats1:8222": 1, "http://nats2:8222": 2, "http://nats3:8222": 3]
		from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "cpu")
		|> map(fn: (r) => ({r with server: dict.get(dict: targets, key: r.server, default: 0)}))
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
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
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`
		import "dict"
		targets = ["http://nats1:8222": 1, "http://nats2:8222": 2, "http://nats3:8222": 3]
		from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "mem")
		|> map(fn: (r) => ({r with server: dict.get(dict: targets, key: r.server, default: 0)}))
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
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
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`
		import "dict"
		targets = ["http://nats1:8222": 1, "http://nats2:8222": 2, "http://nats3:8222": 3]
		from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "connections")
		|> map(fn: (r) => ({r with server: dict.get(dict: targets, key: r.server, default: 0)}))
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
			result.Record().ValueByKey("server"),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsSubscriptions(ctx context.Context) (value interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "subscriptions")
		|> group(columns: ["host"])
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
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

func (x *Service) GetNatsSlowConsumers(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "slow_consumers")
		|> group(columns: ["host"])
		|> aggregateWindow(every: 1s, fn: mean, createEmpty: false)
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().Value(),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsMsgIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`
		import "dict"
		targets = ["http://nats1:8222": 1, "http://nats2:8222": 2, "http://nats3:8222": 3]
		from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "in_msgs" or r._field == "out_msgs")
		|> map(fn: (r) => ({r with server: dict.get(dict: targets, key: r.server, default: 0)}))
		|> derivative(unit: 1s, nonNegative: false)
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().ValueByKey("in_msgs"),
			result.Record().ValueByKey("out_msgs"),
			result.Record().ValueByKey("server"),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}

func (x *Service) GetNatsBytesIO(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`
		import "dict"
		targets = ["http://nats1:8222": 1, "http://nats2:8222": 2, "http://nats3:8222": 3]
		from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r._measurement == "nats")
		|> filter(fn: (r) => r._field == "in_bytes" or r._field == "out_bytes")
		|> map(fn: (r) => ({r with server: dict.get(dict: targets, key: r.server, default: 0)}))
		|> derivative(unit: 1s, nonNegative: false)
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, x.Values.Influx.Bucket)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}

	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, []interface{}{
			result.Record().Time().Format(time.TimeOnly),
			result.Record().ValueByKey("in_bytes"),
			result.Record().ValueByKey("out_bytes"),
			result.Record().ValueByKey("server"),
		})
	}

	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}
