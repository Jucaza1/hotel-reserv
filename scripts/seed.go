package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	us, hs, rs := initialization()
	seedHotel("hotelA", "Spain", 5, hs, rs)
	seedHotel("hotelB", "France", 4, hs, rs)
	seedHotel("hotelC", "Portugal", 3, hs, rs)

	seedUser("james", "oak", "james@mail.com", "mysecretpassword", us)

	fmt.Println("seeding the Database")
}

func initialization() (us db.UserStore, hs db.HotelStore, rs db.RoomStore) {
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
	return userStore, hotelStore, roomStore
}

func seedHotel(name string, location string, rating int, hs db.HotelStore, rs db.RoomStore) {
	ctx := context.TODO()
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  types.Small,
			Price: 99.9,
		},
		{
			Size:  types.Normal,
			Price: 119.9,
		},
		{
			Size:  types.Large,
			Price: 129.9,
		},
		{
			Size:  types.Extra,
			Price: 199.9,
		},
	}

	insertedHotel, err := hs.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := rs.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

	}

}
func seedUser(fname string, lname string, email string, password string, us db.UserStore) {
	ctx := context.TODO()
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
	_, err = us.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}
