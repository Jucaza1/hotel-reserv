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
		uStore         = db.NewMongoUserStore(client, db.DBNAME)
		hStore         = db.NewMongoHotelStore(client, db.DBNAME)
		rStore         = db.NewMongoRoomStore(client, db.DBNAME, hStore)
		bStore         = db.NewMongoBookingStore(client, db.DBNAME)
		userHandler    = api.NewUserHandler(uStore)
		hotelHandler   = api.NewHotelHandler(hStore)
		roomHandler    = api.NewRoomHandler(rStore)
		bookingHandler = api.NewBookingHandler(bStore, rStore)
		authHandler    = api.NewAuthHandler(uStore)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(uStore))
	)

	//auth
	auth.Post("/auth/", authHandler.HandleAuthenticate)

	//version api
	//user handlers
	apiv1.Patch("/users/:id", userHandler.HandlePatchUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/users/", userHandler.HandlePostUser)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)

	//hotel handler
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)

	//room handler
	apiv1.Get("/hotels/:id/rooms", roomHandler.HandleGetRooms)

	//booking handler
	apiv1.Get("/hotels/:idh/rooms/:id/bookings", bookingHandler.HandleGetBookingsByRoom)
	apiv1.Post("/hotels/:idh/rooms/:id/bookings", bookingHandler.HandlePostBooking)
	apiv1.Get("/hotels/:idh/bookings", bookingHandler.HandleGetBookingsByHotel)
	log.Println(*listenAddr)
	app.Listen(*listenAddr)
}

//docker run -d --name YOUR_CONTAINER_NAME_HERE -p YOUR_LOCALHOST_PORT_HERE:27017 -e MONGO_INITDB_ROOT_USERNAME=YOUR_USERNAME_HERE -e MONGO_INITDB_ROOT_PASSWORD=YOUR_PASSWORD_HERE mongo
//docker run --name mongodb -p 27017:27017 -d mongo:latest
//go get go.mongodb.org/mongo-driver/mongo
