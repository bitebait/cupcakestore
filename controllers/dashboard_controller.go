package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type DashboardController interface {
	RenderDashboard(ctx *fiber.Ctx) error
}

type dashboardController struct {
	dashboardService services.DashboardService
}

func NewDashboardController(s services.DashboardService) DashboardController {
	return &dashboardController{
		dashboardService: s,
	}
}

func (c *dashboardController) RenderDashboard(ctx *fiber.Ctx) error {
	data := c.dashboardService.GetInfo(30)
	return ctx.Render("dashboard/dashboard", fiber.Map{"Object": data}, views.BaseLayout)
}
