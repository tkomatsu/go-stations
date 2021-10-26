package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/model"
)

func TestHealthz(t *testing.T) {
	ts := httptest.NewServer(handler.NewHealthzHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var resBody model.HealthzResponse
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&resBody); err != nil {
		t.Fatal(err)
	}

	if resBody.Message != "OK" {
		t.Fatal("Incorrect Response")
	}
}
