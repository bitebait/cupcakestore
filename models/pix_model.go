package models

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type PixPaymentData struct {
	Tipo     string `json:"tipo"`
	Chave    string `json:"chave"`
	Location string `json:"location"`
	Valor    string `json:"valor"`
	Info     string `json:"info"`
	Nome     string `json:"nome"`
	Txid     string `json:"txid"`
}

type PixInfo struct {
	PixQR            string `gorm:"default:''"`
	PixString        string `gorm:"default:''"`
	PixTransactionID string `gorm:"default:''"`
	PixURL           string `gorm:"default:''"`
}

type PixResponse struct {
	Status   string `json:"status"`
	Qrbase64 string `json:"qrbase64"`
	Qrstring string `json:"qrstring"`
	Idfatura string `json:"idfatura"`
	Urlpixae string `json:"urlpixae"`
}

func GeneratePixPayment(data *PixPaymentData) (*PixInfo, error) {
	formData := make([]string, 0)
	formData = append(formData, "tipo="+data.Tipo)
	formData = append(formData, "chave="+data.Chave)
	formData = append(formData, "location="+data.Location)
	formData = append(formData, "valor="+data.Valor)
	formData = append(formData, "info="+data.Info)
	formData = append(formData, "nome="+data.Nome)
	formData = append(formData, "txid="+data.Txid)

	formBody := strings.NewReader(strings.Join(formData, "&"))

	req, err := http.NewRequest("POST", "https://pix.ae", formBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var pixResponse PixResponse
	err = json.NewDecoder(resp.Body).Decode(&pixResponse)
	if err != nil {
		return nil, err
	}

	pixInfo := &PixInfo{
		PixQR:            pixResponse.Qrbase64,
		PixString:        pixResponse.Qrstring,
		PixTransactionID: pixResponse.Idfatura,
		PixURL:           pixResponse.Urlpixae,
	}

	return pixInfo, nil
}
