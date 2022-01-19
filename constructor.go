package gocrawler

import (
	"sync"
)

type Constructor interface {
	Construct(options ...Option) Crawler
}

type ConstructorFunc func(options ...Option) Crawler

func (f ConstructorFunc) Construct(options ...Option) Crawler {
	return f(options...)
}

var (
	constructor      Constructor
	constructorMutex = sync.Mutex{}
)

func init() {
	constructor = ConstructorFunc(NewSimpleCrawler)
}

func ReplaceConstructor(newConstructor Constructor) (oldConstructor Constructor) {
	constructorMutex.Lock()
	defer constructorMutex.Unlock()
	oldConstructor = constructor
	constructor = newConstructor
	return oldConstructor
}

func RecoveryConstructor() Constructor {
	return ReplaceConstructor(ConstructorFunc(NewSimpleCrawler))
}

func New(opts ...Option) Crawler {
	return constructor.Construct(opts...)
}
