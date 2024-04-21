package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context) ([]*types.Booking, error)
	GetBookingsByRoom(ctx context.Context, roomID string) ([]*types.Booking, error)
	GetBookingsByHotel(ctx context.Context, hotelID string) ([]*types.Booking, error)
	GetBookingsByUser(ctx context.Context, userID string) ([]*types.Booking, error)
	GetBookingsByUserAndHotel(ctx context.Context, userID, hotelID string) ([]*types.Booking, error)
	GetBookingsByUserAndRoom(ctx context.Context, userID, roomID string) ([]*types.Booking, error)
	GetBookingByID(ctx context.Context, bookingID string) (*types.Booking, error)
	CancelBooking(ctx context.Context, bookingID string) error
	DeleteBooking(ctx context.Context, bookingID string) error

	Dropper
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client, dbname string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping booking collection")
	return s.coll.Drop(ctx)
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	//check if date is available
	filter := bson.D{
		{"roomID", booking.RoomID},
		{"cancelled", false},
		{"$or", bson.A{
			bson.D{{"fromDate", bson.D{{"$gte", booking.FromDate}}}, {"fromDate", bson.D{{"$lte", booking.ToDate}}}},
			bson.D{{"fromDate", bson.D{{"$lte", booking.FromDate}}}, {"toDate", bson.D{{"$gte", booking.ToDate}}}},
			bson.D{{"toDate", bson.D{{"$gte", booking.FromDate}}}, {"toDate", bson.D{{"$lte", booking.ToDate}}}},
		}},
	}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if !errors.Is(err, mongo.ErrNoDocuments) {
		if err = cur.All(ctx, &bookings); err != nil {
			return nil, types.ErrInternal(err)
		}
	}
	if len(bookings) > 0 {
		return nil, types.ErrUnavailableDate(fmt.Errorf("unavailable date"))
	}
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, types.ErrInternal(err)
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID).Hex()
	return booking, nil
}
func (s *MongoBookingStore) GetBookings(ctx context.Context) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Booking{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, types.ErrInternal(err)
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingsByRoom(ctx context.Context, roomID string) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.M{"roomID": roomID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Booking{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, types.ErrInternal(err)
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingsByHotel(ctx context.Context, hotelID string) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.M{"hotelID": hotelID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Booking{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, types.ErrInternal(err)
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingsByUser(ctx context.Context, userID string) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Booking{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, types.ErrInternal(err)
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingsByUserAndHotel(ctx context.Context, userID, hotelID string) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.D{{"userID", userID}, {"hotelID", hotelID}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Booking{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, types.ErrInternal(err)
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingsByUserAndRoom(ctx context.Context, userID, roomID string) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.D{{"userID", userID}, {"roomID", roomID}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Booking{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, types.ErrInternal(err)
	}
	return bookings, nil
}
func (s *MongoBookingStore) GetBookingByID(ctx context.Context, bookingID string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return nil, types.ErrInvalidID(err)
	}
	var booking types.Booking
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound(err)
		}
		return nil, types.ErrInternal(err)
	}
	return &booking, nil
}
func (s *MongoBookingStore) CancelBooking(ctx context.Context, bookingID string) error {
	oid, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return types.ErrInvalidID(err)
	}
	var booking types.Booking
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return types.ErrNotFound(err)
		}
		return types.ErrInternal(err)
	}
	if booking.FromDate.Unix() <= time.Now().Unix() {
		return types.ErrCancelPastBooking(fmt.Errorf("can not cancel booking in the past"))
	}
	filter := bson.M{"_id": oid}
	update := bson.D{{"$set", bson.D{{"cancelled", true}, {"cancelledAt", time.Now()}}}}
	if _, err = s.coll.UpdateOne(ctx, filter, update); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return types.ErrNotFound(err)
		}
		return types.ErrInternal(err)
	}
	return nil
}
func (s *MongoBookingStore) DeleteBooking(ctx context.Context, bookingID string) error {
	oid, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return types.ErrInvalidID(err)
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return types.ErrInternal(err)
	}
	return nil
}
