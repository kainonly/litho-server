package workflows

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/common"
)

type Service struct {
	*common.Inject
}

func (x *Service) Event() (err error) {
	subj := x.V.NameX(".", "workflows")
	queue := x.V.Name("workflows")
	if _, err = x.JetStream.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		//defer cancel()
		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}
		fmt.Println(dto)
		switch dto.Action {
		case "create":
			break
		case "bulk_create":
			break
		case "update":
			break
		case "update_by_id":
			break
		case "replace":
			break
		case "delete":
			break
		case "bulk_delete":
			break
		}
	}, nats.ManualAck()); err != nil {
		return
	}
	return
}
