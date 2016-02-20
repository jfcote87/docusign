// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docusign

// Logic for obtaining client via the a golang.org/x/net/context Context.
// Code copied from golang.org/x/oauth2
import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/context"
)

// DefaultCtx is the default context you should supply if not using
// your own context.Context (see https://golang.org/x/net/context).
var DefaultCtx = context.Background()

type ctxKeyHTTPClient struct{}
type ctxKeyLogger struct{}

// Logger provides a mechanism to log call made via a Service.
// If a context has the docusign.CallLogger value set to a
// Logger, any service call will log requests
// using the interfaces functions.
type Logger interface {
	LogRequest(cxt context.Context, payload interface{}, req *http.Request)
	// Log Response needs to return a duplicate io.Reader of the res.Body.
	LogResponse(ctx context.Context, res *http.Response) io.Reader
}

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient ctxKeyHTTPClient
var CallLogger ctxKeyLogger

// contextClientFunc is a func which tries to return an *http.Client
// given a Context value. If it returns an error, the search stops
// with that error.  If it returns (nil, nil), the search continues
// down the list of registered funcs.
type contextClientFunc func(context.Context) (*http.Client, error)

var contextClientFuncs []contextClientFunc

func registerContextClientFunc(fn contextClientFunc) {
	contextClientFuncs = append(contextClientFuncs, fn)
}

// contextClient returns the appropriate client for the
// provided context.
func contextClient(ctx context.Context) *http.Client {
	if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
		return hc
	}
	for _, fn := range contextClientFuncs {
		c, err := fn(ctx)
		if err != nil {
			panic(err)
		}
		if c != nil {
			return c
		}
	}
	return http.DefaultClient
}

// contextLogger returns a Logger associated with the
// provided context.
func contextLogger(ctx context.Context) Logger {
	if f, ok := ctx.Value(CallLogger).(Logger); ok {
		return f
	}
	return nil
}

type SimpleLogger struct{}

func (s SimpleLogger) LogRequest(ctx context.Context, payload interface{}, req *http.Request) {
	log.Printf("URL is %s, Payload: %#v", req.URL, payload)
}

func (s SimpleLogger) LogResponse(ctx context.Context, res *http.Response) io.Reader {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Unable to read response: %v", res.Request.URL)
		return &bytes.Reader{}
	}
	log.Printf("Received %d bytes: %s", res.ContentLength, string(b))
	return bytes.NewReader(b)
}
