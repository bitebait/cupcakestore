package models

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/bitebait/cupcakestore/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserFilter struct {
	User       *User
	Pagination *Pagination
}

func NewUserFilter(query string, page, limit int) *UserFilter {
	user := &User{
		Username: query,
		Email:    query,
	}
	pagination := NewPagination(page, limit)

	return &UserFilter{
		User:       user,
		Pagination: pagination,
	}
}

type User struct {
	gorm.Model
	Username   string `gorm:"type:varchar(100);" validate:"required"`
	Email      string `gorm:"type:varchar(100);unique_index" validate:"required,email"`
	Password   string `gorm:"type:varchar(100);" validate:"required,min=8"`
	IsActive   bool
	IsStaff    bool
	FirstLogin time.Time `gorm:"type:timestamp"`
	LastLogin  time.Time `gorm:"type:timestamp"`
}

func (u *User) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *User) HashPassword() error {
	hash, err := utils.PasswordHasher(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash

	return nil
}

func (u *User) CheckPassword(inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	return
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	profile := &Profile{
		FirstName: "",
		LastName:  "",
		UserID:    u.ID,
	}
	return tx.Create(profile).Error
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	if err = tx.Where("user_id = ?", u.ID).Delete(&Profile{}).Error; err != nil {
		return err
	}
	return nil
}
