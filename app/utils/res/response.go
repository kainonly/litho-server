package res

import "github.com/kataras/iris/v12"

func Ok() interface{} {
	return iris.Map{
		"error": 0,
		"msg":   "ok",
	}
}

func Result(data interface{}) interface{} {
	return iris.Map{
		"error": 0,
		"data":  data,
	}
}

func Error(msg string) interface{} {
	return iris.Map{
		"error": 1,
		"msg":   msg,
	}
}
