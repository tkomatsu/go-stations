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

func TestCreate(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ts := httptest.NewServer(handler.NewTODOHandler(service.NewTODOService(todoDB)))
	defer ts.Close()

	testcase := []struct {
		name string
		req model.CreateTODORequest
		wantStatus int
	}{
		{
			name: "normal",
			req: model.CreateTODORequest{
				Subject:     "test",
				Description: "this is test",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "only subject",
			req: model.CreateTODORequest{
				Subject: "test",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "only description",
			req: model.CreateTODORequest{
				Description: "this is test",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func (t *testing.T) {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(tc.req); err != nil {
				t.Fatal(err)
			}
			res, err := http.Post(ts.URL, "application/json", &buf)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.wantStatus {
				t.Fatal("Incorrect response status")
			}
			if tc.wantStatus == http.StatusOK {
				var resBody model.CreateTODOResponse
				dec := json.NewDecoder(res.Body)
				if err := dec.Decode(&resBody); err != nil {
					t.Fatal(err)
				}
			
				if resBody.TODO == nil {
					t.Fatal("TODO empty")
				}
				if resBody.TODO.Subject != tc.req.Subject {
					t.Log("Incorrect handling Subject")
				}
				if resBody.TODO.Description != tc.req.Description {
					t.Log("Incorrect handling description")
				}
			}
		})
	}

	t.Run("not allowed method", func (t *testing.T) {
		res, err := http.Get(ts.URL)
		if err != nil {
			t.Log("Request with GET failed")
		}
	
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Logf("Incorrect status code: %v", res.StatusCode)
		}
	})
}

func TestUpdate(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ts := httptest.NewServer(handler.NewTODOHandler(service.NewTODOService(todoDB)))
	defer ts.Close()

	cli := http.DefaultClient

	testcase := []struct {
		name string
		req model.UpdateTODORequest
		wantStatus int
	}{
		{
			name: "normal",
			req: model.UpdateTODORequest{
				ID: 1,
				Subject: "hello",
				Description: "update",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "empty ID",
			req: model.UpdateTODORequest{
				Subject: "hello",
				Description: "update",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty subject",
			req: model.UpdateTODORequest{
				ID: 1,
				Description: "update",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty description",
			req: model.UpdateTODORequest{
				ID: 1,
				Subject: "hello",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			req: model.UpdateTODORequest{
				ID: 9999999,
				Subject: "hello",
				Description: "update",
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func (t *testing.T) {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(tc.req); err != nil {
				t.Fatal(err)
			}

			httpReq, err := http.NewRequest("PUT", ts.URL, &buf)
			if err != nil {
				t.Fatal(err)
			}

			res, err := cli.Do(httpReq)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.wantStatus {
				t.Fatal("Incorrect response status")
			}
			if tc.wantStatus == http.StatusOK {
				var resBody model.UpdateTODOResponse
				dec := json.NewDecoder(res.Body)
				if err := dec.Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				if resBody.TODO == nil {
					t.Fatal("TODO empty")
				}
				if resBody.TODO.ID != tc.req.ID {
					t.Log("Incorrect handling ID")
				}
				if resBody.TODO.Subject != tc.req.Subject {
					t.Log("Incorrect handling Subject")
				}
				if resBody.TODO.Description != tc.req.Description {
					t.Log("Incorrect handling description")
				}
			}
		})
	}
}
