package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StoreConfigController interface {
	Update(ctx *fiber.Ctx) error
	RenderStoreConfigAddress(ctx *fiber.Ctx) error
	RenderStoreConfigPayment(ctx *fiber.Ctx) error
	RenderStoreConfigPix(ctx *fiber.Ctx) error
	RenderStoreConfigDelivery(ctx *fiber.Ctx) error
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

	err = ctx.BodyParser(storeConfig)
	if err != nil {
		return views.Render(ctx, "config/config", storeConfig, err.Error(), baseLayout)
	}

	storeConfig.DeliveryIsActive = ctx.FormValue("deliveryIsActive") == "on"
	storeConfig.PaymentCashIsActive = ctx.FormValue("paymentCashIsActive") == "on"
	storeConfig.PaymentPixIsActive = ctx.FormValue("paymentPixIsActive") == "on"
	err = c.storeConfigService.Update(storeConfig)
	if err != nil {
		return views.Render(ctx, "config/config", storeConfig, err.Error(), baseLayout)
	}

	return ctx.Redirect("/users")
}

func (c *storeConfigController) RenderStoreConfigAddress(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "config/address", storeConfig, "", baseLayout)
}

func (c *storeConfigController) RenderStoreConfigPayment(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "config/payment", storeConfig, "", baseLayout)
}

func (c *storeConfigController) RenderStoreConfigPix(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "config/pix", storeConfig, "", baseLayout)
}

func (c *storeConfigController) RenderStoreConfigDelivery(ctx *fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "config/delivery", storeConfig, "", baseLayout)
}
