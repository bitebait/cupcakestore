package views

import (
	"bytes"
	"net/url"
	"strings"
	"testing"

	"github.com/Masterminds/sprig/v3"
	"github.com/bitebait/cupcakestore/models"
	templatehtml "github.com/gofiber/template/html/v2"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

func renderTemplate(t *testing.T, name string, data any, layout ...string) string {
	t.Helper()
	engine := htmlEngine(t)
	var output bytes.Buffer
	if err := engine.Render(&output, name, data, layout...); err != nil {
		t.Fatalf("render %s: %v", name, err)
	}
	return output.String()
}

func htmlEngine(t *testing.T) *templatehtml.Engine {
	t.Helper()
	engine := templatehtml.New(".", ".html")
	engine.AddFuncMap(sprig.FuncMap())
	engine.AddFunc("money", Money)
	if err := engine.Load(); err != nil {
		t.Fatal(err)
	}
	return engine
}

func TestPaginationPreservesSearch(t *testing.T) {
	query := `chocolate & "morango" <script>alert(1)</script>`
	output := renderTemplate(t, "snippets/pagination", map[string]any{
		"Query":      query,
		"Pagination": map[string]any{"Page": 2, "Limit": 8, "Total": 30, "TotalPages": 4},
	})
	if strings.Contains(output, "<script>") {
		t.Fatal("search rendered as executable markup")
	}
	document, err := html.Parse(strings.NewReader(output))
	if err != nil {
		t.Fatal(err)
	}
	links := 0
	var inspect func(*html.Node)
	inspect = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					link, err := url.Parse(attr.Val)
					if err != nil || link.Query().Get("q") != query {
						t.Errorf("search was not preserved in %q", attr.Val)
					}
					links++
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			inspect(child)
		}
	}
	inspect(document)
	if links != 4 {
		t.Fatalf("expected four pagination links, got %d", links)
	}
}

func TestStoreRendersForGuestWithoutProducts(t *testing.T) {
	output := renderTemplate(t, "store/store", map[string]any{
		"Object": map[string]any{
			"Products": []any{},
			"Filter": map[string]any{
				"Product":    map[string]any{"Name": ""},
				"Pagination": map[string]any{"Page": 1, "Limit": 8, "Total": 0, "TotalPages": 0},
			},
		},
		"Messages":  map[string]any{},
		"CSRFToken": "test-token",
	}, StoreLayout)
	if !strings.Contains(output, "Nenhum produto disponível") || strings.Contains(output, "Página 1 de 0") {
		t.Fatal("guest empty state is missing or has invalid pagination")
	}
	if strings.Contains(output, "jquery") {
		t.Fatal("storefront must work without jQuery")
	}
}

func TestShoppingCartMutationFormsContainCSRFToken(t *testing.T) {
	output := renderTemplate(t, "shoppingcart/shoppingcart", map[string]any{
		"Object": models.ShoppingCart{
			Model: gorm.Model{ID: 3}, Total: 5,
			Items: []models.ShoppingCartItem{{ProductID: 2, Product: models.Product{Model: gorm.Model{ID: 2}, Name: `<img src=x onerror=alert(1)>`, Price: 5}, ItemPrice: 5, Quantity: 1}},
		},
		"Profile":   &models.Profile{FirstName: "Ana", LastName: "Silva", Address: "Rua Teste", City: "São Paulo", State: "SP", PostalCode: "01000-000", PhoneNumber: "11999999999", UserID: 1},
		"CSRFToken": "cart-token",
	})
	document, err := html.Parse(strings.NewReader(output))
	if err != nil {
		t.Fatal(err)
	}
	forms := map[string]string{}
	var inspect func(*html.Node)
	inspect = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "form" {
			attributes := map[string]string{}
			for _, attr := range node.Attr {
				attributes[attr.Key] = attr.Val
			}
			if strings.EqualFold(attributes["method"], "post") {
				token := ""
				var findToken func(*html.Node)
				findToken = func(child *html.Node) {
					if child.Data == "input" {
						values := map[string]string{}
						for _, attr := range child.Attr {
							values[attr.Key] = attr.Val
						}
						if values["name"] == "_csrf" {
							token = values["value"]
						}
					}
					for nested := child.FirstChild; nested != nil; nested = nested.NextSibling {
						findToken(nested)
					}
				}
				findToken(node)
				forms[attributes["action"]] = token
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			inspect(child)
		}
	}
	inspect(document)
	for _, action := range []string{"/cart/remove/2", "/cart/items/2/quantity", "/orders/checkout/3"} {
		if forms[action] != "cart-token" {
			t.Errorf("missing protected POST form for %s", action)
		}
	}

	if strings.Contains(output, "<img src=x") {
		t.Fatal("product name was not escaped")
	}
}

func TestOrderManagementOffersOnlyValidNextSteps(t *testing.T) {
	for _, test := range []struct {
		name      string
		order     models.Order
		allowed   string
		forbidden string
	}{
		{"pickup ready", models.Order{Status: models.DeliveredStatusAwaiting}, "Entregue", "Enviado"},
		{"delivery ready", models.Order{Status: models.DeliveredStatusAwaiting, IsDelivery: true}, "Enviado", "Entregue"},
		{"pix awaiting", models.Order{Status: models.AwaitingPaymentStatus, PaymentMethod: models.PixPaymentMethod}, "Pagamento Aprovado", "Entregue"},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.order.ID = 7
			output := renderTemplate(t, "orders/order", map[string]any{
				"Object":    map[string]any{"Order": test.order},
				"Profile":   &models.Profile{User: models.User{IsStaff: true}},
				"CSRFToken": "status-token",
			})
			if !strings.Contains(output, `value="`+test.allowed+`"`) || strings.Contains(output, `value="`+test.forbidden+`"`) {
				t.Fatal("status selector does not match the valid next steps")
			}
			if strings.Contains(output, `value="Cancelado"`) || !strings.Contains(output, `href="/orders/cancel/7"`) {
				t.Fatal("cancellation must go through the confirmation page")
			}
		})
	}
}
