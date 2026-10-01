package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseRecorder intercepta el ResponseWriter nativo para leer el StatusCode emitido por el backend.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger emite telemetría básica (Tiempos y Status) hacia la salida estándar.
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Inicializa StatusOK por defecto por si el backend no emite encabezado explícito.
			rec := &responseRecorder{w, http.StatusOK}
			next.ServeHTTP(rec, r)
			
			log.Printf("[Proxy] %s %s | Status: %d | Latencia: %v", r.Method, r.URL.Path, rec.statusCode, time.Since(start))
		})
	}
}