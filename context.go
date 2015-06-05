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

func contextClient(ctx context.Context) (*http.Client, error) {
	if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
		return hc, nil
	}
	for _, fn := range contextClientFuncs {
		c, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		if c != nil {
			return c, nil
		}
	}
	return http.DefaultClient, nil
}

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient contextKey1

// APIEndpoint is the context key to determine the endpoint (demo or live)
// to use. If not set, the live version is assumed.
var APIEndpoint contextKey2

// contextKey is just an empty struct. It exists so HTTPClient can be
// an immutable public variable with a unique type. It's immutable
// because nobody else can create a contextKey, being unexported.
type contextKey1 struct{}
type contextKey2 struct{}

// DefaultCtx is the default context you should supply if not using
// your own context.Context (see https://golang.org/x/net/context).
var DefaultCtx = context.Background()
