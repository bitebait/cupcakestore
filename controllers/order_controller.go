package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bitebait/cupcakestore/helpers"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type OrderController interface {
	RenderOrder(ctx *fiber.Ctx) error
	RenderAllOrders(ctx *fiber.Ctx) error
	Checkout(ctx *fiber.Ctx) error
	Payment(ctx *fiber.Ctx) error
	RenderCancel(ctx *fiber.Ctx) error
	Cancel(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type orderController struct {
	orderService       services.OrderService
	storeConfigService services.StoreConfigService
}

func NewOrderController(orderService services.OrderService, storeConfigService services.StoreConfigService) OrderController {
	return &orderController{
		orderService:       orderService,
		storeConfigService: storeConfigService,
	}
}

func (c *orderController) Checkout(ctx *fiber.Ctx) error {
	profileID := getUserID(ctx)
	currentUser := ctx.Locals("Profile").(*models.Profile)

	if !currentUser.IsProfileComplete() {
		messages.SetErrorMessage(ctx, "por favor, complete as informações do perfil para prosseguir")
		return ctx.Redirect("/profile")
	}

	cartID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		messages.SetErrorMessage(ctx, "erro ao processar o ID do carrinho")
		return ctx.Redirect("/orders")
	}

	order, err := c.orderService.FindOrCreate(profileID, cartID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	if !c.isAuthorizedUser(currentUser, &order, profileID) || !order.CanProceedToCheckout() {
		return ctx.Redirect("/orders")
	}

	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	data := fiber.Map{
		"Order":       order,
		"StoreConfig": storeConfig,
	}

	return ctx.Render("orders/checkout", fiber.Map{"Object": data}, views.StoreLayout)
}

func (c *orderController) Payment(ctx *fiber.Ctx) error {
	profileID := getUserID(ctx)
	cartID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		messages.SetErrorMessage(ctx, "erro ao processar o ID do carrinho")
		return ctx.Redirect("/orders")
	}

	order, err := c.orderService.FindByCartId(cartID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	currentUser := ctx.Locals("Profile").(*models.Profile)
	if !c.isAuthorizedUser(currentUser, &order, profileID) || !order.CanProceedToPayment() {
		return ctx.Redirect("/orders")
	}

	switch ctx.Method() {
	case fiber.MethodPost:
		pixURL, err := c.processPaymentPost(ctx, &order)
		if err != nil {
			messages.SetErrorMessage(ctx, err.Error())
			return c.processPaymentGet(ctx, &order)
		}

		if pixURL != "" {
			return ctx.Redirect(pixURL, fiber.StatusSeeOther)
		}

	case fiber.MethodGet:
		return c.processPaymentGet(ctx, &order)

	default:
		return ctx.Status(fiber.StatusMethodNotAllowed).SendString("Método não permitido")
	}

	return ctx.Redirect("/orders")
}

func (c *orderController) isAuthorizedUser(user *models.Profile, order *models.Order, profileID uint) bool {
	return user.User.IsStaff || order.IsCurrentUserOrder(profileID)
}

func (c *orderController) processPaymentPost(ctx *fiber.Ctx, order *models.Order) (string, error) {
	if err := ctx.BodyParser(order); err != nil {
		return "", fiber.NewError(fiber.StatusBadRequest, "Erro ao processar os dados do pedido. Verifique os campos e tente novamente.")
	}

	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return "", fmt.Errorf("falha ao carregar configuração da loja: %w", err)
	}
	order.IsDelivery = storeConfig.DeliveryIsActive

	if err := c.orderService.Update(order); err != nil {
		return "", fmt.Errorf("falha ao atualizar pedido: %w", err)
	}

	if err := c.orderService.Payment(order); err != nil {
		return "", fmt.Errorf("falha ao processar pagamento: %w", err)
	}

	if order.PaymentMethod == models.PixPaymentMethod {
		if order.PixURL == "" {
			return "", errors.New("QR Code PIX não foi gerado corretamente")
		}
		pixURL := "https://pix.ae" + order.PixURL
		return pixURL, nil
	}

	return "", errors.New("método de pagamento não suportado ou resposta não implementada")
}

