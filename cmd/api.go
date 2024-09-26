package cmd

import (
	"github.com/spf13/cobra"
)

var API = &cobra.Command{
	Use:   "api",
	Short: "Start API service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//ctx := cmd.Context()
		//values := ctx.Value("values").(*common.Values)
		//if values.Address == "" {
		//	values.Address = ":3000"
		//}
		//
		//var x *api.API
		//if x, err = bootstrap.NewAPI(values); err != nil {
		//	return
		//}
		//var h *server.Hertz
		//if h, err = x.Initialize(ctx); err != nil {
		//	return
		//}
		//if err = x.Routes(h); err != nil {
		//	return
		//}
		//
		//if *values.Otlp.Enabled {
		//	defer bootstrap.ProviderOpenTelemetry(values).
		//		Shutdown(ctx)
		//}
		//
		//h.Spin()
		return
	},
}
