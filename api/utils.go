package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/jucaza1/hotel-reserv/db"
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

//test env initialization

func injectENV(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	db.DBURI = os.Getenv("MONGO_DB_URI")
	db.TestDBNAME = os.Getenv("MONGO_DB_TESTNAME")
	if db.TestDBNAME == "" {
		t.Fatal("error: MONGO_DB_TESTNAME not found in .env")
	}
	if db.DBURI == "" {
		t.Fatal("error: MONGO_DB_URI not found in .env")
	}
}