func (c *orderController) processPaymentGet(ctx *fiber.Ctx, order *models.Order) error {
	if order.CanRedirectToPixPayment() {
		return ctx.Redirect("https://pix.ae" + order.PixURL)
	}

	messages.SetErrorMessage(ctx, "não foi possível redirecionar para pagamento via Pix")
	return ctx.Redirect("/orders")
}

func (c *orderController) RenderOrder(ctx *fiber.Ctx) error {
	orderID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	order, err := c.orderService.FindById(orderID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	currentUser := ctx.Locals("Profile").(*models.Profile)
	if !c.isAuthorizedUser(currentUser, &order, currentUser.ID) {
		return ctx.Redirect("/orders")
	}

	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	data := fiber.Map{
		"Order":       order,
		"StoreConfig": storeConfig,
	}

	return ctx.Render("orders/order", fiber.Map{"Object": data}, views.StoreLayout)
}

func (c *orderController) RenderAllOrders(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("Profile").(*models.Profile)
	filter := models.NewOrderFilter(currentUser.ID, ctx.QueryInt("page"), ctx.QueryInt("limit"))

	var orders []models.Order
	var templateName string
	var layout string

	if currentUser.User.IsStaff {
		orders = c.orderService.FindAll(filter)
		templateName = "orders/admin"
		layout = views.BaseLayout
	} else {
		orders = c.orderService.FindAllByUser(filter)
		templateName = "orders/orders"
		layout = views.StoreLayout
	}

	data := fiber.Map{"Orders": orders, "Filter": filter}

	return ctx.Render(templateName, fiber.Map{"Object": data}, layout)
}

func (c *orderController) RenderCancel(ctx *fiber.Ctx) error {
	orderID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		messages.SetErrorMessage(ctx, "ID do pedido é inválido")
		return ctx.Redirect("/orders")
	}

	order, err := c.orderService.FindById(orderID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	currentUser := ctx.Locals("Profile").(*models.Profile)
	if !c.isAuthorizedUser(currentUser, &order, currentUser.ID) {
		return ctx.Redirect("/orders")
	}

	return ctx.Render("orders/cancel", fiber.Map{"Object": order}, views.StoreLayout)
}

func (c *orderController) Cancel(ctx *fiber.Ctx) error {
	orderID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	order, err := c.orderService.FindById(orderID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/orders")
	}

	user := ctx.Locals("Profile").(*models.Profile)
	if c.isAuthorizedUser(user, &order, user.ID) {
		if err := c.orderService.Cancel(order.ID); err != nil {
			messages.SetErrorMessage(ctx, err.Error())
			return ctx.Redirect("/orders")
		}
	}

	messages.SetSuccessMessage(ctx, "pedido cancelado com sucesso")
	return ctx.Redirect("/orders")
}

func (c *orderController) Update(ctx *fiber.Ctx) error {
	orderID, parseErr := helpers.ParseStringToID(ctx.Params("id"))
	if parseErr != nil {
		messages.SetErrorMessage(ctx, "ID do pedido é inválido")
		return ctx.Redirect("/orders")
	}

	order, findErr := c.orderService.FindById(orderID)
	if findErr != nil {
		messages.SetErrorMessage(ctx, findErr.Error())
		return ctx.Redirect("/orders")
	}

	if order.Status == models.CancelledStatus {
		messages.SetErrorMessage(ctx, "o pedido foi cancelado")
		return ctx.Redirect("/orders")
	}

	newStatus := models.ShoppingCartStatus(ctx.FormValue("status"))
	order.Status = newStatus

	currentUser := ctx.Locals("Profile").(*models.Profile)
	if currentUser.User.IsStaff {
		updateErr := c.orderService.Update(&order)
		if updateErr != nil {
			messages.SetErrorMessage(ctx, updateErr.Error())
			return ctx.Redirect("/orders")
		}
	}

	messages.SetSuccessMessage(ctx, "status do pedido atualizado com sucesso")
	return ctx.Redirect("/orders/order/" + strconv.Itoa(int(order.ID)))
}
