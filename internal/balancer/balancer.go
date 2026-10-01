package balancer

import (
	"net/url"
	"sync"
	"sync/atomic"
)

// Balancer define la interfaz genérica para selección de backends con callback de finalización.
type Balancer interface {
	Next() (target *url.URL, done func())
	SetBackends(urls []*url.URL)
}

// RoundRobin implementa Balancer distribuyendo equitativamente sin bloqueos de Mutex en lectura.
type RoundRobin struct {
	mu       sync.RWMutex
	backends []*url.URL
	current  uint32
}

func NewRoundRobin(urls []*url.URL) *RoundRobin {
	return &RoundRobin{
		backends: urls,
		current:  0,
	}
}

func (r *RoundRobin) Next() (*url.URL, func()) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.backends) == 0 {
		return nil, func() {}
	}

	next := atomic.AddUint32(&r.current, 1)
	idx := (int(next) - 1) % len(r.backends)
	return r.backends[idx], func() {} // RoundRobin no necesita trackear conexiones
}

func (r *RoundRobin) SetBackends(urls []*url.URL) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends = urls
}
