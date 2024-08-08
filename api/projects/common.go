package projects

import (
	"context"
	"github.com/google/wire"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/server/api/clusters"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	ProjectsX *Service
}

type Service struct {
	*common.Inject
	ClustersX *clusters.Service
}

type M = map[string]interface{}

func (x *Service) Get(ctx context.Context, id primitive.ObjectID) (data model.Project, err error) {
	if err = x.Db.Collection("projects").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	return
}

type NatsAccount struct {
	Name  string
	Users []NatsUser
}

type NatsUser struct {
	Nkey string
}

func (x *Service) MakeNatsAccount(ctx context.Context, project model.Project) (err error) {
	var user nkeys.KeyPair
	if user, err = nkeys.CreateUser(); err != nil {
		return
	}
	if _, err = user.Sign([]byte(project.Namespace)); err != nil {
		return
	}
	var seed []byte
	if seed, err = user.Seed(); err != nil {
		return
	}
	var pub string
	if pub, err = user.PublicKey(); err != nil {
		return
	}
	var xSeed string
	if xSeed, err = x.Cipher.Encode(seed); err != nil {
		return
	}
	var xPub string
	if xPub, err = x.Cipher.Encode([]byte(pub)); err != nil {
		return
	}
	if _, err = x.Db.Collection("projects").UpdateByID(ctx, project.ID, bson.M{
		"$set": bson.M{
			"nats": model.ProjectNats{
				Seed: xSeed,
				Pub:  xPub,
			},
		},
	}); err != nil {
		return
	}
	return
}

func (x *Service) SyncNatsAccounts(ctx context.Context, project model.Project, accounts *[]NatsAccount) (err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("projects").
		Find(ctx, bson.M{
			"cluster": *project.Cluster,
			"nats":    bson.M{"$exists": 1},
		}); err != nil {
		return
	}
	for cursor.Next(ctx) {
		var data model.Project
		if err = cursor.Decode(&data); err != nil {
			return
		}
		var users []NatsUser
		var pub []byte
		if pub, err = x.Cipher.Decode(data.Nats.Pub); err != nil {
			return
		}
		users = append(users, NatsUser{Nkey: string(pub)})
		*accounts = append(*accounts, NatsAccount{
			Name:  data.Namespace,
			Users: users,
		})
	}
	return
}
