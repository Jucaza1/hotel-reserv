package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userTestDB struct {
	db.UserStore
}

func (tdb *userTestDB) userTeardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func userSetup(t *testing.T) *userTestDB {
	injectENV(t)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		t.Errorf("dbURI %s", db.DBURI)
		t.Fatal(err)
	}
	return &userTestDB{
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)

	app := NewFiberAppCentralErr()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Firstname: "testName",
		Lastname:  "testLast",
		Email:     "test@foo.com",
		Password:  "12345678",
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
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	compareUser(t, &params, &user)
}

func TestGetUser(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)

	app := NewFiberAppCentralErr()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/:id", userHandler.HandleGetUser)

	params := types.CreateUserParams{
		Firstname: "testName",
		Lastname:  "testLast",
		Email:     "test@foo.com",
		Password:  "12345678",
	}
	insertedUser, err := types.NewUserFromParams(params)
	if err != nil {
		t.Error(err)
	}
	insertedUser, err = tdb.UserStore.InsertUser(context.Background(), insertedUser)
	if err != nil {
		t.Error(err)
	}
	reqUri := fmt.Sprintf("/%s", insertedUser.ID)
	req := httptest.NewRequest("GET", reqUri, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code expected 200 but got %d", resp.StatusCode)
	}
	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}
	compareUserWithID(t, insertedUser, &user)
}

func TestGetUsers(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)

	app := NewFiberAppCentralErr()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/", userHandler.HandleGetUsers)

	params := [2]types.CreateUserParams{
		{
			Firstname: "testName1",
			Lastname:  "testLast1",
			Email:     "test@foo.com",
			Password:  "123456781",
		}, {
			Firstname: "testName2",
			Lastname:  "testLast2",
			Email:     "test2@foo.com",
			Password:  "123456782",
		},
	}
	insertedUser := [2]*types.User{}
	var err error
	for i, param := range params {
		insertedUser[i], err = types.NewUserFromParams(param)
		if err != nil {
			t.Error(err)
		}

		insertedUser[i], err = tdb.UserStore.InsertUser(context.Background(), insertedUser[i])
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
	var users [2]types.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Error(err)
	}
	for i, u := range users {
		compareUserWithID(t, insertedUser[i], &u)
	}

}

func TestDeleteUser(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)

	app := NewFiberAppCentralErr()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Delete("/:id", userHandler.HandleDeleteUser)

	params := types.CreateUserParams{
		Firstname: "testName1",
		Lastname:  "testLast1",
		Email:     "test@foo.com",
		Password:  "123456781",
	}

	insertedUser, err := types.NewUserFromParams(params)
	if err != nil {
		t.Error(err)
	}

	insertedUser, err = tdb.UserStore.InsertUser(context.Background(), insertedUser)
	if err != nil {
		t.Error(err)
	}

	reqUri := fmt.Sprintf("/%s", insertedUser.ID)
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
	if msg.Deleted != insertedUser.ID {
		t.Errorf("expected deleted user id %s but got %s", insertedUser.ID, msg.Deleted)
	}
	user, err := tdb.GetUserByID(context.Background(), insertedUser.ID)
	if err != nil && err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}
	if user != nil {
		t.Errorf("the test user remains in the databse")
	}
}

func TestPatchUser(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)

	app := NewFiberAppCentralErr()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Patch("/:id", userHandler.HandlePatchUser)

	params := types.CreateUserParams{
		Firstname: "testName",
		Lastname:  "testLast",
		Email:     "test@foo.com",
		Password:  "12345678",
	}
	insertedUser, err := types.NewUserFromParams(params)
	if err != nil {
		t.Error(err)
	}
	insertedUser, err = tdb.UserStore.InsertUser(context.Background(), insertedUser)
	if err != nil {
		t.Error(err)
	}
	updateParams := types.UpdateUser{
		Firstname: "testName2",
		Lastname:  "testLast2",
		Email:     "test2@foo.com",
	}
	b, _ := json.Marshal(updateParams)
	reqUri := fmt.Sprintf("/%s", insertedUser.ID)
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
	if msg.Updated != insertedUser.ID {
		t.Errorf("expected updated user id %s but got %s", insertedUser.ID, msg.Updated)
	}
	user, err := tdb.GetUserByID(context.Background(), insertedUser.ID)
	if err != nil {
		t.Error(err)
	}
	user2 := types.User{
		ID:        insertedUser.ID,
		Firstname: updateParams.Firstname,
		Lastname:  updateParams.Lastname,
		Email:     updateParams.Email,
	}
	compareUserWithID(t, &user2, user)
}

func compareUser(t *testing.T, expected *types.CreateUserParams, have *types.User) {
	if len(have.ID) == 0 {
		t.Errorf("expected a user id to be set")
	}
	if have.Firstname != expected.Firstname {
		t.Errorf("expected first name %s but got %s", expected.Firstname, have.Firstname)
	}
	if have.Lastname != expected.Lastname {
		t.Errorf("expected last name %s but got %s", expected.Lastname, have.Lastname)
	}
	if have.Email != expected.Email {
		t.Errorf("expected email %s but got %s", expected.Email, have.Email)
	}
}
func compareUserWithID(t *testing.T, expected *types.User, have *types.User) {
	if len(have.ID) == 0 || have.ID != expected.ID {
		t.Errorf("expected user id %s but got %s", expected.ID, have.ID)
	}
	if have.Firstname != expected.Firstname {
		t.Errorf("expected first name %s but got %s", expected.Firstname, have.Firstname)
	}
	if have.Lastname != expected.Lastname {
		t.Errorf("expected last name %s but got %s", expected.Lastname, have.Lastname)
	}
	if have.Email != expected.Email {
		t.Errorf("expected email %s but got %s", expected.Email, have.Email)
	}
}
