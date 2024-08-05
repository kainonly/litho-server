package clusters

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
)

type GetInfoDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetInfo(ctx context.Context, c *app.RequestContext) {
	var dto GetInfoDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ClustersX.GetInfo(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) GetInfo(ctx context.Context, id primitive.ObjectID) (result M, err error) {
	var data model.Cluster
	if data, err = x.Get(ctx, id); err != nil {
		return
	}
	var kube *kubernetes.Clientset
	if kube, err = x.GetClient(data); err != nil {
		return
	}
	var info *version.Info
	if info, err = kube.ServerVersion(); err != nil {
		return
	}
	var nodes *v1.NodeList
	if nodes, err = kube.CoreV1().Nodes().List(ctx, meta.ListOptions{}); err != nil {
		return
	}

	cpu := int64(0)
	mem := int64(0)
	storage := int64(0)
	for _, v := range nodes.Items {
		cpu += v.Status.Allocatable.Cpu().Value()
		mem += v.Status.Allocatable.Memory().Value()
		storage += v.Status.Allocatable.StorageEphemeral().Value()
	}

	result = M{
		"version": info.String(),
		"nodes":   len(nodes.Items),
		"cpu":     cpu,
		"mem":     mem,
		"storage": storage,
	}

	return
}
