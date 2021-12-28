package mock

import (
	"api/bootstrap"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

var (
	Client *mongo.Client
	Db     *mongo.Database
)

func TestMain(m *testing.M) {
	os.Chdir("../")
	values, err := bootstrap.SetValues()
	if err != nil {
		panic(err)
	}
	if Client, err = bootstrap.UseMongoDB(values); err != nil {
		panic(err)
	}
	Db = bootstrap.UseDatabase(Client, values)
	os.Exit(m.Run())
}
