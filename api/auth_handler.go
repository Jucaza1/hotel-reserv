package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email   string `json:"email"`
	Pasword string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return types.ErrInvalidParams(err)
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return types.ErrUnauthorized(fmt.Errorf("invalid credentials"))
	}
	if !types.AuthUser(user.EncyptedPassword, params.Pasword) {
		return types.ErrUnauthorized(fmt.Errorf("invalid credentials"))
	}
	token := createTokenFromUser(user)
	if len(token) == 0 {
		return types.ErrInternal(fmt.Errorf("error creating token"))
	}
	c.Response().Header.Add("X-Authorization", token)
	return c.SendStatus(http.StatusNoContent)
}

func createTokenFromUser(user *types.User) string {
	exp := time.Now().Add(time.Hour * 4)
	claims := jwt.MapClaims{
		"id":      user.ID,
		"expires": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if len(secret) == 0 {
		fmt.Println("failed read secret from env")
	}
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")
	}
	return tokenStr
}
