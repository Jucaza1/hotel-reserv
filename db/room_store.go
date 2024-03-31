package db

import (
	"context"

	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, string) ([]*types.Room, error)
}
type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, dbname string, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID).Hex()

	//update hotel with room IDs slice
	updateRoom := room.ID

	if err := s.HotelStore.UpdateHotelRooms(ctx, room.HotelID, updateRoom); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, hotelID string) ([]*types.Room, error) {
	res, err := s.coll.Find(ctx, bson.M{"hotelID": hotelID})
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err = res.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
