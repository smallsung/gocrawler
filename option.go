package gocrawler

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"go.uber.org/zap"
)

type Option interface {
	//Apply 不应该在构造方法外调用
	Apply(c Crawler)
}

type optionFunc func(s *simpleCrawler)

func (f optionFunc) Apply(c Crawler) {
	if s := c.(*simpleCrawler); s == nil {
		panic("")
	} else {
		f(s)
	}
}

func SetInvokerOption(i interface{}) Option {
	return optionFunc(func(s *simpleCrawler) {
		s.setInvoker(i)
	})
}

func WithLoggerOption(l *zap.Logger) Option {
	return optionFunc(func(s *simpleCrawler) {
		if l != nil {
			s.logger = l
		}
	})
}

type debugger struct {
	logger *zap.Logger
}

func (d debugger) Init() error { return nil }
func (d debugger) Event(e *debug.Event) {
	d.logger.Debug("colly.debug", zap.Uint32("c", e.CollectorID), zap.Uint32("r", e.RequestID), zap.String("t", e.Type), zap.Any("v", e.Values))
}

func DefaultDebuggerOption(l *zap.Logger) Option {
	return DebuggerOption(debugger{logger: l})
}

// Handlers

func RequestHandlerOption(handler RequestHandler) Option {
	return optionFunc(func(s *simpleCrawler) {
		s.requestHandlers = append(s.requestHandlers, handler)
	})
}

func ResponseHeadersHandlerOption(handler ResponseHeadersHandler) Option {
	return optionFunc(func(s *simpleCrawler) {
		s.responseHeadersHandlers = append(s.responseHeadersHandlers, handler)
	})
}

func ResponseHandlerOption(handler ResponseHandler) Option {
	return optionFunc(func(s *simpleCrawler) {
		s.responseHandlers = append(s.responseHandlers, handler)
	})
}

func ScrapedHandlerOption(handler ScrapedHandler) Option {
	return optionFunc(func(s *simpleCrawler) {
		s.scrapedHandlers = append(s.scrapedHandlers, handler)
	})
}

func ErrorHandlerOption(handler ErrorHandler) Option {
	return optionFunc(func(s *simpleCrawler) {
		s.errorHandlers = append(s.errorHandlers, handler)
	})
}

//  colly::Option

func DebuggerOption(d debug.Debugger) Option {
	return optionFunc(func(s *simpleCrawler) {
		colly.Debugger(d)(s.colly)
	})
}

func LimitOption(parallelism int, delay, randomDelay time.Duration) Option {
	rule := &colly.LimitRule{
		DomainRegexp: ".*",
		DomainGlob:   "",
		Parallelism:  parallelism,
		Delay:        delay,
		RandomDelay:  randomDelay,
	}
	return optionFunc(func(s *simpleCrawler) {
		if err := s.colly.Limit(rule); err != nil {
			panic(err)
		}
	})
}
