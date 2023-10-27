package models

import (
	"fmt"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/disintegration/imaging"
	"image"
	"mime/multipart"
	"strings"
)

type ProductImage struct {
	Path string
}

func (i *ProductImage) CropImage(imageFile *multipart.FileHeader) error {
	imageName, err := i.generateRandomImageFileName(imageFile.Filename)
	if err != nil {
		return err
	}

	open, err := imageFile.Open()
	if err != nil {
		return err
	}
	decode, err := imaging.Decode(open)
	if err != nil {
		return err
	}

	croppedImage := i.cropImage(decode)
	err = i.saveCroppedImage(imageName, croppedImage)
	if err != nil {
		return err
	}

	return nil
}

func (i *ProductImage) cropImage(srcImage image.Image) image.Image {
	return imaging.Thumbnail(srcImage, 400, 400, imaging.Lanczos)
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
	rand := utils.NewRandomizer()
	randString, err := rand.GenerateRandomString(22)
	if err != nil {
		return "", err
	}
	return randString + "." + strings.Split(filename, ".")[1], nil
}
