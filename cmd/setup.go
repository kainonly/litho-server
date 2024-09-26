package cmd

import (
	"github.com/spf13/cobra"
)

var Setup = &cobra.Command{
	Use:   "setup",
	Short: "Setup weplanx dynamic configuration",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//ctx := cmd.Context()
		//values := ctx.Value("values").(*common.Values)

		//var x *api.API
		//if x, err = bootstrap.NewAPI(values); err != nil {
		//	return
		//}
		return
	},
}
