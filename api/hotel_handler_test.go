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

type hotelTestDB struct {
	db.HotelStore
}

func (tdb *hotelTestDB) hotelTeardown(t *testing.T) {
	if err := tdb.HotelStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func hotelSetup(t *testing.T) *hotelTestDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}
	return &hotelTestDB{
		HotelStore: db.NewMongoHotelStore(client, db.TestDBNAME),
	}
}
func TestPostHotel(t *testing.T) {
	tdb := hotelSetup(t)
	defer tdb.hotelTeardown(t)

	app := NewFiberAppCentralErr()
	hotelHandler := NewHotelHandler(tdb.HotelStore)
	app.Post("/", hotelHandler.HandlePostHotel)

	params := types.CreateHotelParams{
		Name:     "hoteltest",
		Location: "landtest",
		Rating:   4,
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var hotel types.Hotel
	json.NewDecoder(resp.Body).Decode(&hotel)
	compareHotel(t, &params, &hotel)
}

func TestGetHotel(t *testing.T) {
	tdb := hotelSetup(t)
	defer tdb.hotelTeardown(t)

	app := NewFiberAppCentralErr()
	hotelHandler := NewHotelHandler(tdb.HotelStore)
	app.Get("/:id", hotelHandler.HandleGetHotel)

	params := types.CreateHotelParams{
		Name:     "hoteltest",
		Location: "landtest",
		Rating:   4,
	}
	insertedHotel := types.NewHotelFromParams(params)
	insertedHotel, err := tdb.HotelStore.InsertHotel(context.Background(), insertedHotel)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/%s", insertedHotel.ID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var hotel types.Hotel
	if err := json.NewDecoder(resp.Body).Decode(&hotel); err != nil {
		t.Error(err)
	}
	compareHotelWithID(t, insertedHotel, &hotel)
}

func TestGetHotels(t *testing.T) {
	tdb := hotelSetup(t)
	defer tdb.hotelTeardown(t)

	app := NewFiberAppCentralErr()
	hotelHandler := NewHotelHandler(tdb.HotelStore)
	app.Get("/", hotelHandler.HandleGetHotels)

	params := [2]types.CreateHotelParams{
		{
			Name:     "hoteltest",
			Location: "landtest",
			Rating:   4,
		}, {
			Name:     "hoteltest2",
			Location: "landtest2",
			Rating:   5,
		},
	}
	insertedHotels := [2]*types.Hotel{}
	var err error
	for i, param := range params {
		insertedHotels[i] = types.NewHotelFromParams(param)
		insertedHotels[i], err = tdb.HotelStore.InsertHotel(context.Background(), insertedHotels[i])
		if err != nil {
			t.Error(err)
		}
	}
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var hotels [2]types.Hotel
	if err := json.NewDecoder(resp.Body).Decode(&hotels); err != nil {
		t.Error(err)
	}
	for i, hotel := range hotels {
		compareHotelWithID(t, insertedHotels[i], &hotel)
	}

}
func TestDeleteHotel(t *testing.T) {
	tdb := roomSetup(t)
	defer tdb.roomTeardown(t)

	app := NewFiberAppCentralErr()
	hotelHandler := NewHotelHandler(tdb.HotelStore)
	roomHandler := NewRoomHandler(tdb.RoomStore, tdb.HotelStore)
	app.Delete("/:id", hotelHandler.HandleDeleteHotel, roomHandler.HandleDeleteRoomsByHotel)

	params := types.CreateHotelParams{
		Name:     "hoteltest",
		Location: "landtest",
		Rating:   4,
	}

	insertedHotel := types.NewHotelFromParams(params)
	insertedHotel, err := tdb.HotelStore.InsertHotel(context.Background(), insertedHotel)
	if err != nil {
		t.Error(err)
	}
	roomParams := [2]types.CreateRoomParams{
		{
			Size:  types.Normal,
			Price: 100,
		}, {
			Size:  types.Large,
			Price: 110,
		},
	}
	insertedRooms := [2]*types.Room{}
	for i, param := range roomParams {
		insertedRooms[i] = types.NewRoomFromParams(param)
		insertedRooms[i].HotelID = insertedHotel.ID
		insertedRooms[i], err = tdb.RoomStore.InsertRoom(context.Background(), insertedRooms[i])
		if err != nil {
			t.Error(err)
		}
	}
	insertedHotel, err = tdb.HotelStore.GetHotelByID(context.Background(), insertedHotel.ID)
	if err != nil {
		t.Error(err)
	}

	reqUri := fmt.Sprintf("/%s", insertedHotel.ID)
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
	if msg.Deleted != insertedHotel.ID {
		t.Errorf("expected deleted hotel id %s but got %s", insertedHotel.ID, msg.Deleted)
	}
	hotel, err := tdb.GetHotelByID(context.Background(), insertedHotel.ID)
	if err != nil {
		errst, ok := err.(types.ErrorSt)
		if ok && (errst.Status != http.StatusNotFound) {
			t.Error(errst)
		}
	}
	if hotel != nil {
		t.Errorf("the test hotel remains in the database")
	}
	rooms, err := tdb.GetRooms(context.Background(), insertedHotel.ID)
	if err != nil {
		errst, ok := err.(types.ErrorSt)
		if ok && (errst.Status != http.StatusNotFound) {
			t.Error(errst)
		}
	}
	if len(rooms) != 0 {
		t.Errorf("the test rooms remains in the database")
	}
}

func TestPatchHotel(t *testing.T) {
	tdb := hotelSetup(t)
	defer tdb.hotelTeardown(t)

	app := NewFiberAppCentralErr()

	hotelHandler := NewHotelHandler(tdb.HotelStore)
	app.Patch("/:id", hotelHandler.HandlePatchHotel)
	params := types.CreateHotelParams{
		Name:     "test1",
		Location: "testLand1",
		Rating:   3,
	}
	instertedHotel := types.NewHotelFromParams(params)
	insertedHotel, err := tdb.HotelStore.InsertHotel(context.Background(), instertedHotel)
	if err != nil {
		t.Error(err)
	}
	update := types.UpdateHotel{
		Name:     "test2",
		Location: "testLand2",
	}
	b, _ := json.Marshal(update)
	reqUri := fmt.Sprintf("/%s", insertedHotel.ID)
	req := httptest.NewRequest("PATCH", reqUri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}

	var msg types.MsgUpdated
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		t.Error(err)
	}
	if msg.Updated != insertedHotel.ID {
		t.Errorf("expected updated hotel id %s but got %s", insertedHotel.ID, msg.Updated)
	}
	hotel, err := tdb.GetHotelByID(context.Background(), insertedHotel.ID)
	if err != nil {
		t.Error(err)
	}
	hotel2 := types.Hotel{
		ID:       insertedHotel.ID,
		Name:     update.Name,
		Location: update.Location,
		Rating:   insertedHotel.Rating,
	}
	compareHotelWithID(t, &hotel2, hotel)
}

func compareHotel(t *testing.T, expected *types.CreateHotelParams, have *types.Hotel) {
	if len(have.ID) == 0 {
		t.Errorf("expected a hotel id to be set")
	}
	if have.Name != expected.Name {
		t.Errorf("expected hotel name %s but got %s", expected.Name, have.Name)
	}
	if have.Location != expected.Location {
		t.Errorf("expected hotel location %s but got %s", expected.Location, have.Location)
	}
	if have.Rating != expected.Rating {
		t.Errorf("expected hotel rating %d but got %d", expected.Rating, have.Rating)
	}
}
func compareHotelWithID(t *testing.T, expected *types.Hotel, have *types.Hotel) {
	if len(have.ID) == 0 || have.ID != expected.ID {
		t.Errorf("expected hotel id %s but got %s", expected.ID, have.ID)
	}
	if have.Name != expected.Name {
		t.Errorf("expected hotel name %s but got %s", expected.Name, have.Name)
	}
	if have.Location != expected.Location {
		t.Errorf("expected hotel location %s but got %s", expected.Location, have.Location)
	}
	if have.Rating != expected.Rating {
		t.Errorf("expected hotel rating %d but got %d", expected.Rating, have.Rating)
	}
}
