package models

import (
	"fmt"
	"github.com/bitebait/cupcakestore/utils"
	"mime/multipart"
	"strings"
)

type ProductImage struct {
	FileName  string
	FilePath  string
	ImagePath string
}

func (p *ProductImage) CreateProductImage(imageFile *multipart.FileHeader) error {
	var err error

	p.FileName, err = p.generateRandomImageFileName(imageFile.Filename)
	if err != nil {
		return err
	}

	p.FilePath = fmt.Sprintf("./web/images/%s", p.FileName)

	p.ImagePath = fmt.Sprintf("/images/%s", p.FileName)
	return nil
}

func (p *ProductImage) generateRandomImageFileName(filename string) (string, error) {
	rand := utils.NewRandomizer()
	randString, err := rand.GenerateRandomString(22)
	if err != nil {
		return "", err
	}
	return randString + "." + strings.Split(filename, ".")[1], nil
}
