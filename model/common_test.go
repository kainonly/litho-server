package model_test

import (
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type Inject struct {
	V   *common.Values
	Mgo *mongo.Client
	Db  *mongo.Database
}

var x *Inject

func TestMain(m *testing.M) {}
