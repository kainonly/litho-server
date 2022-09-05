package index_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/common"
	"testing"
)

func TestService_Login(t *testing.T) {
	var err error
	var active common.Active
	active, err = x.IndexService.Login(context.Background(), "weplanx", "pass@VAN1234")
	assert.Nil(t, err)
	t.Log(active)
}
