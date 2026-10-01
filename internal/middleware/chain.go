package middleware

import "net/http"

// Middleware representa una función que intercepta un request.
type Middleware func(http.Handler) http.Handler

// Chain enlaza ordenadamente los middlewares alrededor del Handler final.
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Iteramos en reversa para garantizar orden de declaración
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
