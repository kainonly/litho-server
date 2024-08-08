package projects

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeployNatsDto struct {
	Id primitive.ObjectID `json:"id" vd:"required"`
}

func (x *Controller) DeployNats(ctx context.Context, c *app.RequestContext) {
	var dto DeployNatsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.ProjectsX.DeployNats(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
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
