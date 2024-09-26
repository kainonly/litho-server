package cmd

import (
	"github.com/spf13/cobra"
)

func User() *cobra.Command {
	var email string
	var password string
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Create an email account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			//ctx := cmd.Context()
			//values := ctx.Value("values").(*common.Values)
			//
			//var x *api.API
			//if x, err = bootstrap.NewAPI(values); err != nil {
			//	return
			//}
			//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			//defer cancel()
			//if _, err = x.Db.Collection("users").InsertOne(
			//	ctx,
			//	model.NewUser(email, password),
			//); err != nil {
			//	return
			//}
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
	return userCmd
}
