package clusters

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type GetNodesDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetNodes(ctx context.Context, c *app.RequestContext) {
	var dto GetNodesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ClustersX.GetNodes(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) GetNodes(ctx context.Context, id primitive.ObjectID) (result []interface{}, err error) {
	var data model.Cluster
	if data, err = x.Get(ctx, id); err != nil {
		return
	}
	var kube *kubernetes.Clientset
	if kube, err = x.GetClient(data); err != nil {
		return
	}
	var nodes *v1.NodeList
	if nodes, err = kube.CoreV1().Nodes().List(ctx, meta.ListOptions{}); err != nil {
		return
	}
	for _, v := range nodes.Items {
		result = append(result, M{
			"name":         v.GetName(),
			"create":       v.GetCreationTimestamp(),
			"hostname":     v.Annotations["k3s.io/hostname"],
			"ip":           v.Annotations["k3s.io/internal-ip"],
			"version":      v.Status.NodeInfo.KubeletVersion,
			"cpu":          v.Status.Allocatable.Cpu().Value(),
			"mem":          v.Status.Allocatable.Memory().Value(),
			"storage":      v.Status.Allocatable.StorageEphemeral().Value(),
			"os":           v.Status.NodeInfo.OSImage,
			"architecture": v.Status.NodeInfo.Architecture,
		})
	}
	return
}
