package clusters

import (
	"context"
	"encoding/base64"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	ClustersX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]interface{}

var kubes = sync.Map{}

type Kubeconfig struct {
	Host     string `json:"host"`
	CAData   string `json:"ca_data"`
	CertData string `json:"cert_data"`
	KeyData  string `json:"key_data"`
}

func (x *Service) Get(ctx context.Context, id primitive.ObjectID) (data model.Cluster, err error) {
	if err = x.Db.Collection("clusters").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	return
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
