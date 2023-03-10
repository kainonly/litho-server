package model_test

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateLoginLogsCollection(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().
				SetTimeField("timestamp").
				SetMetaField("metadata"),
		)
	if err := db.CreateCollection(ctx, "login_logs", option); err != nil {
		t.Error(err)
	}
}
