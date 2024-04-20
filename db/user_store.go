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

const userColl = "users"

type UserStore interface {
	UpdateUser(ctx context.Context, id string, updateValid map[string]string) error
	DeleteUser(ctx context.Context, id string) error
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)

	Dropper
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}
func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, updateValid map[string]string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID(err)
	}
	filter := bson.M{"_id": oid}
	update := bson.D{{Key: "$set", Value: updateValid}}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return types.ErrNotFound(err)
		}
		return types.ErrInternal(err)
	}
	return nil
}
func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
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
func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, types.ErrInternal(err)
	}
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, types.ErrInvalidID(err)
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound(err)
		}
		return nil, types.ErrInternal(err)
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound(err)
		}
		return nil, types.ErrInternal(err)
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	res, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*types.User{}, nil
		}
		return nil, types.ErrInternal(err)
	}
	var users []*types.User
	if err := res.All(ctx, &users); err != nil {
		return []*types.User{}, nil
	}
	return users, nil
}
