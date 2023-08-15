package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateLogsetJobsCollection(t *testing.T) {
	ctx := context.TODO()
	err := x.Db.Collection("logset_jobs").Drop(ctx)
	assert.NoError(t, err)
	err = x.Db.Collection("logset_jobs_fail").Drop(ctx)
	assert.NoError(t, err)
	option := options.CreateCollection().
		SetTimeSeriesOptions(
			options.TimeSeries().
				SetTimeField("timestamp").
				SetMetaField("metadata"),
		).
		SetExpireAfterSeconds(31536000)
	err = x.Db.CreateCollection(ctx, "logset_jobs", option)
	assert.NoError(t, err)
	err = x.Db.CreateCollection(ctx, "logset_jobs_fail", option)
	assert.NoError(t, err)
}
