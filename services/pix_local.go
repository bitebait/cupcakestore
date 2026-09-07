package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitebait/cupcakestore/models"
	"github.com/yeqown/go-qrcode/v2"
)

var (
	pixAmountPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.[0-9]{2}$`)
	pixTxIDPattern   = regexp.MustCompile(`^[a-zA-Z0-9]{1,25}$`)
)

// generateLocalPixPayment creates payment instructions, never a payment receipt.
// BR Code/Pix specification: https://www.bcb.gov.br/content/estabilidadefinanceira/pix/Regulamento_Pix/II_ManualdePadroesparaIniciacaodoPix.pdf
func generateLocalPixPayment(data *models.PixPaymentData) (*models.PixInfo, error) {
	if data == nil {
		return nil, errors.New("dados do Pix ausentes")
	}
	key, err := models.NormalizePixKey(data.Tipo, data.Chave)
	if err != nil {
		return nil, err
	}
	name, city, err := models.NormalizePixMerchant(data.Nome, data.City)
	if err != nil {
		return nil, err
	}
	if !validPixAmount(data.Valor) {
		return nil, errors.New("o valor Pix deve ser positivo, com duas casas decimais e até 13 caracteres")
	}
	if !pixTxIDPattern.MatchString(data.Txid) {
		return nil, errors.New("o identificador Pix deve conter de 1 a 25 letras ou números")
	}
	payload := buildPixPayload(key, name, city, data.Valor, data.Txid)
	qr, err := PixQRFromPayload(payload)
	if err != nil {
		return nil, err
	}
	return &models.PixInfo{PixQR: qr, PixString: payload, PixTransactionID: data.Txid}, nil
}

func validPixAmount(value string) bool {
	return len(value) <= 13 && pixAmountPattern.MatchString(value) && value != "0.00"
}

// Inputs are validated before encoding; the optional amount also supports the
// official BR Code reference vector. Checkout always supplies a fixed amount.
func buildPixPayload(key, name, city, amount, txid string) string {
	account := pixTLV("00", "br.gov.bcb.pix") + pixTLV("01", key)
	payload := pixTLV("00", "01") + pixTLV("26", account) + pixTLV("52", "0000") + pixTLV("53", "986")
	if amount != "" {
		payload += pixTLV("54", amount)
	}
	payload += pixTLV("58", "BR") + pixTLV("59", name) + pixTLV("60", city) + pixTLV("62", pixTLV("05", txid)) + "6304"
	return payload + fmt.Sprintf("%04X", pixCRC16(payload))
}

func pixTLV(id, value string) string {
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

// CRC-16/CCITT-FALSE: polynomial 0x1021, initial value 0xffff, no reflection/XOR.
func pixCRC16(payload string) uint16 {
	crc := uint16(0xffff)
	for i := 0; i < len(payload); i++ {
		crc ^= uint16(payload[i]) << 8
		for bit := 0; bit < 8; bit++ {
			if crc&0x8000 != 0 {
				crc = crc<<1 ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

// PixQRFromPayload renders a stored static BR Code as a PNG in plain base64.
// Validation detects malformed historical data; it does not prove key ownership
// or that the bank received funds. The caller must authorize access to the order.
func PixQRFromPayload(payload string) (string, error) {
	if err := ValidatePixPayload(payload); err != nil {
		return "", err
	}
	qr, err := qrcode.NewWith(payload, qrcode.WithEncodingMode(qrcode.EncModeByte), qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	if err != nil {
		return "", errors.New("não foi possível gerar o QR Code Pix")
	}
	writer := &pixPNGWriter{}
	if err := qr.Save(writer); err != nil {
		return "", errors.New("não foi possível gerar a imagem Pix")
	}
	return base64.StdEncoding.EncodeToString(writer.Bytes()), nil
}

// ValidateOrderPix binds stored payment instructions to the order's recorded
// amount and reference. A valid checksum alone cannot establish those terms.
func ValidateOrderPix(order *models.Order) error {
	if order == nil {
		return errors.New("pedido Pix ausente")
	}
	if err := ValidatePixPayload(order.PixString); err != nil {
		return err
	}
	fields, _ := parsePixTLV(order.PixString)
	additional, _ := parsePixTLV(fields["62"])
	amount := fmt.Sprintf("%.2f", order.Total)
	if !validPixAmount(amount) || fields["54"] != amount || !pixTxIDPattern.MatchString(order.PixTransactionID) || additional["05"] != order.PixTransactionID {
		return errors.New("o Pix armazenado não corresponde ao valor ou identificador do pedido; entre em contato com a loja")
	}
	return nil
}

// ValidatePixPayload checks the static BR Code structure and checksum; use
// ValidateOrderPix as well before presenting an order's payment instructions.
func ValidatePixPayload(payload string) error {
	invalid := errors.New("o código Pix armazenado está inválido; entre em contato com a loja")
	if len(payload) < 8 || len(payload) > 512 || !strings.HasPrefix(payload, "000201") || payload[len(payload)-8:len(payload)-4] != "6304" {
		return invalid
	}
	for _, character := range payload {
		if character < 32 || character > 126 {
			return invalid
		}
	}
	if fmt.Sprintf("%04X", pixCRC16(payload[:len(payload)-4])) != payload[len(payload)-4:] {
		return invalid
	}
	fields, err := parsePixTLV(payload)
	if err != nil || fields["00"] != "01" || fields["53"] != "986" || fields["58"] != "BR" || len(fields["52"]) != 4 {
		return invalid
	}
	for id := range fields {
		switch id {
		case "00", "01", "26", "52", "53", "54", "58", "59", "60", "61", "62", "63":
		default:
			return invalid
		}
	}
	for _, digit := range fields["52"] {
		if digit < '0' || digit > '9' {
			return invalid
		}
	}
	if method := fields["01"]; method != "" && method != "11" && method != "12" {
		return invalid
	}
	account, err := parsePixTLV(fields["26"])
	if err != nil || !strings.EqualFold(account["00"], "br.gov.bcb.pix") || account["01"] == "" || len(account["01"]) > 77 || account["25"] != "" {
		return invalid
	}
	// This store accepts ordinary static payments. Other account extensions
	// can describe cash withdrawal or a dynamic URL and need different flows.
	for id := range account {
		if id != "00" && id != "01" && id != "02" {
			return invalid
		}
	}
	keyValid := false
	for _, kind := range []string{"email", "celular", "aleatoria", "cpf", "cnpj"} {
		if _, err := models.NormalizePixKey(kind, account["01"]); err == nil {
			keyValid = true
			break
		}
	}
	if !keyValid {
		return invalid
	}
	additional, err := parsePixTLV(fields["62"])
	if err != nil || (additional["05"] != "***" && !pixTxIDPattern.MatchString(additional["05"])) {
		return invalid
	}
	if fields["54"] != "" && !validPixAmount(fields["54"]) {
		return invalid
	}
	name, city, err := models.NormalizePixMerchant(fields["59"], fields["60"])
	if err != nil || name != fields["59"] || city != fields["60"] {
		return invalid
	}
	return nil
}

func parsePixTLV(value string) (map[string]string, error) {
	fields := make(map[string]string)
	for value != "" {
		if len(value) < 4 || value[0] < '0' || value[0] > '9' || value[1] < '0' || value[1] > '9' || value[2] < '0' || value[2] > '9' || value[3] < '0' || value[3] > '9' {
			return nil, errors.New("campo Pix inválido")
		}
		id := value[:2]
		length, _ := strconv.Atoi(value[2:4])
		if _, duplicate := fields[id]; duplicate || length < 1 || len(value) < 4+length {
			return nil, errors.New("tamanho ou identificador Pix inválido")
		}
		fields[id] = value[4 : 4+length]
		value = value[4+length:]
	}
	return fields, nil
}

type pixPNGWriter struct{ bytes.Buffer }

func (w *pixPNGWriter) Close() error { return nil }

func (w *pixPNGWriter) Write(matrix qrcode.Matrix) error {
	const scale, quietZone = 6, 4
	width := (matrix.Width() + 2*quietZone) * scale
	img := image.NewPaletted(image.Rect(0, 0, width, width), color.Palette{color.White, color.Black})
	matrix.Iterate(qrcode.IterDirection_ROW, func(x, y int, value qrcode.QRValue) {
		if !value.IsSet() {
			return
		}
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetColorIndex((x+quietZone)*scale+dx, (y+quietZone)*scale+dy, 1)
			}
		}
	})
	return png.Encode(&w.Buffer, img)
}
