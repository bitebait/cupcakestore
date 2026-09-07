package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bitebait/cupcakestore/models"
)

const maxPixResponseBytes = 1 << 20

func generatePixPayment(data *models.PixPaymentData) (*models.PixInfo, error) {
	client := &http.Client{
		Timeout:       10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}
	return requestPixPayment(client, "https://pix.ae", data)
}

func requestPixPayment(client *http.Client, endpoint string, data *models.PixPaymentData) (*models.PixInfo, error) {
	if data == nil {
		return nil, errors.New("dados do Pix ausentes")
	}
	form := url.Values{
		"tipo": {data.Tipo}, "chave": {data.Chave}, "location": {data.Location},
		"valor": {data.Valor}, "info": {data.Info}, "nome": {data.Nome}, "txid": {data.Txid},
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serviço Pix retornou HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxPixResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxPixResponseBytes {
		return nil, errors.New("resposta Pix excedeu o limite")
	}
	var result models.PixResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	// Controllers prefix this path with the fixed provider origin.
	path, err := url.Parse(result.Urlpixae)
	if err != nil || path.IsAbs() || path.Host != "" || !strings.HasPrefix(result.Urlpixae, "/") || strings.HasPrefix(result.Urlpixae, "//") || strings.ContainsAny(result.Urlpixae, "\\\r\n") || result.Qrstring == "" {
		return nil, errors.New("resposta Pix inválida")
	}
	return &models.PixInfo{PixQR: result.Qrbase64, PixString: result.Qrstring, PixTransactionID: result.Idfatura, PixURL: result.Urlpixae}, nil
}
