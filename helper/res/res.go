package res

import (
	"github.com/kataras/iris/v12"
)

func Ok() interface{} {
	return iris.Map{
		"error": 0,
		"msg":   "ok",
	}
}

func Data(value interface{}) interface{} {
	return iris.Map{
		"error": 0,
		"data":  value,
	}
}

func Error(msg interface{}) interface{} {
	if val, ok := msg.(error); ok {
		msg = val.Error()
	}
	return iris.Map{
		"error": 1,
		"msg":   msg,
	}
}
