package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"os"
	"time"
)

var values *common.Values

func main() {
	var config string
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
	rootCmd.AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Start OpenAPI service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return
		},
	})
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			if err = model.SetupCategory(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupCluster(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupProject(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupImessage(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupLogsetLogined(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupPicture(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupQueue(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupSchedule(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupUser(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupUser(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetupWorkflow(ctx, x.Db); err != nil {
				return
			}
			return
		},
	})
	var email string
	var password string
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Create an email account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var x *api.API
			if x, err = bootstrap.NewAPI(values); err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if _, err = x.Db.Collection("users").InsertOne(
				ctx,
				model.NewUser(email, password),
			); err != nil {
				return
			}
			return
		},
	}
	userCmd.PersistentFlags().StringVarP(&email,
		"email", "e", "",
		"User's email",
	)
	userCmd.PersistentFlags().StringVarP(&password,
		"password", "p", "",
		"User's password <8~20>",
	)
	rootCmd.AddCommand(userCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
}
