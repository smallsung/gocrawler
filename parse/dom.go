package parse

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/smallsung/gopkg/errors"
	"golang.org/x/net/html"
)

type DocumentError struct {
	errors.LocatorError
}

type DomEachFunc func(i int, dom *Document) error

type Document struct {
	Selection *goquery.Selection
}

func newDocumentFromBytes(b []byte) (*Document, error) {
	if node, err := html.Parse(bytes.NewReader(b)); err != nil {
		newErr := &DocumentError{LocatorError: errors.NewCauserAnnotatorLocatorError(err, "new dom failed")}
		newErr.SetLocation(2)
		return nil, newErr
	} else {
		return &Document{
			Selection: goquery.NewDocumentFromNode(node).Selection,
		}, nil
	}
}

func newDocumentFromNode(selection *goquery.Selection) *Document {
	return &Document{Selection: selection}
}

func NewDocumentFromBytes(b []byte) (dom *Document, err error) {
	return newDocumentFromBytes(b)
}

func Dom(b []byte) (dom *Document, err error) {
	return newDocumentFromBytes(b)
}

func (dom *Document) error(format string, args ...interface{}) error {
	newErr := &DocumentError{LocatorError: errors.NewLocatorError(format, args...)}
	newErr.SetLocation(2)
	return newErr
}
func (dom *Document) errorWithCause(err error) error {
	if err == nil {
		return nil
	}
	newErr := &DocumentError{LocatorError: errors.NewCauserLocatorError(err)}
	newErr.SetLocation(2)
	return newErr
}

func (dom *Document) errorWithCauseAnnotate(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	newErr := &DocumentError{LocatorError: errors.NewCauserAnnotatorLocatorError(err, format, args...)}
	newErr.SetLocation(2)

	return newErr
}

func (dom *Document) MustFind(selector string) (*Document, error) {
	selection := dom.Selection.Find(selector)
	if selection.Length() == 0 {
		return nil, dom.error("dom 没有找到:%s", selector)
	}
	return newDocumentFromNode(selection), nil
}

func (dom *Document) Find(selector string) *Document {
	selection := dom.Selection.Find(selector)
	return newDocumentFromNode(selection)
}

func (dom *Document) Eq(index int) *Document {
	return newDocumentFromNode(dom.Selection.Eq(index))
}

func (dom *Document) Length() int {
	return dom.Selection.Length()
}

func (dom *Document) Text() string {
	return dom.Selection.Text()
}

func (dom *Document) TrimSpaceText() string {
	return strings.TrimSpace(dom.Selection.Text())
}

func (dom *Document) InnerText() string {
	var buf bytes.Buffer
	for _, n := range dom.Selection.Nodes {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					buf.WriteString(c.Data)
				}
			}
		}
	}
	return buf.String()
}

func (dom *Document) TrimSpaceInnerText() string {
	return strings.TrimSpace(dom.InnerText())
}

func (dom *Document) InnerHTML() (innerHTML string, err error) {
	if innerHTML, err = dom.Selection.Html(); err != nil {
		return "", dom.errorWithCause(err)
	} else {
		return innerHTML, nil
	}
}

func (dom *Document) InnerHTMLWithoutError() string {
	if Html, err := dom.Selection.Html(); err != nil {
		return ""
	} else {
		return Html
	}
}

func (dom *Document) OuterHtml() (outerHtml string, err error) {
	if outerHtml, err = goquery.OuterHtml(dom.Selection); err != nil {
		return "", dom.errorWithCause(err)
	} else {
		return outerHtml, nil
	}
}

//Each 返回最后一次错误
func (dom *Document) Each(fn DomEachFunc) error {
	var bubble error
	dom.Selection.Each(func(i int, selection *goquery.Selection) {
		if err := fn(i, newDocumentFromNode(selection)); err != nil {
			bubble = errors.Trace(err)
		}
	})
	return bubble
}

func (dom *Document) EachWithBreak(fn DomEachFunc) error {
	var bubble error
	dom.Selection.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if err := fn(i, newDocumentFromNode(selection)); err != nil {
			bubble = errors.Trace(err)
			return false
		}
		return true
	})
	return bubble
}

func (dom *Document) Attr(attrName string) (string, error) {
	if attr, exist := dom.Selection.Attr(attrName); exist {
		return attr, nil
	} else {
		return "", dom.error("attr does not exist %s", attrName)
	}
}
