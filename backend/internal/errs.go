package internal

import "github.com/pkg/errors"

type BusingessLogicError struct {
	msg string
}

func (e *BusingessLogicError) Error() string {
	return e.msg
}

func NewBusinessError(msg string) error {
	e := &BusingessLogicError{
		msg: msg,
	}

	return errors.WithStack(e)
}
