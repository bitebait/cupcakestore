package models

import (
	"errors"
	"time"

	"github.com/bitebait/cupcakestore/utils"
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

func (u *User) BeforeCreate(tx *gorm.DB) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	if u == nil || u.Password == "" {
		return errors.New("invalid user or empty password")
	}

	hash, err := utils.PasswordHasher(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func (u *User) UpdatePassword(oldPassword, newPassword string) error {
	if err := u.CheckPassword(oldPassword); err != nil {
		return err
	}

	if newPassword == "" {
		return errors.New("new password cannot be empty")
	}

	hash, err := utils.PasswordHasher(newPassword)
	if err != nil {
		return err
	}

	u.Password = hash

	return nil
}

func (u *User) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *User) CheckPassword(inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
}
