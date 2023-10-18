package index

import (
	"context"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/common"
)

type Service struct {
	*common.Inject
}

func (x *Service) Verify(_ context.Context, ts string) (claims passport.Claims, err error) {
	//if claims, err = x.Passport.Verify(ts); err != nil {
	//	return
	//}
	return
}
