package model

type ErrNotFound struct {
	What string
}

func (e *ErrNotFound) Error() string {
	return e.What
}
