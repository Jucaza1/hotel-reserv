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

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	GetRooms(ctx context.Context, hotelID string) ([]*types.Room, error)
	GetRoom(ctx context.Context, roomID string) (*types.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	DeleteRoomsByHotel(ctx context.Context, id string) error

	Dropper
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

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping room collection")
	return s.coll.Drop(ctx)
}
func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, types.ErrInternal(err)
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
	cur, err := s.coll.Find(ctx, bson.M{"hotelID": hotelID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound(err)
		}
		return nil, types.ErrInternal(err)

	}
	var rooms []*types.Room
	if err = cur.All(ctx, &rooms); err != nil {
		return nil, types.ErrInternal(err)
	}
	return rooms, nil
}

func (s *MongoRoomStore) GetRoom(ctx context.Context, roomID string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, types.ErrInvalidID(err)
	}
	var room types.Room
	if err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&room); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound(err)
		}
		return nil, types.ErrInternal(err)
	}
	return &room, nil
}
func (s *MongoRoomStore) DeleteRoom(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID(err)
	}
	var room types.Room
	if err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&room); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return types.ErrInternal(err)
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return types.ErrInternal(err)
	}
	if err := s.HotelStore.DeleteHotelRoom(ctx, room.HotelID, id); err != nil {
		return err
	}
	return nil
}

func (s *MongoRoomStore) DeleteRoomsByHotel(ctx context.Context, id string) error {
	_, err := s.coll.DeleteMany(ctx, bson.M{"hotelID": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return types.ErrInternal(err)
	}
	return nil
}
