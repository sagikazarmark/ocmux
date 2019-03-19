package ocmux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type testTraceExporter struct {
	spans []*trace.SpanData
}

func (t *testTraceExporter) ExportSpan(s *trace.SpanData) {
	t.spans = append(t.spans, s)
}

type testStatsExporter struct {
	viewData []*view.Data
}

func (t *testStatsExporter) ExportView(d *view.Data) {
	t.viewData = append(t.viewData, d)
}

func TestMiddleware_RouteName(t *testing.T) {
	tests := map[string]struct {
		configure    func(r *mux.Router)
		url          string
		expectedName string
	}{
		"route_name": {
			configure: func(r *mux.Router) {
				r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(204)
				})).Name("route_name")
			},
			url:          "/",
			expectedName: "route_name",
		},
		"path_template": {
			configure: func(r *mux.Router) {
				r.Handle("/{id}/whatever", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(204)
				}))
			},
			url:          "/123/whatever",
			expectedName: "GET /{id}/whatever",
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			v := &view.View{
				Name:        "request_total",
				Measure:     ochttp.ServerLatency,
				Aggregation: view.Count(),
				TagKeys:     []tag.Key{ochttp.KeyServerRoute},
			}
			_ = view.Register(v)

			var te testTraceExporter
			trace.RegisterExporter(&te)
			defer trace.UnregisterExporter(&te)

			var se testStatsExporter
			view.RegisterExporter(&se)
			defer view.UnregisterExporter(&se)

			router := mux.NewRouter()

			router.Use(Middleware())

			test.configure(router)

			plugin := ochttp.Handler{
				Handler: router,
				StartOptions: trace.StartOptions{
					Sampler: trace.AlwaysSample(),
				},
			}
			req, _ := http.NewRequest("GET", test.url, nil)
			rr := httptest.NewRecorder()
			plugin.ServeHTTP(rr, req)

			view.Unregister(v) // trigger exporting

			if got, want := se.viewData[0].Rows[0].Tags[0].Value, test.expectedName; got != want {
				t.Errorf("tag name does not match the expected one\ngot:  %s\nwant: %v", got, want)
			}

			if got, want := te.spans[0].Name, test.expectedName; got != want {
				t.Errorf("span name does not match the expected one\ngot:  %s\nwant: %v", got, want)
			}
		})
	}
}
