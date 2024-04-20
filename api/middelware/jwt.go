package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
)

func JWTAuthentication(us db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Authorization")
		if len(token) == 0 {
			return types.ErrUnauthorized(fmt.Errorf("token not present in headers"))
		}
		claims, err := validateToken(token)
		if err != nil {
			return types.ErrUnauthorized(err)
		}
		//check expiration
		tm, ok := claims["expires"].(float64)
		if !ok {
			return types.ErrUnauthorized(fmt.Errorf("expire date invalid"))
		}
		remaining := getTokenRemainingValidity(tm)
		if remaining <= 0 {
			return types.ErrUnauthorized(fmt.Errorf("token expired"))
		}
		//check and save user
		userID, _ := claims["id"].(string)
		user, err := us.GetUserByID(c.Context(), userID)
		if err != nil || userID != user.ID {
			return types.ErrUnauthorized(fmt.Errorf("token user not in database"))
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}
func getTokenRemainingValidity(timestamp float64) int {
	tm := time.Unix(int64(timestamp), 0)
	remainer := time.Until(tm)
	return int(remainer.Seconds())
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %s", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to parse claims")
	}
	return claims, nil
}
