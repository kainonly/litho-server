package xapi

import (
	"context"
	"github.com/bytedance/go-tagexpr/v2/binding"
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	"github.com/weplanx/go/help"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/xapi/emqx"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var Provides = wire.NewSet(
	emqx.Provides,
)

type API struct {
	*common.Inject

	Hertz       *server.Hertz
	Emqx        *emqx.Controller
	EmqxService *emqx.Service
}

func (x *API) Routes(h *server.Hertz) (err error) {
	_emqx := h.Group("emqx")
	{
		_emqx.POST("auth", x.Emqx.Auth)
		_emqx.POST("acl", x.Emqx.Acl)
	}
	return
}

func (x *API) ErrHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		if err.IsType(errors.ErrorTypePublic) {
			statusCode := http.StatusBadRequest
			result := utils.H{"message": err.Error()}
			if meta, ok := err.Meta.(map[string]interface{}); ok {
				if meta["statusCode"] != nil {
					statusCode = meta["statusCode"].(int)
				}
				if meta["code"] != nil {
					result["code"] = meta["code"]
				}
			}
			c.JSON(statusCode, result)
			return
		}

		switch e := err.Err.(type) {
		case decoder.SyntaxError:
			c.JSON(http.StatusBadRequest, utils.H{
				"code":    0,
				"message": e.Description(),
			})
			break
		case *binding.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"code":    0,
				"message": e.Error(),
			})
			break
		case *validator.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"code":    0,
				"message": e.Error(),
			})
			break
		case mongo.ServerError:
			c.JSON(http.StatusInternalServerError, utils.H{
				"code":    0,
				"message": e.Error(),
			})
			break
		default:
			if !x.V.IsRelease() {
				c.JSON(http.StatusInternalServerError, utils.H{
					"code":    0,
					"message": e.Error(),
				})
				break
			}
			logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
	}
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	help.RegValidate()

	h = x.Hertz
	h.Use(x.ErrHandler())

	return
}
