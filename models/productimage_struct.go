package models

import (
	"fmt"
	"github.com/bitebait/cupcakestore/helpers"
	"github.com/disintegration/imaging"
	"image"
	"mime/multipart"
	"strings"
)

type ProductImage struct {
	Path string
}

func (i *ProductImage) Save(imageFile *multipart.FileHeader) error {
	imageName, err := i.generateRandomImageFileName(imageFile.Filename)
	if err != nil {
		return err
	}

	croppedImage, err := i.cropImage(imageFile)
	if err != nil {
		return err
	}

	err = i.saveCroppedImage(imageName, croppedImage)
	if err != nil {
		return err
	}

	return nil
}

func (i *ProductImage) cropImage(imageFile *multipart.FileHeader) (image.Image, error) {
	open, err := imageFile.Open()
	if err != nil {
		return nil, err
	}
	defer open.Close()

	decoded, err := imaging.Decode(open)
	if err != nil {
		return nil, err
	}

	croppedImage := imaging.Thumbnail(decoded, 400, 400, imaging.Lanczos)

	return croppedImage, nil
}

func (i *ProductImage) saveCroppedImage(imageName string, thumbnail image.Image) error {
	imagePath := fmt.Sprintf("./web/images/%s", imageName)

	err := imaging.Save(thumbnail, imagePath)
	if err != nil {
		return err
	}

	i.Path = strings.ReplaceAll(imagePath, "./web/", "/")

	return nil
}

func (i *ProductImage) generateRandomImageFileName(filename string) (string, error) {
	rand := helpers.NewRandomizer()
	randString, err := rand.GenerateString(22)
	if err != nil {
		return "", err
	}
	return randString + "." + strings.Split(filename, ".")[1], nil
}
