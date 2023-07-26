package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateLogsetLoginCollection(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().
				SetTimeField("timestamp").
				SetMetaField("metadata"),
		)
	err := x.Db.CreateCollection(ctx, "logset_logined", option)
	assert.NoError(t, err)
}
