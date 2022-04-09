package main

import (
	"api/common"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/spf13/cobra"
	"os"
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
				panic(err)
			}
			app, err := App(values)
			if err != nil {
				panic(err)
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
	var (
		username string
		password string
		schema   string
	)
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Initialize service data, also supports importing predefined Schema",
		Args: func(cmd *cobra.Command, args []string) error {
			fmt.Println(args)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	installCmd.Flags().StringVarP(&username,
		"username", "u", "weplanx",
		"set administrator username",
	)
	installCmd.Flags().StringVarP(&password,
		"password", "p", "",
		"set administrator password",
	)
	installCmd.MarkFlagRequired("password")
	installCmd.Flags().StringVarP(&schema,
		"schema", "s", "",
		"importing predefined Schema",
	)
	rootCmd.AddCommand(installCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
