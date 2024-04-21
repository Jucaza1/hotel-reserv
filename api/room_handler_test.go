package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roomTestDB struct {
	db.RoomStore
	db.HotelStore
}

func (tdb *roomTestDB) roomTeardown(t *testing.T) {
	if err := tdb.RoomStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if err := tdb.HotelStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func roomSetup(t *testing.T) *roomTestDB {
	injectENV(t)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}
	return &roomTestDB{
		HotelStore: db.NewMongoHotelStore(client, db.TestDBNAME),
		RoomStore:  db.NewMongoRoomStore(client, db.TestDBNAME, db.NewMongoHotelStore(client, db.TestDBNAME)),
	}
}
func seedTestHotel(t *testing.T, tdb db.HotelStore) (hotelID string) {
	params := types.CreateHotelParams{
		Name:     "hoteltest",
		Location: "landtest",
		Rating:   4,
	}
	insertedHotel, err := types.NewHotelFromParams(params)
	if err != nil {
		t.Error(err)
	}
	insertedHotel, err = tdb.InsertHotel(context.Background(), insertedHotel)
	if err != nil {
		t.Error(err)
	}
	return insertedHotel.ID
}

func TestHandleGetRooms(t *testing.T) {
	tdb := roomSetup(t)
	defer tdb.roomTeardown(t)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	app := NewFiberAppCentralErr()
	roomHandler := NewRoomHandler(tdb.RoomStore, tdb.HotelStore)
	app.Get("/hotels/:hid/rooms", roomHandler.HandleGetRooms)
	params := [2]types.CreateRoomParams{
		{
			Size:  types.Normal,
			Price: 100,
		}, {
			Size:  types.Large,
			Price: 110,
		},
	}
	insertedRooms := [2]*types.Room{}
	var err error
	for i, param := range params {
		insertedRooms[i] = types.NewRoomFromParams(param)
		insertedRooms[i].HotelID = hotelID
		insertedRooms[i], err = tdb.RoomStore.InsertRoom(context.Background(), insertedRooms[i])
		if err != nil {
			t.Error(err)
		}
	}
	reqUri := fmt.Sprintf("/hotels/%s/rooms", hotelID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var rooms [2]types.Room
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		t.Error(err)
	}
	for i, room := range rooms {
		compareRoomWithID(t, insertedRooms[i], &room)
	}
}

func TestHandleGetRoomByID(t *testing.T) {
	tdb := roomSetup(t)
	defer tdb.roomTeardown(t)
	hotelID := seedTestHotel(t, tdb.HotelStore)
	app := NewFiberAppCentralErr()
	roomHandler := NewRoomHandler(tdb.RoomStore, tdb.HotelStore)
	app.Get("/hotels/:hid/rooms/:id", roomHandler.HandleGetRoomByID)
	params := [2]types.CreateRoomParams{
		{
			Size:  types.Normal,
			Price: 100,
		}, {
			Size:  types.Large,
			Price: 110,
		},
	}
	insertedRooms := [2]*types.Room{}
	var err error
	for i, param := range params {
		insertedRooms[i] = types.NewRoomFromParams(param)
		insertedRooms[i].HotelID = hotelID
		insertedRooms[i], err = tdb.RoomStore.InsertRoom(context.Background(), insertedRooms[i])
		if err != nil {
			t.Error(err)
		}
	}
	reqUri := fmt.Sprintf("/hotels/%s/rooms/%s", hotelID, insertedRooms[1].ID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var room types.Room
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		t.Error(err)
	}
	compareRoomWithID(t, insertedRooms[1], &room)
}

func TestHandleDeleteRoom(t *testing.T) {
	tdb := roomSetup(t)
	defer tdb.roomTeardown(t)
	hotelID := seedTestHotel(t, tdb.HotelStore)
	app := NewFiberAppCentralErr()
	roomHandler := NewRoomHandler(tdb.RoomStore, tdb.HotelStore)
	app.Delete("/hotels/:hid/rooms/:id", roomHandler.HandleDeleteRoom)
	params := [2]types.CreateRoomParams{
		{
			Size:  types.Normal,
			Price: 100,
		}, {
			Size:  types.Large,
			Price: 110,
		},
	}
	insertedRooms := [2]*types.Room{}
	var err error
	for i, param := range params {
		insertedRooms[i] = types.NewRoomFromParams(param)
		insertedRooms[i].HotelID = hotelID
		insertedRooms[i], err = tdb.RoomStore.InsertRoom(context.Background(), insertedRooms[i])
		if err != nil {
			t.Error(err)
		}
	}
	reqUri := fmt.Sprintf("/hotels/%s/rooms/%s", hotelID, insertedRooms[1].ID)
	req := httptest.NewRequest("DELETE", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var msg types.MsgDeleted
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		t.Error(err)
	}
	if msg.Deleted != insertedRooms[1].ID {
		t.Errorf("expected deleted room id %s but got %s", insertedRooms[1].ID, msg.Deleted)
	}
	room, err := tdb.GetRoom(context.Background(), insertedRooms[1].ID)
	if err != nil {
		errst, ok := err.(types.ErrorSt)
		if ok && (errst.Status != http.StatusNotFound) {
			t.Error(errst)
		}
	}
	if room != nil {
		t.Errorf("the test room remains in the databse")
	}
}

func TestHandlePostRoom(t *testing.T) {
	tdb := roomSetup(t)
	defer tdb.roomTeardown(t)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	app := NewFiberAppCentralErr()
	roomHandler := NewRoomHandler(tdb.RoomStore, tdb.HotelStore)
	app.Post("/hotels/:hid/rooms", roomHandler.HandlePostRoom)

	params := types.CreateRoomParams{
		Size:  types.Normal,
		Price: 100,
	}
	b, _ := json.Marshal(params)
	reqUri := fmt.Sprintf("/hotels/%s/rooms", hotelID)
	req := httptest.NewRequest("POST", reqUri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}

	insertedRoom := types.Room{}
	expedted := types.Room{
		HotelID: hotelID,
		Size:    params.Size,
		Price:   params.Price,
	}
	json.NewDecoder(resp.Body).Decode(&insertedRoom)
	compareRoom(t, &expedted, &insertedRoom)

}

func compareRoomWithID(t *testing.T, expected *types.Room, have *types.Room) {
	if len(have.ID) == 0 || have.ID != expected.ID {
		t.Errorf("expected room id %s but got %s", expected.ID, have.ID)
	}
	if have.HotelID != expected.HotelID {
		t.Errorf("expected room's hotelID %s but got %s", expected.HotelID, have.HotelID)
	}
	if have.Size != expected.Size {
		t.Errorf("expected room size %d but got %d", expected.Size, have.Size)
	}
	if have.Price != expected.Price {
		t.Errorf("expected room price %f but got %f", expected.Price, have.Price)
	}
}
func compareRoom(t *testing.T, expected *types.Room, have *types.Room) {
	if len(have.ID) == 0 {
		t.Errorf("expected room id to be set")
	}
	if have.HotelID != expected.HotelID {
		t.Errorf("expected room's hotelID %s but got %s", expected.HotelID, have.HotelID)
	}
	if have.Size != expected.Size {
		t.Errorf("expected room size %d but got %d", expected.Size, have.Size)
	}
	if have.Price != expected.Price {
		t.Errorf("expected room price %f but got %f", expected.Price, have.Price)
	}
}
