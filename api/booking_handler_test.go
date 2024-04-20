package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bookingTestDB struct {
	db.HotelStore
	db.RoomStore
	db.BookingStore
}

func (tdb *bookingTestDB) bookingTeardown(t *testing.T) {
	if err := tdb.RoomStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if err := tdb.BookingStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if err := tdb.HotelStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func bookingSetup(t *testing.T) *bookingTestDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Fatal(err)
	}
	return &bookingTestDB{
		HotelStore:   db.NewMongoHotelStore(client, db.TestDBNAME),
		RoomStore:    db.NewMongoRoomStore(client, db.TestDBNAME, db.NewMongoHotelStore(client, db.TestDBNAME)),
		BookingStore: db.NewMongoBookingStore(client, db.TestDBNAME),
	}
}

func seedTestRoom(t *testing.T, tdb db.RoomStore, hotelID string) (roomID string) {
	params := types.CreateRoomParams{
		Price: 100,
		Size:  types.Normal,
	}
	room := types.NewRoomFromParams(params)
	room.HotelID = hotelID
	insertedRoom, err := tdb.InsertRoom(context.Background(), room)
	if err != nil {
		t.Error(err)
	}
	return insertedRoom.ID
}

func TestHandleGetBookingsByRoomAsAdmin(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	userAdmin := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          true,
	}
	userNoAdmin := types.User{
		ID:               "0002",
		Firstname:        "testname2",
		Lastname:         "testlast2",
		Email:            "test2@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Get("/rooms/:id/bookings", provideContextUser(userAdmin), bookingHandler.HandleGetBookingsByRoom)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	paramsAdmin := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 150),
		ToDate:   time.Now().Add(time.Hour * 180),
	}
	var (
		err                    error
		newBookingNoAdmin      [2]*types.Booking
		insertedBookingNoAdmin [2]*types.Booking
		newBookingAdmin        *types.Booking
		insertedBookingAdmin   *types.Booking
	)

	for i, param := range params {
		newBookingNoAdmin[i], err = types.NewBookingFromParams(param, userNoAdmin.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBookingNoAdmin[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBookingNoAdmin[i])
		if err != nil {
			t.Error(err)
		}
	}
	newBookingAdmin, err = types.NewBookingFromParams(paramsAdmin, userAdmin.ID, hotelID, roomID)
	if err != nil {
		t.Error(err)
	}
	insertedBookingAdmin, err = tdb.BookingStore.InsertBooking(context.Background(), newBookingAdmin)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/rooms/%s/bookings", roomID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var have []types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}
	for i, expected := range []*types.Booking{insertedBookingNoAdmin[0], insertedBookingNoAdmin[1], insertedBookingAdmin} {
		compareBookingWithID(t, expected, &have[i])
	}
}

func TestHandleGetBookingsByRoomNoAdmin(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	userAdmin := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          true,
	}
	userNoAdmin := types.User{
		ID:               "0002",
		Firstname:        "testname2",
		Lastname:         "testlast2",
		Email:            "test2@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Get("/rooms/:id/bookings", provideContextUser(userNoAdmin), bookingHandler.HandleGetBookingsByRoom)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	paramsAdmin := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 150),
		ToDate:   time.Now().Add(time.Hour * 180),
	}
	var (
		err                    error
		newBookingNoAdmin      [2]*types.Booking
		insertedBookingNoAdmin [2]*types.Booking
		newBookingAdmin        *types.Booking
	)

	for i, param := range params {
		newBookingNoAdmin[i], err = types.NewBookingFromParams(param, userNoAdmin.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBookingNoAdmin[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBookingNoAdmin[i])
		if err != nil {
			t.Error(err)
		}
	}
	newBookingAdmin, err = types.NewBookingFromParams(paramsAdmin, userAdmin.ID, hotelID, roomID)
	if err != nil {
		t.Error(err)
	}
	_, err = tdb.BookingStore.InsertBooking(context.Background(), newBookingAdmin)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/rooms/%s/bookings", roomID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var have []types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}
	if len(have) != len(insertedBookingNoAdmin) {
		t.Errorf("expected the response to contain %d bookings but got %d", len(insertedBookingNoAdmin), len(have))
	}
	for i, expected := range insertedBookingNoAdmin {
		compareBookingWithID(t, expected, &have[i])
	}
}

func TestHandleGetBookingsByHotelAsAdmin(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	userAdmin := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          true,
	}
	userNoAdmin := types.User{
		ID:               "0002",
		Firstname:        "testname2",
		Lastname:         "testlast2",
		Email:            "test2@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Get("/hotels/:hid/bookings", provideContextUser(userAdmin), bookingHandler.HandleGetBookingsByHotel)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	paramsAdmin := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 150),
		ToDate:   time.Now().Add(time.Hour * 180),
	}
	var (
		err                    error
		newBookingNoAdmin      [2]*types.Booking
		insertedBookingNoAdmin [2]*types.Booking
		newBookingAdmin        *types.Booking
		insertedBookingAdmin   *types.Booking
	)

	for i, param := range params {
		newBookingNoAdmin[i], err = types.NewBookingFromParams(param, userNoAdmin.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBookingNoAdmin[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBookingNoAdmin[i])
		if err != nil {
			t.Error(err)
		}
	}
	newBookingAdmin, err = types.NewBookingFromParams(paramsAdmin, userAdmin.ID, hotelID, roomID)
	if err != nil {
		t.Error(err)
	}
	insertedBookingAdmin, err = tdb.BookingStore.InsertBooking(context.Background(), newBookingAdmin)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/hotels/%s/bookings", hotelID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var have []types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}
	for i, expected := range []*types.Booking{insertedBookingNoAdmin[0], insertedBookingNoAdmin[1], insertedBookingAdmin} {
		compareBookingWithID(t, expected, &have[i])
	}
}

