package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"time"
)

var Values = &cobra.Command{
	Use:   "values",
	Short: "Display the dynamic values of server distribution Kv",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		values := ctx.Value("values").(*common.Values)

		var x *api.API
		if x, err = bootstrap.NewAPI(values); err != nil {
			return
		}
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
}
