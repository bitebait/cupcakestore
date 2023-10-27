package models

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"path/filepath"
)

type Thumbnail struct {
	Image string
}

func (t *Thumbnail) GetPath() string {
	thumb := "/images/thumbs/" + filepath.Base(t.Image)
	if thumb == "" {
		return t.Image
	}

	return thumb
}

func (t *Thumbnail) CreateThumbnail() error {
	imageName := t.getImageName()

	srcImage, err := t.openImage(imageName)
	if err != nil {
		return err
	}

	thumbnail := t.createThumbnailImage(srcImage)
	err = t.saveThumbnailImage(imageName, thumbnail)
	if err != nil {
		return err
	}
	return nil
}

func (t *Thumbnail) getImageName() string {
	return filepath.Base(t.Image)
}

func (t *Thumbnail) openImage(imageName string) (image.Image, error) {
	srcImagePath := fmt.Sprintf("./web/images/%s", imageName)
	return imaging.Open(srcImagePath)
}

func (t *Thumbnail) createThumbnailImage(srcImage image.Image) image.Image {
	return imaging.Thumbnail(srcImage, 100, 100, imaging.CatmullRom)
}

func (t *Thumbnail) saveThumbnailImage(imageName string, thumbnail image.Image) error {
	thumbnailPath := fmt.Sprintf("./web/images/thumbs/%s", imageName)
	return imaging.Save(thumbnail, thumbnailPath)
}
