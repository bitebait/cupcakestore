package controllers

import (
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v3"
)

type StoreConfigController interface {
	Update(ctx fiber.Ctx) error
	RenderStoreConfig(ctx fiber.Ctx, configType string) error
}

type storeConfigController struct {
	storeConfigService services.StoreConfigService
}

func NewStoreConfigController(s services.StoreConfigService) StoreConfigController {
	return &storeConfigController{
		storeConfigService: s,
	}
}

func (c *storeConfigController) Update(ctx fiber.Ctx) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/")
	}

	// The four settings forms submit partial updates. Start with current values,
	// then copy only editable settings back, preserving the database identity.
	input := storeConfig
	if err := ctx.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "configuração inválida")
	}
	storeConfig = models.StoreConfig{
		Model:         storeConfig.Model,
		DeliveryPrice: input.DeliveryPrice, DeliveryIsActive: input.DeliveryIsActive,
		PhysicalStoreEmail: input.PhysicalStoreEmail, PhysicalStoreAddress: input.PhysicalStoreAddress,
		PhysicalStoreCity: input.PhysicalStoreCity, PhysicalStoreState: input.PhysicalStoreState,
		PhysicalStorePostalCode: input.PhysicalStorePostalCode, PhysicalStorePhoneNumber: input.PhysicalStorePhoneNumber,
		PaymentCashIsActive: input.PaymentCashIsActive, PaymentPixIsActive: input.PaymentPixIsActive,
		PixKey: input.PixKey, PixKeyType: input.PixKeyType,
	}

	if err = c.storeConfigService.Update(&storeConfig); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/")
	}

	messages.SetSuccessMessage(ctx, "configuração atualizada com sucesso")
	return ctx.Redirect().Status(fiber.StatusFound).To("/dashboard")
}

func (c *storeConfigController) RenderStoreConfig(ctx fiber.Ctx, configType string) error {
	storeConfig, err := c.storeConfigService.GetStoreConfig()

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/")
	}

	viewPath := "config/" + configType
	return ctx.Render(viewPath, fiber.Map{"Object": storeConfig}, "layouts/base")
}
