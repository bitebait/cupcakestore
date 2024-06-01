package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StoreConfigController interface {
	Update(ctx *fiber.Ctx) error
	RenderStoreConfig(ctx *fiber.Ctx, configType string) error
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
		return renderErrorMessage(err, "enviar formulário.")
	}

	if err = ctx.BodyParser(storeConfig); err != nil {
		return renderErrorMessage(err, "enviar formulário.")
	}

	if err = c.storeConfigService.Update(storeConfig); err != nil {
		return renderErrorMessage(err, "atualizar configurações da loja.")
	}

	return ctx.Redirect("/dashboard")
}

func (c *storeConfigController) RenderStoreConfig(ctx *fiber.Ctx, configType string) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return renderErrorMessage(err, "carregar configurações da loja.")
	}

	viewPath := "config/" + configType
	return views.Render(ctx, viewPath, storeConfig, views.StoreLayout)
}
