package e2e

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

type TesterParams struct {
	BaseURL      string
	CurlPrinter  bool
	DebugPrinter bool
}

type Tester struct {
	*httpexpect.Expect
}

func NewTester(t *testing.T, params *TesterParams) *Tester {
	var printers []httpexpect.Printer
	if params.CurlPrinter {
		printers = append(printers, httpexpect.NewCurlPrinter(t))
	}
	if params.DebugPrinter {
		printers = append(printers, httpexpect.NewDebugPrinter(t, true))
	}
	return &Tester{
		Expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL: params.BaseURL,
			Client: &http.Client{
				Jar: httpexpect.NewJar(),
			},
			Reporter: httpexpect.NewRequireReporter(t),
			Printers: printers,
		}),
	}
}

//----------------------------------------------
// Extends of httpexpect.
//----------------------------------------------

func (c *Tester) GetWithAuthToken(path, token string, pathargs ...interface{}) *httpexpect.Request {
	return c.GET(path, pathargs...).WithHeader("Authorization", "Token "+token)
}

func (c *Tester) PostWithAuthToken(path, token string, pathargs ...interface{}) *httpexpect.Request {
	return c.POST(path, pathargs...).WithHeader("Authorization", "Token "+token)
}

func (c *Tester) PutWithAuthToken(path, token string, pathargs ...interface{}) *httpexpect.Request {
	return c.PUT(path, pathargs...).WithHeader("Authorization", "Token "+token)
}

func (c *Tester) DeleteWithAuthToken(path, token string, pathargs ...interface{}) *httpexpect.Request {
	return c.DELETE(path, pathargs...).WithHeader("Authorization", "Token "+token)
}
