package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogsetLogined struct {
	Timestamp time.Time             `bson:"timestamp" json:"timestamp"`
	Metadata  LogsetLoginedMetadata `bson:"metadata" json:"metadata"`
	UserAgent string                `bson:"user_agent" json:"user_agent"`
	Detail    interface{}           `bson:"detail" json:"detail"`
}

type LogsetLoginedMetadata struct {
	UserID   primitive.ObjectID `bson:"user_id" json:"-"`
	ClientIP string             `bson:"client_ip" json:"client_ip"`
	Version  string             `bson:"version" json:"version"`
	Source   string             `bson:"source" json:"source" json:"source"`
}

func (x *LogsetLogined) SetDetail(v interface{}) {
	x.Detail = v
}

func (x *LogsetLogined) SetVersion(v string) {
	x.Metadata.Version = v
}

func NewLogsetLogined(uid primitive.ObjectID, ip string, source string, useragent string) *LogsetLogined {
	return &LogsetLogined{
		Timestamp: time.Now(),
		Metadata: LogsetLoginedMetadata{
			UserID:   uid,
			ClientIP: ip,
			Source:   source,
		},
		UserAgent: useragent,
	}
}

func SetLogsetLogined(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "logset_logined"}); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField("timestamp").
					SetMetaField("metadata"),
			)
		if err = db.CreateCollection(ctx, "logset_logined", option); err != nil {
			return
		}
	}
	return
}

func SetLogsetJobs(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	filter := bson.M{"name": bson.M{"$in": bson.A{"logset_jobs", "logset_jobs_fail"}}}
	if ns, err = db.ListCollectionNames(ctx, filter); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField("timestamp").
					SetMetaField("metadata"),
			).
			SetExpireAfterSeconds(31536000)
		if err = db.CreateCollection(ctx, "logset_jobs", option); err != nil {
			return
		}
		if err = db.CreateCollection(ctx, "logset_jobs_fail", option); err != nil {
			return
		}
	}
	return
}

func SetLogsetOperates(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	filter := bson.M{"name": bson.M{"$in": bson.A{"logset_operates", "logset_operates_fail"}}}
	if ns, err = db.ListCollectionNames(ctx, filter); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField("timestamp").
					SetMetaField("metadata"),
			).
			SetExpireAfterSeconds(31536000)
		if err = db.CreateCollection(ctx, "logset_operates", option); err != nil {
			return
		}
		if err = db.CreateCollection(ctx, "logset_operates_fail", option); err != nil {
			return
		}
	}
	return
}
