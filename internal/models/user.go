package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"  validate:"omitempty"`
	FirstName string    `json:"first_name" db:"first_name"  validate:"required,lte=30"`
	LastName  string    `json:"last_name" db:"last_name"  validate:"required,lte=30"`
	Email     string    `json:"email,omitempty" db:"email"  validate:"omitempty,lte=60,email"`
	Password  string    `json:"password,omitempty" db:"password"  validate:"omitempty,required,gte=6"`
	Avatar    []byte    `json:"avatar,omitempty" db:"avatar"`
	Country   *string   `json:"country,omitempty" db:"country"  validate:"omitempty,lte=24"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" `
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" `
	LoginDate time.Time `json:"login_date" db:"login_date" `
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return nil
}
