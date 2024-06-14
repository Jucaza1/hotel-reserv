package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	us, hs, rs, bs := initialization()
	var hotels [3]string
	hotels[0] = seedHotel("Maria", "Spain", 5, hs)
	hotels[1] = seedHotel("Rose", "France", 4, hs)
	hotels[2] = seedHotel("Sheena", "Portugal", 3, hs)

	var rooms [12]string
	rooms[0] = seedRoom(hotels[0], types.Small, 100, rs)
	rooms[1] = seedRoom(hotels[0], types.Normal, 120, rs)
	rooms[2] = seedRoom(hotels[0], types.Large, 129, rs)
	rooms[3] = seedRoom(hotels[0], types.Extra, 135, rs)

	rooms[4] = seedRoom(hotels[1], types.Small, 100, rs)
	rooms[5] = seedRoom(hotels[1], types.Normal, 120, rs)
	rooms[6] = seedRoom(hotels[1], types.Large, 129, rs)
	rooms[7] = seedRoom(hotels[1], types.Extra, 135, rs)

	rooms[8] = seedRoom(hotels[2], types.Small, 100, rs)
	rooms[9] = seedRoom(hotels[2], types.Normal, 120, rs)
	rooms[10] = seedRoom(hotels[2], types.Large, 129, rs)
	rooms[11] = seedRoom(hotels[2], types.Extra, 135, rs)

	var users [4]string
	users[0] = seedUser(true, "admin", "instrator", "admin@mail.com", "mysecretpassword", us)
	users[1] = seedUser(false, "Levi", "Ackerman", "levi@mail.com", "topsecret", us)
	users[2] = seedUser(false, "Willy", "Tybur", "willy@mail.com", "topsecret", us)
	users[3] = seedUser(false, "Karl", "Fritz", "karl@mail.com", "topsecret", us)

	var bookings [9]string
	bookings[0] = seedBooking(users[1], hotels[0], rooms[0], 20, 25, bs)
	bookings[1] = seedBooking(users[1], hotels[1], rooms[4], 26, 30, bs)
	bookings[2] = seedBooking(users[1], hotels[2], rooms[8], 31, 35, bs)
	bookings[3] = seedBooking(users[2], hotels[0], rooms[1], 20, 25, bs)
	bookings[4] = seedBooking(users[2], hotels[1], rooms[5], 26, 30, bs)
	bookings[5] = seedBooking(users[2], hotels[2], rooms[9], 31, 35, bs)
	bookings[6] = seedBooking(users[3], hotels[0], rooms[0], 15, 19, bs)
	bookings[7] = seedBooking(users[3], hotels[1], rooms[4], 20, 25, bs)
	bookings[8] = seedBooking(users[3], hotels[2], rooms[8], 26, 30, bs)

	fmt.Println("seeding the Database")
}

func initialization() (us db.UserStore, hs db.HotelStore, rs db.RoomStore, bs db.BookingStore) {
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
	if err := client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	bookingStore := db.NewMongoBookingStore(client, db.DBNAME)
	return userStore, hotelStore, roomStore, bookingStore
}

func seedBooking(userID, hotelID, roomID string, fromDate, toDate int, bs db.BookingStore) (bookingID string) {
	params := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 24 * time.Duration(toDate)),
		ToDate:   time.Now().Add(time.Hour * 24 * time.Duration(fromDate)),
	}
	booking, err := types.NewBookingFromParams(params, userID, hotelID, roomID)
	if err != nil {
		log.Fatal(err)
	}
	booking, err = bs.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return booking.ID
}

func seedRoom(hotelID string, size types.RoomSize, price float64, rs db.RoomStore) (roomID string) {
	params := types.CreateRoomParams{
		Size:  size,
		Price: price,
	}
	room := types.NewRoomFromParams(params)
	room.HotelID = hotelID
	room, err := rs.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return room.ID
}

func seedHotel(name, location string, rating int, hs db.HotelStore) (hotelID string) {
	ctx := context.Background()
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
		Rating:   rating,
	}
	insertedHotel, err := hs.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel.ID
}
func seedUser(isAdmin bool, fname, lname, email, password string, us db.UserStore) (userID string) {
	ctx := context.Background()
	userParams := types.CreateUserParams{
		Firstname: fname,
		Lastname:  lname,
		Email:     email,
		Password:  password,
	}
	user, err := types.NewUserFromParams(userParams)
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	user, err = us.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	return user.ID
}
