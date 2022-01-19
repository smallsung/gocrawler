package parse

import (
	"encoding/json"

	"github.com/smallsung/gopkg/errors"
)

type (
	JsonUnmarshalError struct {
		errors.Err
	}
)

type _json struct {
	bytes []byte
}

func (receiver *_json) Unmarshal(v interface{}) error {
	if err := json.Unmarshal(receiver.bytes, v); err != nil {
		return receiver.wrapError(err)
	}
	return nil
}

func (receiver *_json) wrapError(err error) error {
	newErr := &JsonUnmarshalError{Err: errors.NewErrWithCause(err, "")}
	newErr.SetLocation(2)
	return newErr
}
