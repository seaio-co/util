package httplog

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// Transport satisfies http.RoundTripper
type Transport struct {
	transport http.RoundTripper
	// Should the body of the requests and responses be logged.
	logBody bool
	// If logf is nil log.Printf will be used.
	logf func(format string, vs ...interface{})
}

// NewTransport returns a new Transport that uses the given RoundTripper, or
// http.DefaultTransport if nil, and logs all requests and responses using
// logf, or log.Printf if nil.
// The body of the requests and responses are logged too only if logBody is
// true.
func NewTransport(rt http.RoundTripper, logBody bool, logf func(string, ...interface{})) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	if logf == nil {
		logf = log.Printf
	}
	return Transport{rt, logBody, logf}
}

// Client returns a new http.Client using the given transport.
func (t Transport) Client() *http.Client { return &http.Client{Transport: t} }

// RoundTrip so Transport satifies http.RoundTripper
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequest(req, t.logBody)
	if err != nil {
		t.logf("httplog: dump request: %v", err)
		return nil, err
	}
	t.logf("httplog: %s", b)

	res, err := t.transport.RoundTrip(req)
	if err != nil {
		t.logf("httplog: roundtrip error: %v", err)
		return res, err
	}

	if b, err := httputil.DumpResponse(res, t.logBody); err != nil {
		t.logf("httplog: dump response: %v", err)
	} else {
		t.logf("httplog: %s", b)
	}
	return res, err
}
