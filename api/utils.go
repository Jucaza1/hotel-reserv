package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/types"
)

func NewFiberAppCentralErr() *fiber.App {
	var config = fiber.Config{
		// Override the error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			errSt, ok := err.(types.ErrorSt)
			if ok {
				fmt.Println(errSt)
				return c.Status(errSt.Status).JSON(types.MsgError{Error: errSt.Msg})
			}
			fmt.Print(err)
			return c.Status(400).JSON(types.MsgError{Error: err.Error()})
		},
	}
	return fiber.New(config)
}
