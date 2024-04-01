package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/mongo"
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
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if !types.AuthUser(user.EncyptedPassword, params.Pasword) {
		return fmt.Errorf("invalid credentials")
	}
	token := createTokenFromUser(user)
	if len(token) == 0 {
		return fmt.Errorf("internal error")
	}
	c.Response().Header.Add("X-Authorization", token)
	return c.SendStatus(http.StatusNoContent)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	exp := now.Add(time.Hour * 4)
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
