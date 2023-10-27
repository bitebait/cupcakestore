package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null,type:varchar(100)"`
	Description string  `gorm:"not null,type:varchar(300)"`
	Price       float64 `gorm:"not null"`
	Ingredients string  `gorm:"not null,type:varchar(300)"`
	Image       string
	Thumbnail   string
}

func (p *Product) Validate() error {
	v := validator.New()
	return v.Struct(p)
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *Product) AfterCreate(tx *gorm.DB) error {
	t := &Thumbnail{
		Image: p.Image,
	}
	if err := t.CreateThumbnail(); err != nil {
		return err
	}

	tx.Model(p).Update("thumbnail", t.GetPath())
	return nil
}

func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	t := &Thumbnail{
		Image: p.Image,
	}

	if err := t.CreateThumbnail(); err != nil {
		return err
	}

	p.Thumbnail = t.GetPath()
	return nil
}
