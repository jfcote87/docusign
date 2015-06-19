// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docusign

// Logic for obtaining client via the a golang.org/x/net/context Context.
// Code copied from golang.org/x/oauth2
import (
	"net/http"

	"golang.org/x/net/context"
)

// contextClientFunc is a func which tries to return an *http.Client
// given a Context value. If it returns an error, the search stops
// with that error.  If it returns (nil, nil), the search continues
// down the list of registered funcs.
type contextClientFunc func(context.Context) (*http.Client, error)

var contextClientFuncs []contextClientFunc

func registerContextClientFunc(fn contextClientFunc) {
	contextClientFuncs = append(contextClientFuncs, fn)
}

// ContextSetting provides default values for creating a new
// docusing.Serivce.  If needed, a Context's docusign.APISettings
// Value shoud be set.
type ContextSetting struct {
	Client         *http.Client
	Endpoint       string
	LogRawRequest  func(context.Context, string, ...interface{})
	LogRawResponse func(context.Context, string, ...interface{})
}

// contextClient returns the appropriate
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

// contextSettings returns a ContextSetting from the context.
func contextSettings(ctx context.Context) *ContextSetting {
	if cs, ok := ctx.Value(APISettings).(*ContextSetting); ok {
		if cs.Client == nil {
			return &ContextSetting{Client: contextClient(ctx), Endpoint: cs.Endpoint, LogRawRequest: cs.LogRawRequest, LogRawResponse: cs.LogRawResponse}
		}
		return cs
	}
	return &ContextSetting{Client: contextClient(ctx), Endpoint: liveUrl, LogRawRequest: nil, LogRawResponse: nil}

}

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient contextKey1

// APIEndpoint is the context key to determine the endpoint (demo or live)
// to use. If not set, the live version is assumed.
var APISettings contextKey2

// contextKeyX is just an empty struct. It exists so HTTPClient can be
// an immutable public variable with a unique type. It's immutable
// because nobody else can create a contextKey, being unexported.
type contextKey1 struct{}
type contextKey2 struct{}

// DefaultCtx is the default context you should supply if not using
// your own context.Context (see https://golang.org/x/net/context).
var DefaultCtx = context.Background()
