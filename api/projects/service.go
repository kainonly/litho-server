package projects

import (
	"bytes"
	"context"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/server/api/clusters"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"text/template"
)

type Service struct {
	*common.Inject
	ClustersX *clusters.Service
}

func (x *Service) GetTenants(ctx context.Context, id primitive.ObjectID) (result M, err error) {
	var project model.Project
	if err = x.Db.Collection("projects").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&project); err != nil {
		return
	}

	result = M{}
	if project.Nats != nil {
		var b []byte
		if b, err = x.Cipher.Decode(project.Nats.Seed); err != nil {
			return
		}
		result["nkey"] = string(b)
	}

	return
}

func (x *Service) DeployNats(ctx context.Context, id primitive.ObjectID) (err error) {
	var project model.Project
	if err = x.Db.Collection("projects").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&project); err != nil {
		return
	}
	if project.Cluster == nil {
		return
	}
	var cluster model.Cluster
	if cluster, err = x.ClustersX.Get(ctx, *project.Cluster); err != nil {
		return
	}
	var accounts []NatsAccount
	if cluster.Admin {
		accounts = append(accounts, NatsAccount{
			Name:  "weplanx",
			Users: []NatsUser{{Nkey: x.V.Nats.Pub}},
		})
	}
	if err = x.MakeNatsAccount(ctx, project); err != nil {
		return
	}
	if err = x.SyncNatsAccounts(ctx, project, &accounts); err != nil {
		return
	}

	var tmpl *template.Template
	if tmpl, err = template.ParseFiles("./templates/account.tpl"); err != nil {
		return
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, accounts); err != nil {
		return
	}

	var kube *kubernetes.Clientset
	if kube, err = x.ClustersX.GetClient(cluster); err != nil {
		return
	}
	secret := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: "nats-system",
			Name:      "include",
		},
		Data: map[string][]byte{"accounts.conf": buf.Bytes()},
		Type: "Opaque",
	}
	if _, err = kube.CoreV1().
		Secrets("nats-system").
		Update(ctx, secret, meta.UpdateOptions{}); err != nil {
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
