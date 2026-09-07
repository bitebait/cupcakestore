package views

import (
	"bytes"
	"net/url"
	"strings"
	"testing"

	"github.com/Masterminds/sprig/v3"
	templatehtml "github.com/gofiber/template/html/v2"
	"golang.org/x/net/html"
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
	if strings.Count(output, `src="/plugins/jquery/jquery.min.js"`) != 1 {
		t.Fatal("store must load jQuery exactly once")
	}
}

func TestShoppingCartMutationFormsContainCSRFToken(t *testing.T) {
	output := renderTemplate(t, "shoppingcart/shoppingcart", map[string]any{
		"Object": map[string]any{
			"ID":    3,
			"Total": 5.0,
			"Items": []any{map[string]any{
				"Product":   map[string]any{"ID": 2, "Name": `<img src=x onerror=alert(1)>`},
				"ItemPrice": 5.0, "Quantity": 1,
			}},
		},
		"Profile":   map[string]any{"IsProfileComplete": true},
		"CSRFToken": "cart-token",
	})
	for _, action := range []string{"/cart/remove/2", "/orders/checkout/3"} {
		if !strings.Contains(output, `action="`+action+`" method="post"`) {
			t.Errorf("missing POST form for %s", action)
		}
	}
	if strings.Count(output, `name="_csrf" value="cart-token"`) != 2 {
		t.Fatal("cart mutation forms must render the root CSRF token, including inside item range")
	}
	if strings.Contains(output, "<img src=x") {
		t.Fatal("product name was not escaped")
	}
}
