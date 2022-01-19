package parse

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync/atomic"

	"github.com/smallsung/gopkg/errors"
)

type Parser struct {
	bytes  []byte
	string *string

	Json   *_json
	Regexp *_regexp

	Dom *Document
}

type ParserMode uint8

const (
	ModeEnableJson ParserMode = 1 << iota
	ModeEnableDom
	ModeEnableRegexp
)

func Parse(b []byte, mode ParserMode) (*Parser, error) {
	p := new(Parser)
	if mode&ModeEnableRegexp != 0 {
		s := string(b)
		p.Regexp = &_regexp{bytes: b, string: &s}
		p.string = &s
	}
	if mode&ModeEnableDom != 0 {
		if dom, err := newDocumentFromBytes(b); err != nil {
			return nil, errors.Annotate(err, "parse.newDocumentFromBytes")
		} else {
			p.Dom = dom
		}
	}
	if mode&ModeEnableJson != 0 {
		p.Json = &_json{bytes: b}
	}
	p.bytes = b
	return p, nil
}

var debugStoreTestDir uint64

func (receiver *Parser) DebugStoreTestDir() {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("")
	}
	filename := filepath.Join(filepath.Dir(file), fmt.Sprintf("___debug_%d", atomic.AddUint64(&debugStoreTestDir, 1)))
	if err := ioutil.WriteFile(filename, receiver.bytes, 0777); err != nil {
		panic("")
	}
}
