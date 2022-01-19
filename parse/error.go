package parse

import (
	"github.com/smallsung/gopkg/errors"
)

type NotExpectValueError struct {
	errors.Err
}

func NewNotExpectValueError(got interface{}) error {
	err := &NotExpectValueError{Err: errors.NewErr("非预期结果. Got %+v", got)}
	err.SetLocation(1)
	return err
}
