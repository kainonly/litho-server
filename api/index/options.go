package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"net/http"
)

type OptionsDto struct {
	Type string `query:"type"`
}

func (x *Controller) Options(_ context.Context, c *app.RequestContext) {
	var dto OptionsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	switch dto.Type {
	case "upload":
		switch x.V.Cloud {
		case "tencent":
			c.JSON(http.StatusOK, M{
				"type": "cos",
				"url": fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
					x.V.TencentCosBucket, x.V.TencentCosRegion,
				),
				"limit": x.V.TencentCosLimit,
			})
			return
		}
	case "collaboration":
		c.JSON(http.StatusOK, M{
			"url":      "https://open.larksuite.com/open-apis/authen/v1/index",
			"redirect": x.V.RedirectUrl,
			"app_id":   x.V.LarkAppId,
		})
		return
	case "generate-secret":
		c.JSON(http.StatusOK, M{
			"id":  help.Random(8),
			"key": help.Random(16),
		})
		return
	case "monitor":
		c.JSON(http.StatusOK, M{
			"enabled": *x.V.Influx.Enabled,
		})
		return
	}
	return
}
