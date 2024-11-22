package controllers

import (
	"errors"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/services"
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
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/")
	}

	if err = ctx.BodyParser(&storeConfig); err != nil {
		messages.SetErrorMessage(ctx,
			errors.New("erro ao processar a configuração, verifique os dados e tente novamente").Error())
		return ctx.Redirect("/")
	}

	if err = c.storeConfigService.Update(&storeConfig); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/")
	}

	messages.SetSuccessMessage(ctx, "configuração atualizada com sucesso")
	return ctx.Redirect("/dashboard")
}

func (c *storeConfigController) RenderStoreConfig(ctx *fiber.Ctx, configType string) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/")
	}

	viewPath := "config/" + configType
	return ctx.Render(viewPath, fiber.Map{"Object": storeConfig}, "layouts/base")
}
