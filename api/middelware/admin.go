package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/types"
)

func AdminMiddleware(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(types.User)
	if !ok {
		return types.ErrUnauthorized(fmt.Errorf("user not found"))
	}
	if !user.IsAdmin {
		return types.ErrUnauthorized(fmt.Errorf("user is not admin"))
	}
	return c.Next()
}
