package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const maxImageBytes = 4 << 20
const maxImagePixels = 16_000_000

type ProductImage struct{ Path string }

func (i *ProductImage) Save(file *multipart.FileHeader) error {
	thumbnail, err := i.cropImage(file)
	if err != nil {
		return err
	}
	// User filenames never determine the path or output format.
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return err
	}
	name := hex.EncodeToString(random[:]) + ".jpg"
	if err := imaging.Save(thumbnail, filepath.Join("web", "images", name)); err != nil {
		return err
	}
	i.Path = "/images/" + name
	return nil
}

func (i *ProductImage) cropImage(file *multipart.FileHeader) (image.Image, error) {
	if file == nil || file.Size <= 0 || file.Size > maxImageBytes {
		return nil, errors.New("a imagem deve ter até 4 MB")
	}
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	config, format, err := image.DecodeConfig(io.LimitReader(reader, maxImageBytes))
	if err != nil {
		return nil, errors.New("imagem inválida")
	}
	if format != "jpeg" && format != "png" && format != "gif" {
		return nil, errors.New("use uma imagem JPEG, PNG ou GIF")
	}
	if config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > maxImagePixels {
		return nil, errors.New("a imagem deve ter até 16 milhões de pixels")
	}
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	decoded, err := imaging.Decode(io.LimitReader(reader, maxImageBytes))
	if err != nil {
		return nil, err
	}
	return imaging.Thumbnail(decoded, 400, 400, imaging.Lanczos), nil
}
