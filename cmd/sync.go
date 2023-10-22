package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
)

var Sync = &cobra.Command{
	Use:   "sync",
	Short: "Sync weplanx models",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		values := ctx.Value("values").(*common.Values)

		var x *api.API
		if x, err = bootstrap.NewAPI(values); err != nil {
			return
		}
		if err = model.SetProjects(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetUsers(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetClusters(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetSchedules(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetWorkflows(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetQueues(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetImessages(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetAccTasks(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetCategories(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetPictures(ctx, x.Db); err != nil {
			return
		}
		if err = model.SetVideos(ctx, x.Db); err != nil {
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
		return
	},
}