func TestHandleGetBookingsByHotelNoAdmin(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	userAdmin := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          true,
	}
	userNoAdmin := types.User{
		ID:               "0002",
		Firstname:        "testname2",
		Lastname:         "testlast2",
		Email:            "test2@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Get("/hotels/:hid/bookings", provideContextUser(userNoAdmin), bookingHandler.HandleGetBookingsByHotel)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	paramsAdmin := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 150),
		ToDate:   time.Now().Add(time.Hour * 180),
	}
	var (
		err                    error
		newBookingNoAdmin      [2]*types.Booking
		insertedBookingNoAdmin [2]*types.Booking
		newBookingAdmin        *types.Booking
	)

	for i, param := range params {
		newBookingNoAdmin[i], err = types.NewBookingFromParams(param, userNoAdmin.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBookingNoAdmin[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBookingNoAdmin[i])
		if err != nil {
			t.Error(err)
		}
	}
	newBookingAdmin, err = types.NewBookingFromParams(paramsAdmin, userAdmin.ID, hotelID, roomID)
	if err != nil {
		t.Error(err)
	}
	_, err = tdb.BookingStore.InsertBooking(context.Background(), newBookingAdmin)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/hotels/%s/bookings", hotelID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var have []types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}
	if len(have) != len(insertedBookingNoAdmin) {
		t.Errorf("expected the response to contain %d bookings but got %d", len(insertedBookingNoAdmin), len(have))
	}
	for i, expected := range insertedBookingNoAdmin {
		compareBookingWithID(t, expected, &have[i])
	}
}

func TestHandleCancelBooking(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	user := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Patch("/bookings/:id", provideContextUser(user), bookingHandler.HandleCancelBooking)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	var (
		err             error
		newBooking      [2]*types.Booking
		insertedBooking [2]*types.Booking
	)

	for i, param := range params {
		newBooking[i], err = types.NewBookingFromParams(param, user.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBooking[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBooking[i])
		if err != nil {
			t.Error(err)
		}
	}

	reqUri := fmt.Sprintf("/bookings/%s", insertedBooking[1].ID)
	req := httptest.NewRequest("PATCH", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}

	var msg types.MsgCancelled
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		t.Error(err)
	}
	if msg.Cancelled != insertedBooking[1].ID {
		t.Errorf("expected cancelled message to be %s but got %s", insertedBooking[1].ID, msg.Cancelled)
	}
	have, err := tdb.GetBookingsByUser(context.Background(), user.ID)
	if err != nil {
		t.Error(err)
	}
	if len(have) != len(insertedBooking) {
		t.Errorf("expected the response to contain %d bookings but got %d", len(insertedBooking), len(have))
	}
	for i, expected := range insertedBooking {
		compareBookingWithID(t, expected, have[i])
	}
	if have[0].Cancelled == true {
		t.Errorf("expected first booking not to be cancelled")
	}
	if have[1].Cancelled != true {
		t.Errorf("expected second booking to be cancelled")
	}
}

func TestHandleDeleteBooking(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	user := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Delete("/bookings/:id", bookingHandler.HandleDeleteBooking)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	var (
		err             error
		newBooking      [2]*types.Booking
		insertedBooking [2]*types.Booking
	)

	for i, param := range params {
		newBooking[i], err = types.NewBookingFromParams(param, user.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBooking[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBooking[i])
		if err != nil {
			t.Error(err)
		}
	}

	reqUri := fmt.Sprintf("/bookings/%s", insertedBooking[1].ID)
	req := httptest.NewRequest("DELETE", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var msg types.MsgDeleted
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		t.Error(err)
	}
	if msg.Deleted != insertedBooking[1].ID {
		t.Errorf("expected cancelled message to be %s but got %s", insertedBooking[1].ID, msg.Deleted)
	}
	have, err := tdb.GetBookingsByUser(context.Background(), user.ID)
	if err != nil {
		t.Error(err)
	}
	if len(have) != len(insertedBooking)-1 {
		t.Errorf("expected the response to contain %d bookings but got %d", len(insertedBooking)-1, len(have))
	}
	compareBookingWithID(t, insertedBooking[0], have[0])
}

func TestHandlePostBookingsSuccess(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	user := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Post("/rooms/:id/bookings", provideContextUser(user), bookingHandler.HandlePostBooking)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	var (
		err             error
		newBooking      [2]*types.Booking
		insertedBooking [2]*types.Booking
	)
	for i, param := range params {
		newBooking[i], err = types.NewBookingFromParams(param, user.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBooking[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBooking[i])
		if err != nil {
			t.Error(err)
		}
	}
	postBookingParams := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 120),
		ToDate:   time.Now().Add(time.Hour * 150),
	}

	b, err := json.Marshal(postBookingParams)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/rooms/%s/bookings", roomID)
	req := httptest.NewRequest("POST", reqUri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var have types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}
	expected, err := types.NewBookingFromParams(postBookingParams, user.ID, hotelID, roomID)
	if err != nil {
		t.Error(err)
	}
	compareBooking(t, expected, &have)
}

