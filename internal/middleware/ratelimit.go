package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

// RateLimiter aplica protección L7 limitando peticiones por IP usando Token Bucket simplificado.
func RateLimiter(limit int, window time.Duration) Middleware {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Worker limpiador de memoria: previene Memory Leaks de IPs que no regresan.
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			now := time.Now()
			for ip, v := range visitors {
				if now.Sub(v.lastSeen) > window {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]

			mu.Lock()
			v, exists := visitors[ip]
			if !exists || time.Since(v.lastSeen) > window {
				// Nueva IP o expiró la ventana de tiempo
				visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
			} else {
				v.count++
				v.lastSeen = time.Now()
			}
			count := visitors[ip].count
			mu.Unlock() // Liberamos Mutex ASAP (Return Early)

			if count > limit {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
