package models

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"image"
	"strings"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null,type:varchar(100)"`
	Description string  `gorm:"not null,type:varchar(300)"`
	Price       float64 `gorm:"not null"`
	Ingredients string  `gorm:"not null,type:varchar(300)"`
	Image       string
}

func (p *Product) createThumbnail() error {
	imageName := p.getImageName()
	srcImage, err := p.openImage(imageName)
	if err != nil {
		return err
	}

	thumbnail := p.createThumbnailImage(srcImage)
	err = p.saveThumbnailImage(imageName, thumbnail)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) getImageName() string {
	return strings.Split(p.Image, "/")[2]
}

func (p *Product) openImage(imageName string) (image.Image, error) {
	srcImagePath := fmt.Sprintf("./web/images/%s", imageName)
	return imaging.Open(srcImagePath)
}

func (p *Product) createThumbnailImage(srcImage image.Image) image.Image {
	return imaging.Thumbnail(srcImage, 100, 100, imaging.CatmullRom)
}

func (p *Product) saveThumbnailImage(imageName string, thumbnail image.Image) error {
	thumbnailPath := fmt.Sprintf("./web/images/thumbs/%s", imageName)
	return imaging.Save(thumbnail, thumbnailPath)
}
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}

func (p *Product) AfterCreate(tx *gorm.DB) error {
	if err := p.createThumbnail(); err != nil {
		return err
	}

	return nil
}

func (p *Product) Validate() error {
	v := validator.New()
	return v.Struct(p)
}

func (p *Product) Thumbnail() string {
	thumb := "/images/thumbs/" + strings.Split(p.Image, "/")[2]
	if thumb == "" {
		return p.Image
	}
	return thumb
}
