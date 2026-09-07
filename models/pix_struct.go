package models

type PixPaymentData struct {
	Tipo  string `json:"tipo"`
	Chave string `json:"chave"`
	Valor string `json:"valor"`
	Nome  string `json:"nome"`
	City  string `json:"city"`
	Txid  string `json:"txid"`
}

type PixInfo struct {
	PixQR            string `gorm:"default:''"`
	PixString        string `gorm:"default:''"`
	PixTransactionID string `gorm:"default:''"`
	PixURL           string `gorm:"default:''"`
}
