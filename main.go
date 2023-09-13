package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"os"
	"time"
)

func main() {
	var config string
	var values *common.Values
	rootCmd := &cobra.Command{
		Use:               "weplanx",
		Short:             "API service, based on Hertz's project",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if values, err = bootstrap.LoadStaticValues(config); err != nil {
				return
			}
			return
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config,
		"config", "c", "config/default.values.yml",
		"The default configuration file of weplanx server values",
	)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "setup",
		Short: "Initialize weplanx server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var x *api.API
			if x, err = bootstrap.NewAPI(values); err != nil {
				return
			}
			if err = x.Values.Service.Update(x.V.Extra); err != nil {
				return
			}
			return
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "api",
		Short: "Start API service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var x *api.API
			if x, err = bootstrap.NewAPI(values); err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			var h *server.Hertz
			if h, err = x.Initialize(ctx); err != nil {
				return
			}
			if err = x.Routes(h); err != nil {
				return
			}
			defer bootstrap.ProviderOpenTelemetry(values).
				Shutdown(ctx)

			h.Spin()
			return
		},
	})
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
}
