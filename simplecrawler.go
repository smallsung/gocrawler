package gocrawler

import (
	"fmt"
	"reflect"

	"github.com/gocolly/colly/v2"
	"github.com/smallsung/gopkg/errors"
	"go.uber.org/zap"
)

type simpleCrawler struct {
	invoker reflect.Value

	colly  *colly.Collector
	logger *zap.Logger

	requestHandlers         []RequestHandler
	responseHandlers        []ResponseHandler
	responseHeadersHandlers []ResponseHeadersHandler
	scrapedHandlers         []ScrapedHandler
	errorHandlers           []ErrorHandler
}

func NewSimpleCrawler(opts ...Option) Crawler {
	s := &simpleCrawler{
		colly:                   colly.NewCollector(),
		logger:                  zap.NewNop(),
		requestHandlers:         make([]RequestHandler, 0, 4),
		responseHandlers:        make([]ResponseHandler, 0, 4),
		responseHeadersHandlers: make([]ResponseHeadersHandler, 0, 4),
		scrapedHandlers:         make([]ScrapedHandler, 0, 4),
		errorHandlers:           make([]ErrorHandler, 0, 4),
	}
	s.colly.OnRequest(s.collyRequest)
	s.colly.OnResponseHeaders(s.collyResponseHeaders)
	s.colly.OnResponse(s.collyResponse)
	s.colly.OnScraped(s.collyScraped)
	s.colly.OnError(s.collyError)

	options := append([]Option{
		// 注意排序
	}, opts...)

	for _, opt := range options {
		if _, ok := opt.(optionFunc); ok {
			opt.Apply(s)
		}
	}

	return s
}

var (
	errInvoker    = fmt.Errorf("signature must be func(gocrawler.Collector, ...Type) (Type, error)")
	collectorType = reflect.TypeOf((*Collector)(nil)).Elem()
	errorType     = reflect.TypeOf((*error)(nil)).Elem()
)

func (s *simpleCrawler) setInvoker(i interface{}) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Func {
		panic(errInvoker)
	}

	typ := val.Type()
	switch typ.NumIn() {
	case 0:
		panic(errInvoker)
	case 1:
		if typ.In(0).Kind() != reflect.Interface || !typ.In(0).Implements(collectorType) {
			panic(errInvoker)
		}
	}

	switch typ.NumOut() {
	case 2:
		if typ.Out(1) != errorType {
			panic(errInvoker)
		}
	default:
		panic(errInvoker)
	}

	s.invoker = val
}

func (s *simpleCrawler) Crawl(out interface{}, ins ...interface{}) error {
	if out == nil || reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.Format("parameter out must be non-nil pointer: %v", out)
	}
	if !s.invoker.IsValid() {
		panic(errInvoker)
	}

	var calls []reflect.Value
	calls = append(calls, reflect.ValueOf(s))
	for _, i := range ins {
		calls = append(calls, reflect.ValueOf(i))
	}
	outs := s.invoker.Call(calls)
	if !outs[1].IsNil() {
		return errors.Trace(outs[1].Interface().(error))
	}
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(outs[0].Interface()))
	return nil
}

func (s *simpleCrawler) Client() Client {
	return s
}

func (s *simpleCrawler) Logger() Logger {
	return s
}

func (s *simpleCrawler) Debug(msg string, fields ...zap.Field) {
	s.logger.Debug(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}
func (s *simpleCrawler) Info(msg string, fields ...zap.Field) {
	s.logger.Info(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}
func (s *simpleCrawler) Warn(msg string, fields ...zap.Field) {
	s.logger.Warn(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}
func (s *simpleCrawler) Error(msg string, fields ...zap.Field) {
	s.logger.Error(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}
func (s *simpleCrawler) DPanic(msg string, fields ...zap.Field) {
	s.logger.DPanic(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}
func (s *simpleCrawler) Panic(msg string, fields ...zap.Field) {
	s.logger.Panic(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}
func (s *simpleCrawler) Fatal(msg string, fields ...zap.Field) {
	s.logger.Fatal(msg, append(fields, zap.Namespace("colly"), zap.Uint32("i", s.colly.ID))...)
}

func (s *simpleCrawler) collyRequest(collyRequest *colly.Request) {
	request := makeRequest(collyRequest, s.logger)
	for _, handler := range s.requestHandlers {
		if err := handler.HandleRequest(request); err != nil {
			err := makeError(errors.Annotate(err, "handle request")).
				setContext(&Context{context: collyRequest.Ctx}).setRequest(collyRequest)
			request.Context().Errors().Append(err)
		}
	}

}

func (s *simpleCrawler) collyResponseHeaders(collyResponse *colly.Response) {
	response := makeResponse(collyResponse, s.logger)
	for _, handler := range s.responseHeadersHandlers {
		if err := handler.HandleResponseHeaders(response); err != nil {
			err := makeError(errors.Annotate(err, "handle responseHeaders")).
				setContext(&Context{context: collyResponse.Ctx}).setResponse(collyResponse).setRequest(collyResponse.Request)
			response.Context().Errors().Append(err)
		}
	}

}
func (s *simpleCrawler) collyResponse(collyResponse *colly.Response) {
	response := makeResponse(collyResponse, s.logger)
	for _, handler := range s.responseHandlers {
		if err := handler.HandleResponse(response); err != nil {
			err := makeError(errors.Annotate(err, "handle response")).
				setContext(&Context{context: collyResponse.Ctx}).setResponse(collyResponse).setRequest(collyResponse.Request)
			response.Context().Errors().Append(err)
		}
	}

}
func (s *simpleCrawler) collyScraped(collyResponse *colly.Response) {
	response := makeResponse(collyResponse, s.logger)
	for _, handler := range s.scrapedHandlers {
		if err := handler.HandleScraped(response); err != nil {
			err := makeError(errors.Annotate(err, "handle scraped")).
				setContext(&Context{context: collyResponse.Ctx}).setResponse(collyResponse).setRequest(collyResponse.Request)
			response.Context().Errors().Append(err)
		}
	}

}

func (s *simpleCrawler) collyError(collyResponse *colly.Response, err error) {
	if err == nil {
		return
	}
	response := makeResponse(collyResponse, s.logger)

	response.Context().Errors().err = makeError(errors.Annotate(err, "colly.handleOnError")).
		setRequest(collyResponse.Request).
		setResponse(collyResponse).
		setContext(&Context{context: collyResponse.Ctx})

	for _, handler := range s.errorHandlers {
		if err := handler.HandleError(response, err); err != nil {
			err := makeError(errors.Annotate(err, "handle error")).
				setContext(&Context{context: collyResponse.Ctx}).setResponse(collyResponse).setRequest(collyResponse.Request)
			response.Context().Errors().Append(err)
		}
	}
}
