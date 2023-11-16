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

func (c storeConfigController) Update(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (c storeConfigController) RenderStoreConfig(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		//TODO Change Users to Dashboard
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "config/config", storeConfig, "", baseLayout)
}
