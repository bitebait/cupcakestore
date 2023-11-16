package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	FirstName   string `gorm:"type:varchar(100)" validate:"required"`
	LastName    string `gorm:"type:varchar(100)" validate:"required"`
	Address     string `gorm:"type:varchar(200)"`
	City        string `gorm:"type:varchar(100)"`
	State       string `gorm:"type:varchar(100)"`
	PostalCode  string `gorm:"type:varchar(20)"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	UserID      uint
	User        User // Belongs to relationship with User
}

func (p *Profile) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p *Profile) BeforeSave(tx *gorm.DB) error {
	return p.Validate()
}

func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	return p.Validate()
}

func (p *Profile) Validate() error {
	v := validator.New()
	return v.Struct(p)
}
