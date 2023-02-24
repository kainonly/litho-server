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

func (x *Service) GetCgoCalls(ctx context.Context) (data []interface{}, err error) {
	queryAPI := x.Influx.QueryAPI(x.Values.Influx.Org)
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -15m, stop: now())
		|> filter(fn: (r) => r["_measurement"] == "prometheus")
		|> filter(fn: (r) => r["_field"] == "process.runtime.go.cgo.calls")
		|> filter(fn: (r) => r["service.name"] == "%s")
		|> group(columns: ["service.name"], mode: "by")
		|> aggregateWindow(every: 5m, fn: mean, createEmpty: false)
		|> yield(name: "mean")
	`, x.Values.Influx.Bucket, x.Values.Namespace)
	var result *api.QueryTableResult
	if result, err = queryAPI.Query(ctx, query); err != nil {
		return
	}
	data = make([]interface{}, 0)
	for result.Next() {
		data = append(data, map[string]interface{}{
			"timestamp": result.Record().Time().Format(time.DateTime),
			"value":     result.Record().Value(),
		})
	}
	if result.Err() != nil {
		hlog.Error(result.Err())
	}
	return
}
