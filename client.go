package gocrawler

import (
	"io"
	"net/http"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/smallsung/gopkg/errors"
)

type request struct {
	context *Context
	logger  Logger
	request *colly.Request
}

func makeRequest(r *colly.Request, l Logger) *request {
	return &request{context: &Context{context: r.Ctx}, logger: l, request: r}
}

func (r *request) Context() *Context       { return r.context }
func (r *request) Logger() Logger          { return r.logger }
func (r *request) Request() *colly.Request { return r.request }

type response struct {
	context  *Context
	logger   Logger
	response *colly.Response
}

func makeResponse(r *colly.Response, l Logger) *response {
	return &response{context: &Context{context: r.Ctx}, logger: l, response: r}
}

func (r *response) Context() *Context         { return r.context }
func (r *response) Logger() Logger            { return r.logger }
func (r *response) Response() *colly.Response { return r.response }

func (s *simpleCrawler) request(context *Context, method HttpMethod, url string, requestData io.Reader, header http.Header) error {
	if context == nil || context.context == nil {
		panic("*Context must be non-nil. you can use NewContext().")
	}

	err := s.colly.Request(string(method), url, requestData, context.context, header)
	if err != nil {
		return errors.Trace(err)
	}

	if !context.Errors().IsNil() {
		return errors.Trace(context.Errors().Err())
	}
	return nil
}

func (s *simpleCrawler) Request(context *Context, method HttpMethod, url string, requestData io.Reader, header http.Header) error {
	return s.request(context, method, url, requestData, header)
}

func (s *simpleCrawler) MultiRequest(context *Context, requests MultiRequest) error {
	if len(requests.slice) == 0 {
		return nil
	}

	mutex := sync.Mutex{}
	var errs []error
	wg := sync.WaitGroup{}
	for _, r := range requests.slice {
		wg.Add(1)
		r := r
		go func() {
			if err := s.request(context, r.Method, r.URL, r.Body, r.Headers); err != nil {
				mutex.Lock()
				errs = append(errs, err)
				mutex.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(errs) != 0 {
		return errors.Trace(errs[0])
	}

	if !context.Errors().IsNil() {
		return errors.Trace(context.Errors().Err())
	}
	return nil
}
func (s *simpleCrawler) Get(context *Context, url string) error {
	return s.request(context, HttpMethodGet, url, nil, nil)
}

func (s *simpleCrawler) Post(context *Context, url string, requestData map[string]string) error {
	return s.request(context, HttpMethodPost, url, Map2Reader(requestData), nil)
}
