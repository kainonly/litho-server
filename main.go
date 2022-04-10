package main

import (
	"api/bootstrap"
	"api/common"
	"context"
	"github.com/fatih/color"
	"github.com/gin-contrib/pprof"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func main() {
	var config string
	rootCmd := &cobra.Command{
		Version: "v0.0.0",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config,
		"config", "c", "config/config.yml",
		"path to weplanx server configuration file",
	)
	var (
		address   string
		openpprof bool
	)
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start HTTP service",
		Run: func(cmd *cobra.Command, args []string) {
			values, err := common.SetValues(config)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			app, err := App(values)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			if openpprof {
				pprof.Register(app)
			}
			app.Run(address)
		},
	}
	serverCmd.Flags().StringVarP(&address,
		"address", "a", "0.0.0.0:9000",
		"binding to address and port",
	)
	serverCmd.Flags().BoolVarP(&openpprof,
		"pprof", "", false,
		"use the pprof tool",
	)
	rootCmd.AddCommand(serverCmd)
	var install bootstrap.Install
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Initialize service data, also supports importing predefined Schema",
		Run: func(cmd *cobra.Command, args []string) {
			values, err := common.SetValues(config)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			client, err := bootstrap.UseMongoDB(values)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			install.Db = bootstrap.UseDatabase(client, values)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err = install.Basic(ctx); err != nil {
				color.Red("%s", err.Error())
				return
			}
			if install.Template == "" {
				return
			}
			if err = install.UseTemplate(ctx); err != nil {
				color.Red("%s", err.Error())
				return
			}
		},
	}
	installCmd.Flags().StringVarP(&install.Username,
		"username", "u", "weplanx",
		"set administrator username",
	)
	installCmd.Flags().StringVarP(&install.Password,
		"password", "p", "",
		"set administrator password",
	)
	installCmd.MarkFlagRequired("password")
	installCmd.Flags().StringVarP(&install.Template,
		"template", "t", "",
		"importing predefined Template",
	)
	rootCmd.AddCommand(installCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
}
