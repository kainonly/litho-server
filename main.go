package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/spf13/cobra"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/server/openapi"
	"github.com/weplanx/server/xapi"
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
		Use:   "xapi",
		Short: "Start Internal API service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var x *xapi.API
			if x, err = bootstrap.NewXAPI(values); err != nil {
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
			h.Spin()
			return
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Start Open API service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var x *openapi.API
			if x, err = bootstrap.NewOpenAPI(values); err != nil {
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
			h.Spin()
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
			if err = model.SetCategories(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetClusters(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetProjects(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetImessages(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetLogsetLogins(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetLogsetJobs(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetLogsetOperates(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetLogsetImessages(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetPictures(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetQueues(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetSchedules(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetUsers(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetUsers(ctx, x.Db); err != nil {
				return
			}
			if err = model.SetWorkflows(ctx, x.Db); err != nil {
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
		"email", "u", "",
		"User's email <Must be email>",
	)
	userCmd.PersistentFlags().StringVarP(&password,
		"password", "p", "",
		"User's password <between 8-20>",
	)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "values",
		Short: "Display the dynamic values of server distribution Kv",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var x *api.API
			if x, err = bootstrap.NewAPI(values); err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if _, err = x.Initialize(ctx); err != nil {
				return
			}
			time.Sleep(time.Second)
			var output []byte
			if output, err = json.MarshalIndent(x.V, "", "    "); err != nil {
				return
			}
			fmt.Println(string(output))
			return
		},
	})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
