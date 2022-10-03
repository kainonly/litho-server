package main

import (
	"context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/weplanx/server/bootstrap"
	"os"
	"time"
)

func main() {
	var path string
	rootCmd := &cobra.Command{
		Version: "v0.0.0",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&path,
		"config", "c", "config/config.yml",
		"配置文件路径",
	)
	var (
		address string
	)
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "启动 HTTP 服务",
		Run: func(cmd *cobra.Command, args []string) {
			values, err := bootstrap.LoadStaticValues(path)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			api, err := bootstrap.NewAPI(values)
			if err != nil {
				color.Red("%s", err.Error())
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			h, err := api.Initialize(ctx)
			if err != nil {
				color.Red("%s", err.Error())
			}

			h.Spin()
		},
	}
	serverCmd.Flags().StringVarP(&address,
		"address", "a", "0.0.0.0:3000",
		"监听端口",
	)
	rootCmd.AddCommand(serverCmd)
	var install bootstrap.Install
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "初始化应用数据",
		Run: func(cmd *cobra.Command, args []string) {
			values, err := bootstrap.LoadStaticValues(path)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			client, err := bootstrap.UseMongoDB(values)
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			install.Db = bootstrap.UseDatabase(values, client)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err = install.Basic(ctx); err != nil {
				color.Red("%s", err.Error())
				return
			}
		},
	}
	installCmd.Flags().StringVarP(&install.Username,
		"username", "u", "weplanx",
		"管理员用户名",
	)
	installCmd.Flags().StringVarP(&install.Password,
		"password", "p", "",
		"管理员用户密码",
	)
	rootCmd.AddCommand(installCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
}
