package gocrawler

import (
	"io"
	"net/url"
	"strings"
)

// CopyMap 非深度拷贝，主要用于query Request.body拷贝
func CopyMap(m1 map[string]string) map[string]string {
	m2 := make(map[string]string)
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

// Map2Reader github.com/gocolly/colly::createFormReader
func Map2Reader(m map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range m {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}
