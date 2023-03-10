package model_test

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateAccessLogCollection(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().
				SetTimeField("timestamp").
				SetMetaField("metadata"),
		).
		SetExpireAfterSeconds(15552000)
	if err := db.CreateCollection(ctx, "access_logs", option); err != nil {
		t.Error(err)
	}
}
