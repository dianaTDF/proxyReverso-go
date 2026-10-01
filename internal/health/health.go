package health

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"proxyRev/internal/balancer"
)

// Checker ejecuta goroutines periódicas para evaluar la salud de backends.
type Checker struct {
	mu          sync.RWMutex
	allBackends []*url.URL
	balancer    balancer.Balancer
	interval    time.Duration
}

// NewChecker inicializa el worker de chequeos.
func NewChecker(urls []*url.URL, bal balancer.Balancer, interval time.Duration) *Checker {
	return &Checker{
		allBackends: urls,
		balancer:    bal,
		interval:    interval,
	}
}

// UpdateBackends renueva la lista de servidores a observar de forma Thread-Safe.
func (c *Checker) UpdateBackends(urls []*url.URL) {
	c.mu.Lock()
	c.allBackends = urls
	c.mu.Unlock()
	
	// Fuerza un chequeo inmediato fuera del Lock exclusivo
	c.check()
}

// Start inicia un loop infinito protegido por context.Context
func (c *Checker) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	c.check()

	for {
		select {
		case <-ctx.Done():
			log.Println("[Health] Worker detenido correctamente.")
			return
		case <-ticker.C:
			c.check()
		}
	}
}

func (c *Checker) check() {
	c.mu.RLock()
	targets := make([]*url.URL, len(c.allBackends))
	copy(targets, c.allBackends)
	c.mu.RUnlock()

	var healthy []*url.URL
	client := http.Client{Timeout: 2 * time.Second}

	for _, u := range targets {
		resp, err := client.Get(u.String())
		if err != nil || resp.StatusCode >= 500 {
			log.Printf("[Health] Backend caido: %s", u.String())
			continue
		}
		
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		
		healthy = append(healthy, u)
	}

	c.balancer.SetBackends(healthy)
}
