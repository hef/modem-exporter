// +build gofuzzbeta

package client

import (
	"bytes"
	"github.com/antchfx/htmlquery"
	"go.uber.org/zap"
	"io/ioutil"
	"testing"
)

func FuzzParseStatus(f *testing.F) {

	logger := zap.NewNop()
	x, err := ioutil.ReadFile("testdata/cmconnectionstatus.html")
	if err != nil {
		panic(err)
	}
	f.Add(x)

	f.Fuzz(func(t *testing.T, in []byte) {
		r := bytes.NewReader(in)
		doc, err := htmlquery.Parse(r)
		if err != nil {
			t.Fatal("What?")
		}
		parseStatusPage(logger, doc)
	})

}
