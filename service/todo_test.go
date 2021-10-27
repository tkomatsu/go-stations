package service_test

import (
	"context"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/service"
)

func TestCreateTODO(t *testing.T) {
	todoDB, err := db.NewDB("../.sqlite3/todo_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()
	svc := service.NewTODOService(todoDB)
	ctx := context.Background()

	testcase := []struct {
		name string
		subject string
		descritpion string
	}{
		{
			name: "normal",
			subject: "hello world",
			descritpion: "this is test",
		},
		{
			name: "empty subject",
			subject: "",
			descritpion: "this is test",
		},
	}
}
