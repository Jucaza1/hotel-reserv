package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/jucaza1/hotel-reserv/api"
	middleware "github.com/jucaza1/hotel-reserv/api/middelware"
	"github.com/jucaza1/hotel-reserv/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	// Override d error handler
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":4000", "The listen addres of the API server")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	//var initialization
	var (
		uStore       = db.NewMongoUserStore(client, db.DBNAME)
		hStore       = db.NewMongoHotelStore(client, db.DBNAME)
		rStore       = db.NewMongoRoomStore(client, db.DBNAME, hStore)
		userHandler  = api.NewUserHandler(uStore)
		hotelHandler = api.NewHotelHandler(hStore, rStore)
		authHandler  = api.NewAuthHandler(uStore)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	//auth
	auth.Post("/auth/", authHandler.HandleAuthenticate)

	//version api
	//user handlers
	apiv1.Patch("/user/:id", userHandler.HandlePatchUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user/", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)

	//hotel handler
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandleGetRooms)

	log.Println(*listenAddr)
	app.Listen(*listenAddr)
}

//docker run -d --name YOUR_CONTAINER_NAME_HERE -p YOUR_LOCALHOST_PORT_HERE:27017 -e MONGO_INITDB_ROOT_USERNAME=YOUR_USERNAME_HERE -e MONGO_INITDB_ROOT_PASSWORD=YOUR_PASSWORD_HERE mongo
//docker run --name mongodb -p 27017:27017 -d mongo:latest
//go get go.mongodb.org/mongo-driver/mongo
