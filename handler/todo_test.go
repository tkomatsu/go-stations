package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

func TestCreateSuccessFull(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ts := httptest.NewServer(handler.NewTODOHandler(service.NewTODOService(todoDB)))
	defer ts.Close()

	reqBody := model.CreateTODORequest{
		Subject:     "test",
		Description: "this is test",
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(reqBody); err != nil {
		t.Fatal(err)
	}

	res, err := http.Post(ts.URL, "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Log("Incorrect response status")
	}

	var resBody model.CreateTODOResponse
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&resBody); err != nil {
		t.Fatal(err)
	}

	if resBody.TODO == nil {
		t.Fatal("TODO empty")
	}
	if resBody.TODO.Subject != reqBody.Subject {
		t.Log("Incorrect handling Subject")
	}
	if resBody.TODO.Description != reqBody.Description {
		t.Log("Incorrect handling description")
	}
}

func TestCreateSuccessOnlySubject(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ts := httptest.NewServer(handler.NewTODOHandler(service.NewTODOService(todoDB)))
	defer ts.Close()

	reqBody := model.CreateTODORequest{
		Subject: "test",
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(reqBody); err != nil {
		t.Fatal(err)
	}

	res, err := http.Post(ts.URL, "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Log("Incorrect response status")
	}

	var resBody model.CreateTODOResponse
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&resBody); err != nil {
		t.Fatal(err)
	}

	if resBody.TODO == nil {
		t.Fatal("TODO empty")
	}
	if resBody.TODO.Subject != reqBody.Subject {
		t.Log("Incorrect handling Subject")
	}
	if resBody.TODO.Description != reqBody.Description {
		t.Log("Incorrect handling description")
	}
}

func TestCreateFailGet(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ts := httptest.NewServer(handler.NewTODOHandler(service.NewTODOService(todoDB)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Log("Request with GET failed")
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Logf("Incorrect status code: %v", res.StatusCode)
	}
}

func TestCreateFailEmptySubject(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ts := httptest.NewServer(handler.NewTODOHandler(service.NewTODOService(todoDB)))
	defer ts.Close()

	reqBody := model.CreateTODORequest{
		Description: "this is test",
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(reqBody); err != nil {
		t.Fatal(err)
	}

	res, err := http.Post(ts.URL, "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Logf("Incorrect response status: %v", res.StatusCode)
	}
}
