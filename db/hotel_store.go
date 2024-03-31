package db

import (
	"context"
	"fmt"

	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	GetHotels(context.Context) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotelRooms(context.Context, string, string) error
	UpdateHotel(context.Context, string, map[string]string) error
}
type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}

func (s *MongoHotelStore) UpdateHotelRooms(ctx context.Context, id string, updateRoom string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"rooms": updateRoom}}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, id string, updateMap map[string]string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	//REVISAR ESTO HACER MEJOR PARSEANDO STRUCT

	validUpdate := map[string]string{}
	if updateMap["name"] != "" {
		validUpdate["name"] = updateMap["name"]
	}
	if updateMap["location"] != "" {
		validUpdate["location"] = updateMap["location"]
	}
	if updateMap["name"] != "" && updateMap["location"] != "" {
		return fmt.Errorf("no valid update parameters for hotel")
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": validUpdate}
	s.coll.UpdateOne(ctx, filter, update)
	return nil
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return hotel, nil
}
func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	res, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotel []*types.Hotel
	if err := res.All(ctx, &hotel); err != nil {
		return []*types.Hotel{}, nil
	}
	return hotel, nil
}
func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}
