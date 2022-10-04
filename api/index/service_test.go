package index_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Login(t *testing.T) {
	ctx := context.TODO()
	ts, err := x.IndexService.Login(ctx, "weplanx", "pass@VAN1234")
	assert.Nil(t, err)
	t.Log(ts)
}
