package datasets_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Lists(t *testing.T) {
	ctx := context.TODO()
	data, err := x.DatasetsX.Lists(ctx, "")
	assert.NoError(t, err)
	t.Log(data)
}
