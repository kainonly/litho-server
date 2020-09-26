package cors

import (
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

type Option struct {
	Origin        []string
	Method        []string
	AllowHeader   []string
	ExposedHeader []string
	MaxAge        int64
	Credentials   bool
}

func Cors(option Option) iris.Handler {
	origin := strings.Join(option.Origin, ",")
	method := strings.Join(option.Method, ",")
	allowHeader := strings.Join(option.AllowHeader, ",")
	exposedHeader := strings.Join(option.ExposedHeader, ",")
	maxAge := strconv.FormatInt(option.MaxAge, 10)
	return func(ctx iris.Context) {
		ctx.Header("access-control-allow-origin", origin)
		ctx.Header("Access-Control-Allow-Methods", method)
		ctx.Header("Access-Control-Allow-Headers", allowHeader)
		ctx.Header("Access-Control-Expose-Headers", exposedHeader)
		ctx.Header("Access-Control-Max-Age", maxAge)
		if option.Credentials {
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		ctx.Next()
	}
}
