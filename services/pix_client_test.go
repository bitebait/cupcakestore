package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bitebait/cupcakestore/models"
)

func TestPixRequestEncodesFormAndChecksResponse(t *testing.T) {
	data := &models.PixPaymentData{Chave: "a+b@example.com", Info: "Cupcake & café", Valor: "12.50"}
	for _, tc := range []struct {
		name    string
		status  int
		body    string
		wantErr bool
	}{
		{"valid", 200, `{"urlpixae":"/invoice/123","qrstring":"pix-code"}`, false},
		{"provider error", 500, `{"urlpixae":"/invoice/123","qrstring":"pix-code"}`, true},
		{"invalid json", 200, `<html>error</html>`, true},
		{"external URL", 200, `{"urlpixae":"https://example.com","qrstring":"pix-code"}`, true},
		{"host suffix", 200, `{"urlpixae":".example.com/path","qrstring":"pix-code"}`, true},
		{"protocol relative", 200, `{"urlpixae":"//example.com/path","qrstring":"pix-code"}`, true},
		{"missing payment data", 200, `{}`, true},
		{"oversized response", 200, strings.Repeat(" ", maxPixResponseBytes+1), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				if r.Method != "POST" || r.Form.Get("chave") != data.Chave || r.Form.Get("info") != data.Info {
					t.Errorf("invalid form: %v", r.Form)
				}
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			result, err := requestPixPayment(server.Client(), server.URL, data)
			if (err != nil) != tc.wantErr {
				t.Fatalf("result=%+v err=%v", result, err)
			}
		})
	}
}

func TestPixRequestHonorsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_ = json.NewEncoder(w).Encode(models.PixResponse{Urlpixae: "/invoice/123", Qrstring: "code"})
	}))
	defer server.Close()
	client := server.Client()
	client.Timeout = time.Millisecond
	if _, err := requestPixPayment(client, server.URL, &models.PixPaymentData{}); err == nil {
		t.Fatal("expected timeout")
	}
}
