package types

import (
	"errors"
	"net/mail"
)

type SignInInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

func (this SignInInput) IsValidSignInInput() error {
	email, err := mail.ParseAddress(this.Email)
	if err != nil {
		return errors.New("Invalid_Email")
	}
	this.Email = email.Address
	if len(this.PasswordHash) == 0 {
		return errors.New("Invalid_PasswordHash")
	}
	return nil
}

type SignUpInput struct {
	// TODO: Use interface instead
	SignInInput
	FullName string `json:"fullName"`
}

func (this SignUpInput) IsValidSignUpInput() error {
	if err := this.IsValidSignInInput(); err != nil {
		return err
	}
	if len(this.FullName) == 0 {
		return errors.New("Invalid_FullName")
	}
	return nil
}
