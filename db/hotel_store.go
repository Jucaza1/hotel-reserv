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

const hotelColl = "hotels"

type HotelStore interface {
	UpdateHotelRooms(ctx context.Context, id, updateRoom string) error
	UpdateHotel(ctx context.Context, id string, validUpdate map[string]any) error
	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	GetHotels(ctx context.Context) ([]*types.Hotel, error)
	GetHotelByID(ctx context.Context, id string) (*types.Hotel, error)
	DeleteHotel(ctx context.Context, id string) error
	DeleteHotelRoom(ctx context.Context, hotelID, roomID string) error

	Dropper
}
type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping hotel collection")
	return s.coll.Drop(ctx)
}

func (s *MongoHotelStore) UpdateHotelRooms(ctx context.Context, id, updateRoom string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID(err)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"rooms": updateRoom}}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return types.ErrNotFound(err)
		}
		return types.ErrInternal(err)
	}
	return nil
}
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, id string, validUpdate map[string]any) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID(err)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": validUpdate}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return types.ErrNotFound(err)
		}
		return types.ErrInternal(err)
	}
	return nil
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, types.ErrInternal(err)
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return hotel, nil
}
func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.Hotel{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var hotel []*types.Hotel
	if err := cur.All(ctx, &hotel); err != nil {
		return []*types.Hotel{}, nil
	}
	return hotel, nil
}
func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, types.ErrInvalidID(err)
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound(err)
		}
		return nil, types.ErrInternal(err)
	}
	return &hotel, nil
}
func (s *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID(err)
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}
	return nil
}
func (s *MongoHotelStore) DeleteHotelRoom(ctx context.Context, hotelID, roomID string) error {
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return types.ErrInvalidID(err)
	}
	if _, err := s.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$pull": bson.M{"rooms": roomID}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return types.ErrInternal(err)
	}
	return nil
}
