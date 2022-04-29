package vars

import (
	"api/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Controller struct {
	Service *Service
}

// Gets 获取指定变量
func (x *Controller) Gets(c *gin.Context) interface{} {
	var query struct {
		Keys []string `form:"keys" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	values, err := x.Service.Gets(ctx, query.Keys)
	if err != nil {
		return err
	}
	for k, v := range values {
		if common.SecretKey(k) {
			if v == "" || v == nil {
				values[k] = "未设置"
			} else {
				values[k] = "已设置"

			}
		}
	}
	return values
}

// Get 获取变量
func (x *Controller) Get(c *gin.Context) interface{} {
	var uri struct {
		Key string `uri:"key" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	value, err := x.Service.Get(ctx, uri.Key)
	if err != nil {
		return err
	}
	if common.SecretKey(uri.Key) {
		value = "已设置"
	}
	return value
}

// Set 设置变量
func (x *Controller) Set(c *gin.Context) interface{} {
	var uri struct {
		Key string `uri:"key" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	var body struct {
		Value interface{} `json:"value"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if err := x.Service.Set(ctx, uri.Key, body.Value); err != nil {
		return err
	}
	return nil
}

func (x *Controller) Options(c *gin.Context) interface{} {
	var query struct {
		Type string `form:"type" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	switch query.Type {
	case "upload":
		platform, err := x.Service.Get(ctx, "cloud_platform")
		if err != nil {
			return err
		}
		switch platform {
		case "tencent":
			v, err := x.Service.Gets(ctx, []string{
				"tencent_cos_bucket",
				"tencent_cos_region",
				"tencent_cos_limit",
			})
			if err != nil {
				return err
			}
			limit, _ := strconv.Atoi(v["tencent_cos_limit"].(string))
			return gin.H{
				"type": "cos",
				"url": fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
					v["tencent_cos_bucket"], v["tencent_cos_region"],
				),
				"limit": limit,
			}
		}
	}
	return nil
}
