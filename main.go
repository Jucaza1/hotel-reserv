package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jucaza1/hotel-reserv/api"
	middleware "github.com/jucaza1/hotel-reserv/api/middelware"
	"github.com/jucaza1/hotel-reserv/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	db.DBURI = os.Getenv("MONGO_DB_URI")
	db.DBNAME = os.Getenv("MONGO_DB_NAME")
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	if listenAddr == "" {
		log.Fatal("error: HTTP_LISTEN_ADDRESS not found in .env")
	}
	if db.DBNAME == "" {
		log.Fatal("error: MONGO_DB_NAME not found in .env")
	}
	if db.DBURI == "" {
		log.Fatal("error: MONGO_DB_URI not found in .env")
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
		roomHandler    = api.NewRoomHandler(rStore, hStore)
		bookingHandler = api.NewBookingHandler(bStore, rStore)
		authHandler    = api.NewAuthHandler(uStore)
		app            = api.NewFiberAppCentralErr()
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(uStore))
		admin          = apiv1.Group("/admin", middleware.AdminMiddleware)
	)

	//auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	//version api
	apiv1.Get("/users/me", userHandler.HandleGetMyUser)
	apiv1.Post("/users", userHandler.HandlePostUser)

	//hotel handler
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)

	//room handler
	apiv1.Get("/hotels/:hid/rooms", roomHandler.HandleGetRooms)
	apiv1.Get("/hotels/:hid/rooms/id", roomHandler.HandleGetRoomByID)

	//booking handler
	apiv1.Get("/rooms/:id/bookings", bookingHandler.HandleGetBookingsByRoom)
	apiv1.Post("/rooms/:id/bookings", bookingHandler.HandlePostBooking)
	apiv1.Get("/hotels/:hid/bookings", bookingHandler.HandleGetBookingsByHotel)
	apiv1.Get("/bookings", bookingHandler.HandleGetBookings)
	apiv1.Patch("/bookings/:id", bookingHandler.HandleCancelBooking)

	//admin only user handlers
	admin.Patch("/users/:id", userHandler.HandlePatchUser)
	admin.Delete("/users/:id", userHandler.HandleDeleteUser)
	admin.Post("/users", userHandler.HandlePostUser)
	admin.Post("/users/admin", userHandler.HandlePostAdminUser)
	admin.Get("/users/me", userHandler.HandleGetMyUser)
	admin.Get("/users", userHandler.HandleGetUsers)
	admin.Get("/users/:id", userHandler.HandleGetUser)

	//admin only room handlers
	admin.Delete("/rooms/id", roomHandler.HandleDeleteRoom)
	admin.Post("/hotels/:hid/rooms/", roomHandler.HandlePostRoom)

	//admin only hotel handler
	admin.Delete("/hotels/:id", hotelHandler.HandleDeleteHotel, roomHandler.HandleDeleteRoomsByHotel)
	admin.Post("/hotels", hotelHandler.HandlePostHotel)
	admin.Patch("/hotels", hotelHandler.HandlePatchHotel)
	admin.Delete("/bookings/:id", bookingHandler.HandleDeleteBooking)

	log.Println("app listening on port ", listenAddr)
	log.Fatal(app.Listen(listenAddr))
}

//docker run -d --name YOUR_CONTAINER_NAME_HERE -p YOUR_LOCALHOST_PORT_HERE:27017 -e MONGO_INITDB_ROOT_USERNAME=YOUR_USERNAME_HERE -e MONGO_INITDB_ROOT_PASSWORD=YOUR_PASSWORD_HERE mongo
//docker run --name mongodb -p 27017:27017 -d mongo:latest
//go get go.mongodb.org/mongo-driver/mongo
