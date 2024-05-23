package users

import (
	"errors"
	"fmt"
)

var ErrFirstNameIsRequired = errors.New("first name is required")
var ErrLastNameIsRequired = errors.New("last name is required")
var ErrSliceIsEmpty = errors.New("the slice is empty")

type ErrNotFound struct {
	ID uint64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user with id %d not found", e.ID)
}
