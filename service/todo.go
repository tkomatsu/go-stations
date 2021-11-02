package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmtInsert, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	stmtConfirm, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}

	// validate arguments
	if subject == "" {
		return nil, errors.New("subject not found")
	}

	// insert operation
	ret, err := stmtInsert.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	// confirm operatrion
	id, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}
	var todo model.TODO
	err = stmtConfirm.QueryRowContext(ctx, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	todo.ID = id
	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	stmtRead, err := s.db.PrepareContext(ctx, read)
	if err != nil {
		return nil, err
	}
	stmtReadID, err := s.db.PrepareContext(ctx, readWithID)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if prevID < 0 || size < 0 {
		return nil, errors.New("invalid argument")
	}
	if size == 0 {
		size = -1
	}
	if prevID == 0 {
		rows, err = stmtRead.QueryContext(ctx, size)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = stmtReadID.QueryContext(ctx, prevID, size)
		if err != nil {
			return nil, err
		}
	}

	todos := make([]*model.TODO, 0)
	for rows.Next() {
		var todo model.TODO
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmtUpdate, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	stmtConfirm, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}

	if subject == "" {
		return nil, errors.New("subject not found")
	}

	_, err = stmtUpdate.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	var todo model.TODO
	err = stmtConfirm.QueryRowContext(ctx, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, &model.ErrNotFound{What:err.Error()}
	}
	todo.ID = id
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
