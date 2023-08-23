package schedules_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Keys(t *testing.T) {
	data, err := x.SchedulesService.Keys("b4qxEqajeelm02bkuBmuc")
	assert.NoError(t, err)
	t.Log(data)
}
