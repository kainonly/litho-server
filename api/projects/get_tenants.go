package projects

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetTenantsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetTenants(ctx context.Context, c *app.RequestContext) {
	var dto GetTenantsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ProjectsX.GetTenants(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) GetTenants(ctx context.Context, id primitive.ObjectID) (result M, err error) {
	var project model.Project
	if err = x.Db.Collection("projects").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&project); err != nil {
		return
	}

	result = M{}
	if project.Nats != nil {
		var b []byte
		if b, err = x.Cipher.Decode(project.Nats.Seed); err != nil {
			return
		}
		result["nkey"] = string(b)
	}

	return
}
