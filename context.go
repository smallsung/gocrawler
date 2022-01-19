package gocrawler

import (
	"reflect"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/smallsung/gopkg/errors"
)

const (
	contextKeyErrors = "___contextKeyErrors"
	contextKeyLocker = "___contextKeyLocker"
)

type TypeAssertionError struct {
	errors.LocatorError
}

func newTypeAssertionError(key string, want reflect.Type, got interface{}) *TypeAssertionError {
	err := &TypeAssertionError{
		LocatorError: errors.NewLocatorError("type assertion[%s]: want %s got %T", key, want.String(), got),
	}
	err.SetLocation(1)
	return err
}

type Context struct {
	context *colly.Context
}

func NewContext() *Context {
	ctx := colly.NewContext()
	ctx.Put(contextKeyErrors, new(Errors))
	ctx.Put(contextKeyLocker, &sync.Mutex{})
	return &Context{context: ctx}
}

func (ctx *Context) Put(key string, value interface{}) {
	ctx.context.Put(key, value)
}

func (ctx *Context) MustGot(key string, receiver interface{}) {

	receiverVal := reflect.ValueOf(receiver)
	if receiver == nil || receiverVal.Kind() != reflect.Ptr {
		panic(errors.Errorf("parameter receiver must be pointer or non-nil: %v", receiver))
	}

	any := ctx.context.GetAny(key)
	if any == nil {
		switch receiverVal.Elem().Kind() {
		//case reflect.Slice:
		//	receiverVal.Elem().Set(reflect.New(receiverVal.Elem().Type()).Elem())
		//	return
		//case reflect.Ptr:
		//	receiverVal.Elem().Set(reflect.New(receiverVal.Elem().Type()).Elem())
		//	return
		//case reflect.Bool,
		//	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		//	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		//	reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		//	reflect.Uintptr, reflect.String,
		//	reflect.Struct:
		//	panic(newTypeAssertionError(key, receiverVal.Elem().Type(), any))
		default:
			panic(newTypeAssertionError(key, receiverVal.Elem().Type(), any))
		}
	} else {
		anyVal := reflect.ValueOf(any)
		switch anyVal.Kind() {
		case reflect.Ptr:
			receiverVal.Elem().Set(anyVal)

		case reflect.Bool,
			reflect.Int:
			receiverVal.Elem().Set(anyVal)
		case reflect.Slice:
			receiverVal.Elem().Set(anyVal)
		case reflect.Struct:
			receiverVal.Elem().Set(anyVal)
		default:
			panic(newTypeAssertionError(key, receiverVal.Elem().Type(), any))
		}
	}
}

func (ctx *Context) Errors() *Errors {
	return ctx.context.GetAny(contextKeyErrors).(*Errors)
}

func (ctx *Context) Uint(key string) uint {
	var receiver uint
	ctx.MustGot(key, &receiver)
	return receiver
}

func (ctx *Context) Int(key string) int {
	var receiver int
	ctx.MustGot(key, &receiver)
	return receiver
}

func (ctx *Context) Bool(key string) bool {
	var receiver bool
	ctx.MustGot(key, &receiver)
	return receiver
}
