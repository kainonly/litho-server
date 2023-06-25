package admin

import (
	"github.com/weplanx/server/admin/index"
	"github.com/weplanx/server/common"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	fx.Invoke(Routes),
)

type Inject struct {
	common.Inject
	Index *index.Controller
}

func Routes(i Inject) (err error) {

	i.Hertz.GET("", i.Index.Ping)

	return
}

//func (x *API) AccessLogs() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		now := time.Now()
//		c.Next(ctx)
//		method := string(c.Request.Header.Method())
//		if method == "GET" {
//			return
//		}
//		var userId string
//		if value, ok := c.Get("identity"); ok {
//			claims := value.(passport.Claims)
//			userId = claims.UserId
//		}
//		x.Transfer.Publish(context.Background(), "access", transfer.Payload{
//			Timestamp: now,
//			Data: map[string]interface{}{
//				"metadata": map[string]interface{}{
//					"method":    method,
//					"path":      string(c.Request.Path()),
//					"user_id":   userId,
//					"client_ip": c.ClientIP(),
//				},
//				"status":     c.Response.StatusCode(),
//				"user_agent": string(c.Request.Header.UserAgent()),
//			},
//			Format: map[string]interface{}{
//				"metadata.user_id": "oid",
//			},
//		})
//	}
//}

//func (x *API) ErrHandler() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		c.Next(ctx)
//		err := c.Errors.Last()
//		if err == nil {
//			return
//		}
//
//		if err.IsType(errors.ErrorTypePublic) {
//			statusCode := http.StatusBadRequest
//			result := utils.H{"message": err.Error()}
//			if meta, ok := err.Meta.(map[string]interface{}); ok {
//				if meta["statusCode"] != nil {
//					statusCode = meta["statusCode"].(int)
//				}
//				if meta["code"] != nil {
//					result["code"] = meta["code"]
//				}
//			}
//			c.JSON(statusCode, result)
//			return
//		}
//
//		switch e := err.Err.(type) {
//		case decoder.SyntaxError:
//			c.JSON(http.StatusBadRequest, utils.H{
//				"code":    0,
//				"message": e.Description(),
//			})
//			break
//		case *binding.Error:
//			c.JSON(http.StatusBadRequest, utils.H{
//				"code":    0,
//				"message": e.Error(),
//			})
//			break
//		case *validator.Error:
//			c.JSON(http.StatusBadRequest, utils.H{
//				"code":    0,
//				"message": e.Error(),
//			})
//			break
//		default:
//			logger.Error(err)
//			c.Status(http.StatusInternalServerError)
//		}
//	}
//}

//func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
//	h = x.Hertz
//
//	h.Use(x.ErrHandler())
//
//	return
//}
