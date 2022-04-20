package system

import (
	"api/common"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/transfer"
	"io"
	"io/ioutil"
	"time"
)

type Middleware struct {
	*Service
	Passport *passport.Passport
	Transfer *transfer.Transfer
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
		ok, err := x.VerifySession(ctx, uid, claims["jti"].(string))
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
		if err = x.RenewSession(ctx, uid); err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set(common.TokenClaimsKey, claims)
		c.Next()
	}
}
