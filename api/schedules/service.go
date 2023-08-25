package schedules

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/api/clusters"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/workflow"
	"github.com/weplanx/workflow/typ"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
	Clusters *clusters.Service
}

func (x *Service) Ping(id string) (r bool, err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Workflow.NewSchedule(id); err != nil {
		return
	}
	return schedule.Ping()
}

func (x *Service) Deploy(ctx context.Context, id primitive.ObjectID) (err error) {
	var data model.Schedule
	if err = x.Db.Collection("schedules").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	var kube *kubernetes.Clientset
	if kube, err = x.Clusters.GetKube(ctx, data.ClusterId); err != nil {
		return
	}
	container := core.Container{
		Name:            data.Name,
		Image:           data.Image,
		ImagePullPolicy: core.PullAlways,
		Env: []core.EnvVar{
			{Name: "MODE", Value: "release"},
			{Name: "NAMESPACE", Value: x.V.Namespace},
			{Name: "NODE", Value: id.Hex()},
			{Name: "NATS_HOSTS", Value: strings.Join(x.V.Nats.Hosts, ",")},
			{Name: "NATS_NKEY", Value: x.V.Nats.Nkey},
		},
	}
	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name: data.Name,
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"app": data.Name},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"app": data.Name},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{container},
				},
			},
		},
	}
	if _, err = kube.AppsV1().
		Deployments("default").
		Create(ctx, deployment, meta.CreateOptions{}); err != nil {
		return
	}
	return
}

func (x *Service) Undeploy(ctx context.Context, id primitive.ObjectID) (err error) {
	var r bool
	if r, err = x.Ping(id.Hex()); err != nil {
		if err == nats.ErrBucketNotFound {
			return nil
		}
		return
	}
	if !r {
		return
	}

	var data model.Schedule
	if err = x.Db.Collection("schedules").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	var kube *kubernetes.Clientset
	if kube, err = x.Clusters.GetKube(ctx, data.ClusterId); err != nil {
		return
	}
	return kube.AppsV1().
		Deployments("default").
		Delete(ctx, data.Name, meta.DeleteOptions{})
}

func (x *Service) Keys(id string) (keys []string, err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Workflow.NewSchedule(id); err != nil {
		return
	}
	return schedule.Lists()
}

func (x *Service) Set(id string, key string, option typ.ScheduleOption) (err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Workflow.NewSchedule(id); err != nil {
		return
	}
	return schedule.Set(key, option)
}

func (x *Service) Get(id string, key string) (r typ.ScheduleOption, err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Workflow.NewSchedule(id); err != nil {
		return
	}
	return schedule.Get(key)
}

func (x *Service) Revoke(id string, key string) (err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Workflow.NewSchedule(id); err != nil {
		return
	}
	return schedule.Remove(key)
}

func (x *Service) Event() (err error) {
	subj := x.V.NameX(".", "events", "schedules")
	queue := x.V.Name("events", "schedules")
	if _, err = x.JetStream.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}

		switch dto.Action {
		case "create":
			id, _ := primitive.ObjectIDFromHex(dto.Result.(M)["InsertedID"].(string))
			if err = x.Deploy(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
