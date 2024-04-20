package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Firstname) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.Lastname) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.+[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func ValidateUserUpdate(update UpdateUser) (updateValid map[string]string, errors map[string]string) {
	errors = map[string]string{}
	updateValid = map[string]string{}
	if update.Firstname != "" {
		if len(update.Firstname) < minFirstNameLen {
			errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
		} else {
			updateValid["firstName"] = update.Firstname
		}
	}
	if update.Lastname != "" {
		if len(update.Lastname) < minLastNameLen {
			errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
		} else {
			updateValid["lastName"] = update.Lastname
		}
	}
	if update.Email != "" {
		if !isEmailValid(update.Email) {
			errors["email"] = "email is invalid"
		} else {
			updateValid["email"] = update.Email
		}
	}
	return updateValid, errors
}

type User struct {
	ID               string `bson:"_id,omitempty" json:"id,omitempty"`
	Firstname        string `bson:"firstName" json:"firstName"`
	Lastname         string `bson:"lastName" json:"lastName"`
	Email            string `bson:"email" json:"email"`
	EncyptedPassword string `bson:"password" json:"-"`
	IsAdmin          bool   `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Firstname:        params.Firstname,
		Lastname:         params.Lastname,
		Email:            params.Email,
		EncyptedPassword: string(encpw),
	}, nil
}

func NewAdminFromParams(params CreateUserParams) (*User, error) {
	user, err := NewUserFromParams(params)
	if err != nil {
		return nil, err
	}
	user.IsAdmin = true
	return user, nil
}

func AuthUser(encpw, pw string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw))
}

type UpdateUser struct {
	Firstname string `json:"firstName,omitempty"`
	Lastname  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
}
