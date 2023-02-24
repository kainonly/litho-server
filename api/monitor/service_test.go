package monitor_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetCgoCalls(t *testing.T) {
	data, err := x.MonitorX.GetCgoCalls(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}
