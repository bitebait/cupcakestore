package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	FirstName   string `gorm:"type:varchar(100)"`
	LastName    string `gorm:"type:varchar(100)"`
	Address     string `gorm:"type:varchar(200)"`
	City        string `gorm:"type:varchar(100)"`
	State       string `gorm:"type:varchar(100)"`
	PostalCode  string `gorm:"type:varchar(20)"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	UserID      uint   `gorm:"unique;not null"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (p *Profile) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p *Profile) Validate() error {
	v := validator.New()
	return v.Struct(p)
}

func (p *Profile) IsProfileComplete() bool {
	return p.FirstName != "" &&
		p.LastName != "" &&
		p.Address != "" &&
		p.City != "" &&
		p.State != "" &&
		p.PostalCode != "" &&
		p.PhoneNumber != ""
}
