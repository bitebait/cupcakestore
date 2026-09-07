package models

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
