package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jucaza1/hotel-reserv/types"
)

func TestHandleAuthenticateSuccess(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	app := NewFiberAppCentralErr()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	userParams := types.CreateUserParams{
		Firstname: "testName",
		Lastname:  "testLast",
		Email:     "test@foo.com",
		Password:  "secretpasstest",
	}
	insertedUser, err := types.NewUserFromParams(userParams)
	if err != nil {
		t.Error(err)
	}
	_, err = tdb.UserStore.InsertUser(context.Background(), insertedUser)
	if err != nil {
		t.Error(err)
	}
	authParams := AuthParams{
		Email:   userParams.Email,
		Pasword: userParams.Password,
	}
	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected the response status code to be %d but got %d", http.StatusNoContent, resp.StatusCode)
	}
	token := resp.Header.Get("X-Authorization")
	if token == "" {
		t.Errorf("expected token to be found in headers")
	}

}
func TestHandleAuthenticateWrongPassword(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	app := NewFiberAppCentralErr()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	userParams := types.CreateUserParams{
		Firstname: "testName",
		Lastname:  "testLast",
		Email:     "test@foo.com",
		Password:  "secretpasstest",
	}
	insertedUser, err := types.NewUserFromParams(userParams)
	if err != nil {
		t.Error(err)
	}
	_, err = tdb.UserStore.InsertUser(context.Background(), insertedUser)
	if err != nil {
		t.Error(err)
	}
	authParams := AuthParams{
		Email:   userParams.Email,
		Pasword: "wrongpassword",
	}
	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode == http.StatusNoContent {
		t.Errorf("expected the response status code to NOT be %d", http.StatusNoContent)
	}
	token := resp.Header.Get("X-Authorization")
	if token != "" {
		t.Errorf("expected token not to be found in headers")
	}

}
func TestHandleAuthenticateWrongEmail(t *testing.T) {
	tdb := userSetup(t)
	defer tdb.userTeardown(t)
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	app := NewFiberAppCentralErr()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	userParams := types.CreateUserParams{
		Firstname: "testName",
		Lastname:  "testLast",
		Email:     "test@foo.com",
		Password:  "secretpasstest",
	}
	insertedUser, err := types.NewUserFromParams(userParams)
	if err != nil {
		t.Error(err)
	}
	_, err = tdb.UserStore.InsertUser(context.Background(), insertedUser)
	if err != nil {
		t.Error(err)
	}
	authParams := AuthParams{
		Email:   "wrong@mail.com",
		Pasword: userParams.Password,
	}
	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected the response status code to be %d but got %d", http.StatusUnauthorized, resp.StatusCode)
	}
	token := resp.Header.Get("X-Authorization")
	if token != "" {
		t.Errorf("expected token not to be found in headers")
	}
}
