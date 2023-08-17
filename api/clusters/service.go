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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Service struct {
	*common.Inject
}

func (x *Service) GetNodes(ctx context.Context, id primitive.ObjectID) (data []interface{}, err error) {
	var cluster model.Cluster
	if err = x.Db.Collection("clusters").FindOne(ctx,
		bson.M{"_id": id},
	).Decode(&cluster); err != nil {
		return
	}

	var kube *kubernetes.Clientset
	if kube, err = x.GetKube(cluster.Config); err != nil {
		return
	}
	var nodes *v1.NodeList
	if nodes, err = kube.CoreV1().Nodes().List(ctx, meta.ListOptions{}); err != nil {
		return
	}
	for _, v := range nodes.Items {
		data = append(data, M{
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

type Kubeconfig struct {
	Host     string `json:"host"`
	CAData   string `json:"ca_data"`
	CertData string `json:"cert_data"`
	KeyData  string `json:"key_data"`
}

func (x *Service) GetKube(ciphertext string) (kube *kubernetes.Clientset, err error) {
	var b []byte
	if b, err = x.Cipher.Decode(ciphertext); err != nil {
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
	if kube, err = kubernetes.NewForConfig(&rest.Config{
		Host: config.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   ca,
			CertData: cert,
			KeyData:  key,
		},
	}); err != nil {
		return
	}
	return
}
