package app

import (
	"bytes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/transfer"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"server/app/system"
	"server/common"
	"time"
)

type Middleware struct {
	*common.Inject
	Values   *common.Values
	Transfer *transfer.Transfer
	System   *system.Service
	Passport *passport.Passport
}

func (x *Middleware) Global() *gin.Engine {
	r := gin.New()
	if os.Getenv("GIN_MODE") == "release" {
		logger, _ := zap.NewProduction()
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	} else {
		r.Use(gin.Logger())
	}
	r.SetTrustedProxies(x.Values.TrustedProxies)
	r.Use(requestid.New())
	r.Use(gin.CustomRecovery(catchError))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     x.Values.Cors.AllowOrigins,
		AllowMethods:     x.Values.Cors.AllowMethods,
		AllowHeaders:     x.Values.Cors.AllowHeaders,
		ExposeHeaders:    x.Values.Cors.ExposeHeaders,
		AllowCredentials: x.Values.Cors.AllowCredentials,
		MaxAge:           time.Duration(x.Values.Cors.MaxAge) * time.Second,
	}))
	r.Use(x.RequestLogging())
	helper.ExtendValidation()
	return r
}

type Response struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w Response) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (x *Middleware) RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &Response{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = resp
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		c.Next()
		go func() {
			tags := map[string]string{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"request_id": c.Request.Header.Get("x-request-id"),
			}
			if v, ok := c.Get(common.TokenClaimsKey); ok {
				claims := v.(jwt.MapClaims)
				tags["jti"] = claims["jti"].(string)
				tags["uid"] = claims["context"].(map[string]interface{})["uid"].(string)
			}
			payload, err := transfer.NewPayload(transfer.InfluxDto{
				Measurement: "request",
				Tags:        tags,
				Fields: map[string]interface{}{
					"query":    c.Request.URL.Query(),
					"headers":  c.Request.Header,
					"body":     body,
					"response": resp.body.String(),
				},
				Time: time.Now(),
			})
			if err != nil {
				return
			}
			if err = x.Transfer.Publish(c.Request.Context(), "request", payload); err != nil {
				return
			}
		}()
	}
}

func (x *Middleware) AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ts, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.AuthExpired.Error(),
			})
			return
		}
		claims, err := x.Passport.Verify(ts)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.AuthExpired.Error(),
			})
			return
		}
		ctx := c.Request.Context()
		uid := claims["context"].(map[string]interface{})["uid"].(string)
		ok, err := x.System.VerifySession(ctx, uid, claims["jti"].(string))
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_CONFLICT",
				"message": common.AuthExpired.Error(),
			})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_CONFLICT",
				"message": common.AuthConflict.Error(),
			})
			return
		}
		if err = x.System.RenewSession(ctx, uid); err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set(common.TokenClaimsKey, claims)
		c.Next()
	}
}

func catchError(c *gin.Context, err interface{}) {
	c.AbortWithStatusJSON(500, gin.H{
		"message": err,
	})
}
