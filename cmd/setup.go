package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
)

var Setup = &cobra.Command{
	Use:   "setup",
	Short: "Setup weplanx dynamic configuration",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		values := ctx.Value("values").(*common.Values)

		var x *api.API
		if x, err = bootstrap.NewAPI(values); err != nil {
			return
		}
		if err = x.Values.Service.Update(x.V.Extra); err != nil {
			return
		}
		return
	},
}
