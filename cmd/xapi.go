package cmd

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/cobra"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/xapi"
)

var XAPI = &cobra.Command{
	Use:   "xapi",
	Short: "Start Internal API service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		values := ctx.Value("values").(*common.Values)
		if values.Address == "" {
			values.Address = ":6000"
		}

		var x *xapi.API
		if x, err = bootstrap.NewXAPI(values); err != nil {
			return
		}
		var h *server.Hertz
		if h, err = x.Initialize(ctx); err != nil {
			return
		}
		if err = x.Routes(h); err != nil {
			return
		}
		h.Spin()
		return
	},
}
