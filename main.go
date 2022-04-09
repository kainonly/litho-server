package main

import (
	"api/bootstrap"
	"github.com/gin-contrib/pprof"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "wpx",
		Version: "v0.0.0",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return nil
		},
	}
	var (
		address   string
		path      string
		openpprof bool
	)
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start HTTP service",
		Run: func(cmd *cobra.Command, args []string) {
			values, err := bootstrap.SetValues(path)
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
	serverCmd.PersistentFlags().StringVarP(&address,
		"address", "a", "0.0.0.0:9000",
		"binding to address and port",
	)
	serverCmd.PersistentFlags().StringVarP(&path,
		"config", "c", "config/config.yml",
		"path to weplanx server configuration file",
	)
	serverCmd.PersistentFlags().BoolVarP(&openpprof,
		"pprof", "", false,
		"use the pprof tool",
	)
	rootCmd.AddCommand(serverCmd)
	var (
		schema string
	)
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Initialize service data, also supports importing predefined Schema",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	installCmd.PersistentFlags().StringVarP(&schema,
		"schema", "s", "",
		"importing predefined Schema",
	)
	rootCmd.AddCommand(installCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
