package res

import (
	"github.com/kataras/iris/v12"
)

func Ok() interface{} {
	return iris.Map{
		"error": false,
		"msg":   "ok",
	}
}

func Data(value interface{}) interface{} {
	return iris.Map{
		"error": false,
		"data":  value,
	}
}

func Error(msg interface{}) interface{} {
	if val, ok := msg.(error); ok {
		msg = val.Error()
	}
	return iris.Map{
		"error": true,
		"msg":   msg,
	}
}
