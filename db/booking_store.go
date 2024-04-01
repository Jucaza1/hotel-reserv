package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, string) ([]*types.Booking, error)
	GetBookingsByHotel(context.Context, string) ([]*types.Booking, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client, dbname string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	filter := bson.D{
		{"roomID", booking.RoomID},
		{"$or", bson.A{
			bson.D{{"fromDate", bson.D{{"$gte", booking.FromDate}}}, {"fromDate", bson.D{{"$lte", booking.ToDate}}}},
			bson.D{{"fromDate", bson.D{{"$lte", booking.FromDate}}}, {"toDate", bson.D{{"$gte", booking.ToDate}}}},
			bson.D{{"toDate", bson.D{{"$gte", booking.FromDate}}}, {"toDate", bson.D{{"$lte", booking.ToDate}}}},
		}},
	}
	res, err := s.coll.Find(ctx, filter)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	var bookings []*types.Booking
	if !errors.Is(err, mongo.ErrNoDocuments) {
		if err = res.All(ctx, &bookings); err != nil {
			return nil, err
		}
	}
	if len(bookings) > 0 {
		return nil, fmt.Errorf("unavailable date")
	}
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID).Hex()
	return booking, nil
}
func (s *MongoBookingStore) GetBookings(ctx context.Context, roomID string) ([]*types.Booking, error) {
	res, err := s.coll.Find(ctx, bson.M{"roomID": roomID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}
	var bookings []*types.Booking
	if err = res.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingsByHotel(ctx context.Context, hotelID string) ([]*types.Booking, error) {
	res, err := s.coll.Find(ctx, bson.M{"hotelID": hotelID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}
	var bookings []*types.Booking
	if err = res.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
