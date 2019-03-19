package ocmux

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// Middleware holds a Gorilla Mux middleware to update the OpenCensus span name and the http route tag.
// It uses the configured route name (if any), otherwise falls back to the request method and the path template.
// This is typically used in a HTTP server using Gorilla Mux for routing.
func Middleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			route := mux.CurrentRoute(req)
			span := trace.FromContext(req.Context())

			if route == nil || span == nil {
				next.ServeHTTP(w, req)

				return
			}

			if name := getRouteName(route, req); name != "" {
				span.SetName(name)
				ochttp.SetRoute(req.Context(), name)
			}

			next.ServeHTTP(w, req)
		})
	}
}

func getRouteName(route *mux.Route, req *http.Request) string {
	name := route.GetName()
	if name == "" {
		name, _ = route.GetPathTemplate()
		if name != "" {
			name = req.Method + " " + name
		}
	}

	return name
}
