package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_requests_total",
			Help: "Total de peticiones procesadas por el reverse proxy.",
		},
		[]string{"method", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "proxy_request_duration_seconds",
			Help:    "Distribucion de latencia de peticiones en segundos.",
			Buckets: prometheus.DefBuckets, // Buckets estándar (.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)
		},
		[]string{"method"},
	)
)

// Metrics recolecta métricas de tráfico exponiéndolas a Prometheus.
// Reutiliza 'responseRecorder' definido en logger.go (mismo package).
func Metrics() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Intercepta para capturar el status emitido por capas interiores
			rec := &responseRecorder{w, http.StatusOK}
			
			next.ServeHTTP(rec, r)

			duration := time.Since(start).Seconds()

			httpRequestsTotal.WithLabelValues(r.Method, strconv.Itoa(rec.statusCode)).Inc()
			httpRequestDuration.WithLabelValues(r.Method).Observe(duration)
		})
	}
}
