package observability

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

type Controller struct {
	ObservabilityService *Service
}

func (x *Controller) GetCgoCalls(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetCgoCalls(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoUptime(ctx context.Context, c *app.RequestContext) {
	value, err := x.ObservabilityService.GetMongoUptime(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"value": value,
	})
}

func (x *Controller) GetMongoAvailableConnections(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoAvailableConnections(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoOpenConnections(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoOpenConnections(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoCommandsPerSecond(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoCommandsPerSecond(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoQueryOperations(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoQueryOperations(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoDocumentOperations(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoDocumentOperations(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoFlushes(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoFlushes(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (x *Controller) GetMongoNetworkIO(ctx context.Context, c *app.RequestContext) {
	data, err := x.ObservabilityService.GetMongoNetworkIO(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}
