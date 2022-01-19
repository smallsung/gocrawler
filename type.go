package gocrawler

import (
	"io"
	"net/http"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

type Crawler interface {
	Crawl(out interface{}, ins ...interface{}) error
}

type Collector interface {
	Client() Client
	Logger() Logger
}

type Client interface {
	Request(ctx *Context, method HttpMethod, url string, requestData io.Reader, header http.Header) error
	MultiRequest(*Context, MultiRequest) error

	Get(ctx *Context, url string) error
	Post(ctx *Context, url string, requestData map[string]string) error
}

type MultiRequest struct {
	slice []struct {
		URL     string
		Method  HttpMethod
		Body    io.Reader
		Headers http.Header
	}
}

func NewMultiRequest() MultiRequest {
	return MultiRequest{}
}

func (m *MultiRequest) Append(method HttpMethod, url string, requestData io.Reader, header http.Header) {
	m.slice = append(m.slice, struct {
		URL     string
		Method  HttpMethod
		Body    io.Reader
		Headers http.Header
	}{URL: url, Method: method, Body: requestData, Headers: header})
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type (
	Request interface {
		Request() *colly.Request
		Context() *Context
		Logger() Logger
	}

	Response interface {
		Response() *colly.Response
		Context() *Context
		Logger() Logger
	}
)

type (
	RequestHandler interface {
		HandleRequest(request Request) error
	}
	ResponseHandler interface {
		HandleResponse(response Response) error
	}
	ResponseHeadersHandler interface {
		HandleResponseHeaders(response Response) error
	}
	ScrapedHandler interface {
		HandleScraped(response Response) error
	}
	ErrorHandler interface {
		HandleError(response Response, err error) error
	}
)

type (
	RequestHandlerFunc         func(request Request) error
	ResponseHeadersHandlerFunc func(response Response) error
	ResponseHandlerFunc        func(response Response) error
	ScrapedHandlerFunc         func(response Response) error
	ErrorHandlerFunc           func(response Response, err error) error
)

func (f RequestHandlerFunc) HandleRequest(request Request) error {
	return f(request)
}

func (f ResponseHeadersHandlerFunc) HandleResponseHeaders(response Response) error {
	return f(response)
}

func (f ResponseHandlerFunc) HandleResponse(response Response) error {
	return f(response)
}

func (f ScrapedHandlerFunc) HandleScraped(response Response) error {
	return f(response)
}

func (f ErrorHandlerFunc) HandleError(response Response, err error) error {
	return f(response, err)
}
