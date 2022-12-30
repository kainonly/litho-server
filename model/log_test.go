package model_test

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateCollectionForLoginLogs(t *testing.T) {
	ctx := context.TODO()
	if err := db.CreateCollection(ctx, "login_logs", options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().
				SetMetaField("metadata").
				SetTimeField("timestamp"),
		).SetExpireAfterSeconds(15552000)); err != nil {
		t.Error(err)
	}
}
