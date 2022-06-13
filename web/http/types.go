package http

import (
	"context"
	"net/http"
)

// Handler is a function signature which retains most of the important bits of http.HandlerFunc, but enhances it with
// an incoming context, so we don't need to mangle the one in http.Request, and returns an error, which will be useful
// once we have an error handling middleware.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Middleware is a function signature that takes a Handler, and returns a Handler. The idea here is that the passed in
// Handler is going to be wrapped into whatever functionality we need, and then the new, enhanced Handler gets returned.
type Middleware func(h Handler) Handler

// WrapMiddleware creates a new handler by wrapping middleware around a final handler. The middlewares' Handlers will be
// executed by requests in the order they are provided.
func WrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Loop backwards through the middleware invoking each one. Replace the handler with the new wrapped handler.
	// Looping backwards ensures that the first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
