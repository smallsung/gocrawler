package parse

import (
	"regexp"

	"github.com/smallsung/gopkg/errors"
)

type RegexpNotMatchError struct {
	errors.Err
}

type _regexp struct {
	bytes  []byte
	string *string
	ptr    regexpPtr
}

func (receiver *_regexp) FindSubMatch(regexp *regexp.Regexp) ([][]byte, error) {
	return receiver.ptr.findSubMatch(regexp, receiver.bytes)
}

func (receiver *_regexp) FindStringSubMatch(regexp *regexp.Regexp) ([]string, error) {
	return receiver.ptr.findStringSubMatch(regexp, *receiver.string)
}

func (receiver *_regexp) FindAllStringSubMatch(regexp *regexp.Regexp, n int) ([][]string, error) {
	return receiver.ptr.findAllStringSubMatch(regexp, *receiver.string, n)
}

type RegexpWrapper struct {
	regexp *regexp.Regexp
	ptr    regexpPtr
}

func Regexp(regexp *regexp.Regexp) *RegexpWrapper {
	return &RegexpWrapper{regexp: regexp}
}

func (receiver *RegexpWrapper) FindSubMatch(b []byte) ([][]byte, error) {
	return receiver.ptr.findSubMatch(receiver.regexp, b)
}

func (receiver *RegexpWrapper) FindStringSubMatch(s string) ([]string, error) {
	return receiver.ptr.findStringSubMatch(receiver.regexp, s)
}

func (receiver *RegexpWrapper) FindAllStringSubMatch(s string, n int) ([][]string, error) {
	return receiver.ptr.findAllStringSubMatch(receiver.regexp, s, n)
}

type regexpPtr uintptr

func (receiver *regexpPtr) notMatch(regexp *regexp.Regexp) *RegexpNotMatchError {
	err := &RegexpNotMatchError{Err: errors.NewErr("正则表达式没有匹配:%s", regexp.String())}
	err.SetLocation(3)
	return err
}

func (receiver regexpPtr) findSubMatch(regexp *regexp.Regexp, b []byte) ([][]byte, error) {
	if matches := regexp.FindSubmatch(b); len(matches) == 0 {
		return nil, receiver.notMatch(regexp)
	} else {
		return matches, nil
	}
}

func (receiver regexpPtr) findStringSubMatch(regexp *regexp.Regexp, s string) ([]string, error) {
	if matches := regexp.FindStringSubmatch(s); len(matches) == 0 {
		return nil, receiver.notMatch(regexp)
	} else {
		return matches, nil
	}
}

func (receiver regexpPtr) findAllStringSubMatch(regexp *regexp.Regexp, s string, n int) ([][]string, error) {
	if matches := regexp.FindAllStringSubmatch(s, n); len(matches) == 0 {
		return nil, receiver.notMatch(regexp)
	} else {
		return matches, nil
	}
}
