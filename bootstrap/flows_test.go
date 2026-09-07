package bootstrap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
)

func TestStoreFlowRegistrationLoginCheckoutAndCancellation(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	cfg := *config.Get()
	cfg.DBType, cfg.DBPath = "sqlite", filepath.Join(t.TempDir(), "flow.db")
	cfg.AdminEmail, cfg.AdminPassword = "", ""
	db, err := database.Open(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	database.DB = db
	session.SetupSession()
	app := createFiberApp()
	registerMiddlewares(app)
	registerRoutes(app, database.DB)
	cookies := map[string]*http.Cookie{}
	csrfPattern := regexp.MustCompile(`name="_csrf"[^>]*value="([^"]+)"`)
	token := ""
	request := func(method, path string, form url.Values, wantStatus int) string {
		t.Helper()
		if form == nil {
			form = url.Values{}
		}
		if method == "POST" {
			form.Set("_csrf", token)
		}
		req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		if method == "POST" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		response, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		for _, cookie := range response.Cookies() {
			cookies[cookie.Name] = cookie
		}
		if response.StatusCode != wantStatus {
			t.Fatalf("%s %s = %d want %d: %s", method, path, response.StatusCode, wantStatus, body)
		}
		if match := csrfPattern.FindStringSubmatch(string(body)); len(match) > 1 {
			token = match[1]
		}
		return string(body)
	}
	request("GET", "/store", nil, 200)
	request("GET", "/auth/register", nil, 200)
	request("POST", "/auth/register", url.Values{"email": {"customer@example.com"}, "password": {"customer-password"}, "firstname": {"Ana"}, "lastname": {"Silva"}, "isStaff": {"true"}, "ID": {"99"}}, 302)
	var user models.User
	if err := db.Where("email = ?", "customer@example.com").First(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.IsStaff || user.ID == 99 {
		t.Fatal("registration accepted privileged fields")
	}
	request("GET", "/auth/login", nil, 200)
	oldSession := cookies["session_id"].Value
	request("POST", "/auth/login", url.Values{"email": {" CUSTOMER@EXAMPLE.COM "}, "password": {"customer-password"}}, 302)
	if cookies["session_id"].Value == oldSession {
		t.Fatal("login did not rotate session")
	}
	request("GET", "/users", nil, 403)
	userID := strconv.Itoa(int(user.ID))
	request("GET", "/profile/"+userID, nil, 200)
	request("POST", "/profile/update/"+userID, url.Values{
		"firstname": {"Ana"}, "lastname": {"Silva"}, "address": {"Rua Teste"}, "city": {"São Paulo"}, "state": {"SP"}, "postalcode": {"01000-000"}, "phonenumber": {"11999999999"}, "UserID": {"999"}, "ID": {"999"},
	}, 302)
	var profile models.Profile
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		t.Fatal(err)
	}
	if !profile.IsProfileComplete() || profile.ID == 999 {
		t.Fatal("profile update failed or accepted internal ID")
	}
	product := models.Product{Name: "Cupcake", Price: 12.5, IsActive: true}
	if err := db.Create(&product).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.Stock{ProfileID: profile.ID, ProductID: product.ID, Quantity: 5, Type: models.StockEntrada}).Error; err != nil {
		t.Fatal(err)
	}
	productID := strconv.Itoa(int(product.ID))
	request("GET", "/cart", nil, 200)
	request("POST", "/cart", url.Values{"id": {productID}, "quantity": {"2"}}, 302)
	var cart models.ShoppingCart
	if err := db.Where("profile_id = ?", profile.ID).First(&cart).Error; err != nil {
		t.Fatal(err)
	}
	cartID := strconv.Itoa(int(cart.ID))
	// GET must never create an order or reserve inventory.
	request("GET", "/orders/checkout/"+cartID, nil, 302)
	var count int64
	if err := db.Model(&models.Order{}).Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("GET created order: %d %v", count, err)
	}
	request("GET", "/cart", nil, 200)
	request("POST", "/orders/checkout/"+cartID, nil, 200)
	request("POST", "/orders/payment/"+cartID, url.Values{"paymentMethod": {"Dinheiro"}, "isDelivery": {"0"}, "Total": {"0.01"}, "ProfileID": {"999"}, "Status": {"Entregue"}}, 303)
	var order models.Order
	if err := db.Where("shopping_cart_id = ?", cart.ID).First(&order).Error; err != nil {
		t.Fatal(err)
	}
	if order.Total != 25 || order.IsDelivery || order.ProfileID != profile.ID || order.Status != models.ProcessingStatus {
		t.Fatalf("incorrect payment: %+v", order)
	}
	orderID := strconv.Itoa(int(order.ID))
	request("GET", "/orders/order/"+orderID, nil, 200)
	request("POST", "/orders/order/"+orderID, url.Values{"status": {"Entregue"}}, 403)
	request("GET", "/orders/cancel/"+orderID, nil, 200)
	request("POST", "/orders/cancel/"+orderID, nil, 302)
	request("POST", "/orders/cancel/"+orderID, nil, 302)
	if err := db.First(&product, product.ID).Error; err != nil {
		t.Fatal(err)
	}
	if product.CurrentStock != 5 {
		t.Fatalf("stock after cancellation = %d", product.CurrentStock)
	}
	request("GET", "/auth/logout", nil, 405)
	request("POST", "/auth/logout", nil, 302)
	request("GET", "/cart", nil, 302)
}
