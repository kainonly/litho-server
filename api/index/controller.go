package index

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (x *Controller) Index(_ context.Context, _ *empty.Empty) (*IndexReply, error) {
	return nil, nil
}
