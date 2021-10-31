package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/service"
)

var init_data = []struct {
	subject     string
	description string
}{
	{
		subject: "foo",
		description: "this is foo",
	},
	{
		subject: "bar",
		description: "this is bar",
	},
	{
		subject: "baz",
		description: "this is baz",
	},
}

func TestCreateTODO(t *testing.T) {
	dbpath := "./todo_temp.db"
	todoDB, err := db.NewDB(dbpath)
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
		isError bool
	}{
		{
			name: "normal",
			subject: "hello world",
			descritpion: "this is test",
			isError: false,
		},
		{
			name: "empty subject",
			subject: "",
			descritpion: "this is test",
			isError: true,
		},
		{
			name: "empty description",
			subject: "hello world",
			descritpion: "",
			isError: false,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func (t *testing.T) {
			todo, err := svc.CreateTODO(ctx, tc.subject, tc.descritpion)
			switch {
			case tc.isError && err == nil:
				t.Fatal("expected err, but err is nil")
			case !tc.isError && err != nil:
				t.Fatal("not expected err, but err is not nil: ", err)
			}

			if !tc.isError {
				if tc.subject != todo.Subject {
					t.Fatal("expected: ", tc.subject, ", actual: ", todo.Subject)
				}
				if tc.descritpion != todo.Description {
					t.Fatal("expected: ", tc.descritpion, ", actual: ", todo.Description)
				}
			}
		})
	}

	if err := os.Remove(dbpath); err != nil {
		t.Log(err)
	}
}

func TestReadTODO(t *testing.T) {
	dbpath := "./todo_temp.db"
	todoDB, err := db.NewDB(dbpath)
	if err != nil {
		t.Fatal(err)
	}
	defer todoDB.Close()

	ctx := context.Background()

	stmt, err := todoDB.PrepareContext(ctx, "INSERT INTO todos(subject, description) VALUES(?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range init_data {
		if _, err := stmt.ExecContext(ctx, data.subject, data.description); err != nil {
			t.Fatal(err)
		}
	}
	svc := service.NewTODOService(todoDB)

	testcase := []struct {
		name    string
		prevID  int64
		size    int64
		isError bool
	}{
		{
			name: "normal",
			prevID: 2,
			size: 2,
			isError: false,
		},
		{
			name: "error",
			prevID: -1,
			size: -1,
			isError: true,
		},
		{
			name: "size > data size",
			prevID: 0,
			size: 10,
			isError: false,
		},
		{
			name: "prevID",
			prevID: 1,
			size: 10,
			isError: false,
		},
		{
			name: "size > data size",
			prevID: 0,
			size: 10,
			isError: false,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func (t *testing.T) {
			todos, err := svc.ReadTODO(ctx, tc.prevID, tc.size)
			switch {
			case tc.isError && err == nil:
				t.Fatal("expected err, but err is nil")
			case !tc.isError && err != nil:
				t.Fatal("not expected err, but err is not nil: ", err)
			}

			if !tc.isError {
				for i, todo := range todos {
					if tc.prevID > 0 && tc.prevID < todo.ID {
						t.Fatal("range error")
					}
					if init_data[i].subject != todo.Subject {
						t.Fatal("expected: ", init_data[i].subject, ", actual: ", todo.Subject)
					}
					if init_data[i].description != todo.Description {
						t.Fatal("expected: ", init_data[i].description, ", actual: ", todo.Description)
					}
				}
			}
		})
	}

	if err := os.Remove(dbpath); err != nil {
		t.Log(err)
	}
}
