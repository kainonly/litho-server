package clusters

import (
	"context"
	"encoding/base64"
	"github.com/bytedance/sonic"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

type Service struct {
	*common.Inject
}

func (x *Service) Get(ctx context.Context, id primitive.ObjectID) (data model.Cluster, err error) {
	if err = x.Db.Collection("clusters").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	return
}

var kubes = sync.Map{}

type Kubeconfig struct {
	Host     string `json:"host"`
	CAData   string `json:"ca_data"`
	CertData string `json:"cert_data"`
	KeyData  string `json:"key_data"`
}

func (x *Service) GetClient(data model.Cluster) (client *kubernetes.Clientset, err error) {
	id := data.ID.Hex()
	if i, ok := kubes.Load(id); ok {
		client = i.(*kubernetes.Clientset)
		return
	}
	var b []byte
	if b, err = x.Cipher.Decode(data.Config); err != nil {
		return
	}
	var config Kubeconfig
	if err = sonic.Unmarshal(b, &config); err != nil {
		return
	}
	var ca []byte
	if ca, err = base64.StdEncoding.DecodeString(config.CAData); err != nil {
		return
	}
	var cert []byte
	if cert, err = base64.StdEncoding.DecodeString(config.CertData); err != nil {
		return
	}
	var key []byte
	if key, err = base64.StdEncoding.DecodeString(config.KeyData); err != nil {
		return
	}
	if client, err = kubernetes.NewForConfig(&rest.Config{
		Host: config.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   ca,
			CertData: cert,
			KeyData:  key,
		},
	}); err != nil {
		return
	}
	kubes.Store(id, client)
	return
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
