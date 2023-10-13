package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) CheckPassword(inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
}
