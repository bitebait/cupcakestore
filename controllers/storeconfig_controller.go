package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StoreConfigController interface {
	Update(ctx *fiber.Ctx) error
	RenderStoreConfig(ctx *fiber.Ctx) error
}

type storeConfigController struct {
	storeConfigService services.StoreConfigService
}

func NewStoreConfigController(s services.StoreConfigService) StoreConfigController {
	return &storeConfigController{
		storeConfigService: s,
	}
}

func (c *storeConfigController) Update(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return err
	}

	err = ctx.BodyParser(&storeConfig)
	if err != nil {
		if err != nil {
			return views.Render(ctx, "config/config", storeConfig, err.Error(), baseLayout)
		}
	}

	err = c.storeConfigService.Update(storeConfig)
	if err != nil {
		return views.Render(ctx, "config/config", storeConfig, err.Error(), baseLayout)
	}

	return ctx.Redirect("/config")
}

func (c storeConfigController) RenderStoreConfig(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		//TODO Change Users to Dashboard
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "config/config", storeConfig, "", baseLayout)
}
