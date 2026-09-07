package models

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func imageFile(t *testing.T, data []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", "no-extension")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(4 << 20); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = request.MultipartForm.RemoveAll() })
	return request.MultipartForm.File["image"][0]
}

func TestProductImageValidatesBeforeFullDecode(t *testing.T) {
	var valid bytes.Buffer
	if err := png.Encode(&valid, image.NewRGBA(image.Rect(0, 0, 10, 10))); err != nil {
		t.Fatal(err)
	}
	// Alter only the PNG dimensions and checksum; full decoding would allocate
	// an enormous image, but DecodeConfig lets us reject it first.
	huge := append([]byte(nil), valid.Bytes()...)
	binary.BigEndian.PutUint32(huge[16:20], 100000)
	binary.BigEndian.PutUint32(huge[20:24], 100000)
	binary.BigEndian.PutUint32(huge[29:33], crc32.ChecksumIEEE(huge[12:29]))
	for _, tc := range []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"filename without extension", valid.Bytes(), false},
		{"invalid image", []byte("not an image"), true},
		{"too many pixels", huge, true},
		{"too many bytes", make([]byte, maxImageBytes+1), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			img, err := (&ProductImage{}).cropImage(imageFile(t, tc.data))
			if (err != nil) != tc.wantErr {
				t.Fatalf("image=%v err=%v", img, err)
			}
		})
	}
}
