package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/models"
)

// Banco Central, Manual de Padrões para Iniciação do Pix v2.10.0, section 2.6.3.
const officialPixVector = "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D3D"

func TestPixPayloadMatchesOfficialReferenceAndCRC(t *testing.T) {
	payload := buildPixPayload("123e4567-e12b-12d1-a456-426655440000", "Fulano de Tal", "BRASILIA", "", "***")
	if payload != officialPixVector {
		t.Fatalf("payload differs from Banco Central reference: %s", payload)
	}
	if err := ValidatePixPayload(officialPixVector); err != nil {
		t.Fatal(err)
	}
	// Independent standard CRC-16/CCITT-FALSE check vector.
	if got := pixCRC16("123456789"); got != 0x29b1 {
		t.Fatalf("CRC = %04X, want 29B1", got)
	}
}

func TestLocalPixPreservesAmountReferenceAndProducesPNG(t *testing.T) {
	data := localPixFixture()
	payment, err := generateLocalPixPayment(&data)
	if err != nil {
		t.Fatal(err)
	}
	if payment.PixURL != "" || payment.PixTransactionID != data.Txid {
		t.Fatalf("unexpected external URL or transaction reference: %#v", payment)
	}
	fields, err := parsePixTLV(payment.PixString)
	if err != nil {
		t.Fatal(err)
	}
	if fields["54"] != "12.50" || fields["59"] != "Doces Sao Joao" || fields["60"] != "Sao Paulo" || fields["62"] != "0509CUPCAKE42" {
		t.Fatalf("incorrect payment terms: %#v", fields)
	}
	decoded, err := base64.StdEncoding.DecodeString(payment.PixQR)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("invalid PNG: %v", err)
	}
	if img.Bounds().Dx() != img.Bounds().Dy() || img.Bounds().Dx() < 200 {
		t.Fatalf("unexpected QR dimensions: %v", img.Bounds())
	}
	// Keep the four-module white quiet zone, required for reliable camera reading.
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if x >= 24 && y >= 24 && x < img.Bounds().Dx()-24 && y < img.Bounds().Dy()-24 {
				continue
			}
			r, g, b, _ := img.At(x, y).RGBA()
			if r != 65535 || g != 65535 || b != 65535 {
				t.Fatalf("QR quiet zone is not white at %d,%d", x, y)
			}
		}
	}
	regenerated, err := PixQRFromPayload(payment.PixString)
	if err != nil || regenerated != payment.PixQR {
		t.Fatalf("regenerating stored payment changed QR: %v", err)
	}
}

func TestLocalPixRejectsInvalidTermsBeforeGeneratingInstructions(t *testing.T) {
	if _, err := generateLocalPixPayment(nil); err == nil {
		t.Fatal("accepted missing Pix data")
	}
	for _, amount := range []string{"", "0.00", "-1.00", "1", "1.0", "1.001", "1,00", "NaN", "1e2", "01.00", "10000000000.00"} {
		data := localPixFixture()
		data.Valor = amount
		if _, err := generateLocalPixPayment(&data); err == nil {
			t.Fatalf("accepted invalid amount %q", amount)
		}
	}
	for _, txid := range []string{"", "***", "order-1", strings.Repeat("a", 26), "á123"} {
		data := localPixFixture()
		data.Txid = txid
		if _, err := generateLocalPixPayment(&data); err == nil {
			t.Fatalf("accepted invalid txid %q", txid)
		}
	}
	for _, edit := range []func(*models.PixPaymentData){
		func(data *models.PixPaymentData) { data.Chave = "not-an-email" },
		func(data *models.PixPaymentData) { data.Nome = "" },
		func(data *models.PixPaymentData) { data.City = "" },
	} {
		data := localPixFixture()
		edit(&data)
		if _, err := generateLocalPixPayment(&data); err == nil {
			t.Fatal("accepted incomplete receiver configuration")
		}
	}
}

func TestPixQRRejectsCorruptOrUnsupportedHistoricalPayloads(t *testing.T) {
	withCRC := func(body string) string { return body + "6304" + fmt.Sprintf("%04X", pixCRC16(body+"6304")) }
	for _, payload := range []string{
		"", "copy-code", "https://pix.ae/invoice/1", officialPixVector[:len(officialPixVector)-1] + "0",
		strings.Replace(officialPixVector, "Fulano", "Fulana", 1), officialPixVector + "extra",
		withCRC(strings.TrimSuffix(officialPixVector, "63041D3D") + "5303986"),
		withCRC(strings.Replace(strings.TrimSuffix(officialPixVector, "63041D3D"), "2658", "2699", 1)),
		withCRC(strings.Replace(strings.TrimSuffix(officialPixVector, "63041D3D"), "5303986", "5303840", 1)),
		withCRC(strings.Replace(strings.TrimSuffix(officialPixVector, "63041D3D"), "0136", "2536", 1)),
		withCRC("000201" + pixTLV("26", pixTLV("00", "br.gov.bcb.pix")+pixTLV("01", "loja@example.com")+pixTLV("03", "00000000")) + "5204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***"),
	} {
		if _, err := PixQRFromPayload(payload); err == nil {
			t.Fatalf("accepted invalid stored Pix %q", payload)
		}
	}
}

func localPixFixture() models.PixPaymentData {
	return models.PixPaymentData{Tipo: "email", Chave: "loja@example.com", Nome: "Doces São João", City: "São Paulo", Valor: "12.50", Txid: "CUPCAKE42"}
}

func TestValidateOrderPixRejectsMismatchedOrUnspecifiedPaymentTerms(t *testing.T) {
	data := localPixFixture()
	payment, err := generateLocalPixPayment(&data)
	if err != nil {
		t.Fatal(err)
	}
	valid := models.Order{Total: 12.50, PixString: payment.PixString, PixTransactionID: data.Txid}
	if err := ValidateOrderPix(&valid); err != nil {
		t.Fatal(err)
	}
	for _, edit := range []func(*models.Order){
		func(order *models.Order) { order.Total = 12.51 },
		func(order *models.Order) { order.PixTransactionID = "otherOrder" },
		func(order *models.Order) { order.PixTransactionID = "" },
		func(order *models.Order) { order.PixString = officialPixVector; order.PixTransactionID = "***" },
	} {
		order := valid
		edit(&order)
		if err := ValidateOrderPix(&order); err == nil {
			t.Fatalf("accepted mismatched Pix %#v", order)
		}
	}
	if err := ValidateOrderPix(nil); err == nil {
		t.Fatal("accepted missing order")
	}
}
