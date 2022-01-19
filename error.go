package gocrawler

import (
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/smallsung/gopkg/errors"
)

type Error struct {
	errors.LocatorError
	context  *Context
	request  *colly.Request
	response *colly.Response
}

func (err *Error) Context() *Context         { return err.context }
func (err *Error) Request() *colly.Request   { return err.request }
func (err *Error) Response() *colly.Response { return err.response }

func (err *Error) setContext(context *Context) *Error {
	err.context = context
	return err
}
func (err *Error) setRequest(request *colly.Request) *Error {
	err.request = request
	return err
}
func (err *Error) setResponse(response *colly.Response) *Error {
	err.response = response
	return err
}

func makeError(err error) *Error {
	newErr := &Error{LocatorError: errors.NewCauserLocatorError(err)}
	newErr.SetLocation(1)
	return newErr
}

type Errors struct {
	errs []*Error
	err  *Error
	lock sync.RWMutex
}

func (errs *Errors) IsNil() bool {
	errs.lock.RLock()
	defer errs.lock.RUnlock()
	return len(errs.errs) == 0
}

func (errs *Errors) Append(err *Error) {
	if err == nil {
		return
	}
	errs.lock.Lock()
	defer errs.lock.Unlock()
	errs.errs = append(errs.errs, err)
}

func (errs *Errors) Err() *Error {
	errs.lock.RLock()
	defer errs.lock.RUnlock()
	if errs.IsNil() {
		return nil
	}
	return errs.errs[0]
}

func (errs *Errors) Errs() []*Error {
	errs.lock.RLock()
	defer errs.lock.RUnlock()
	return errs.errs
}
