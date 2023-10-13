package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(100);" validate:"required"`
	Email      string    `gorm:"type:varchar(100);unique_index" validate:"required,email"`
	Password   string    `gorm:"type:varchar(100);" validate:"required,min=8"`
	IsActive   bool      `gorm:"default:true"`
	IsStaff    bool      `gorm:"default:false"`
	FirstLogin time.Time `gorm:"type:timestamp"`
	LastLogin  time.Time `gorm:"type:timestamp"`
}

func (u *User) Validate() error {
	v := validator.New()
	return v.Struct(u)
}
