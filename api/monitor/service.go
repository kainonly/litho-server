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

func (x *Service) GetQPS(ctx context.Context) (data []interface{}, err error) {
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
  		|> map(fn: (r) => ({r with _value: r._value / 60000000000}))
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
		|> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
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
