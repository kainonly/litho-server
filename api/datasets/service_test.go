package datasets_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_List(t *testing.T) {
	ctx := context.TODO()
	data, err := x.DatasetsService.List(ctx)
	assert.NoError(t, err)
	t.Log(data)
}