func TestHandlePostBookingsFailureDate(t *testing.T) {
	tdb := bookingSetup(t)
	defer tdb.bookingTeardown(t)

	app := NewFiberAppCentralErr()
	bookingHandler := NewBookingHandler(tdb.BookingStore, tdb.RoomStore)
	user := types.User{
		ID:               "0000",
		Firstname:        "testname",
		Lastname:         "testlast",
		Email:            "test@mail.com",
		EncyptedPassword: "0",
		IsAdmin:          false,
	}
	app.Post("/rooms/:id/bookings", provideContextUser(user), bookingHandler.HandlePostBooking)

	hotelID := seedTestHotel(t, tdb.HotelStore)
	roomID := seedTestRoom(t, tdb.RoomStore, hotelID)
	params := [2]types.CreateBookingParams{
		{
			FromDate: time.Now().Add(time.Hour * 24),
			ToDate:   time.Now().Add(time.Hour * 72),
		}, {
			FromDate: time.Now().Add(time.Hour * 84),
			ToDate:   time.Now().Add(time.Hour * 108),
		},
	}
	var (
		err             error
		newBooking      [2]*types.Booking
		insertedBooking [2]*types.Booking
	)
	for i, param := range params {
		newBooking[i], err = types.NewBookingFromParams(param, user.ID, hotelID, roomID)
		if err != nil {
			t.Error(err)
		}
		insertedBooking[i], err = tdb.BookingStore.InsertBooking(context.Background(), newBooking[i])
		if err != nil {
			t.Error(err)
		}
	}
	postBookingParams := types.CreateBookingParams{
		FromDate: time.Now().Add(time.Hour * 90),
		ToDate:   time.Now().Add(time.Hour * 150),
	}

	b, err := json.Marshal(postBookingParams)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/rooms/%s/bookings", roomID)
	req := httptest.NewRequest("POST", reqUri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status code expected %d but got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
	var have types.MsgError
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}
	if have.Error != "unavailable date" {
		t.Errorf("expected error message to be: unavailable date")
	}

}

func provideContextUser(u types.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Context().SetUserValue("user", u)
		return c.Next()
	}
}
func compareBookingWithID(t *testing.T, expected *types.Booking, have *types.Booking) {
	if len(have.ID) == 0 || have.ID != expected.ID {
		t.Errorf("expected booking id %s but got %s", expected.ID, have.ID)
	}
	if have.HotelID != expected.HotelID {
		t.Errorf("expected booking hotelID %s but got %s", expected.HotelID, have.HotelID)
	}
	if have.UserID != expected.UserID {
		t.Errorf("expected booking userID %s but got %s", expected.UserID, have.UserID)
	}
	if have.FromDate.Truncate(time.Second).Unix() != expected.FromDate.Truncate(time.Second).Unix() {
		t.Errorf("expected booking fromDate %d but got %d",
			expected.FromDate.Truncate(time.Second).Unix(), have.FromDate.Truncate(time.Second).Unix())
	}
	if have.ToDate.Truncate(time.Second).Unix() != expected.ToDate.Truncate(time.Second).Unix() {
		t.Errorf("expected booking toDate %d but got %d",
			expected.ToDate.Truncate(time.Second).Unix(), have.ToDate.Truncate(time.Second).Unix())
	}
}
func compareBooking(t *testing.T, expected *types.Booking, have *types.Booking) {
	if len(have.ID) == 0 {
		t.Errorf("expected booking id to be set")
	}
	if have.HotelID != expected.HotelID {
		t.Errorf("expected booking hotelID %s but got %s", expected.HotelID, have.HotelID)
	}
	if have.UserID != expected.UserID {
		t.Errorf("expected booking userID %s but got %s", expected.UserID, have.UserID)
	}
	if have.FromDate.Truncate(time.Second).Unix() != expected.FromDate.Truncate(time.Second).Unix() {
		t.Errorf("expected booking fromDate %d but got %d",
			expected.FromDate.Truncate(time.Second).Unix(), have.FromDate.Truncate(time.Second).Unix())
	}
	if have.ToDate.Truncate(time.Second).Unix() != expected.ToDate.Truncate(time.Second).Unix() {
		t.Errorf("expected booking toDate %d but got %d",
			expected.ToDate.Truncate(time.Second).Unix(), have.ToDate.Truncate(time.Second).Unix())
	}
}
