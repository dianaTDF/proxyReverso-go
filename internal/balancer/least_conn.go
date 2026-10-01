package balancer

import (
	"net/url"
	"sync"
)

type backendNode struct {
	url         *url.URL
	activeConns int64
}

// LeastConnections enruta el tráfico al nodo con menor carga concurrente.
type LeastConnections struct {
	mu       sync.Mutex
	backends []*backendNode
}

func NewLeastConnections(urls []*url.URL) *LeastConnections {
	lc := &LeastConnections{}
	lc.SetBackends(urls)
	return lc
}

func (l *LeastConnections) Next() (*url.URL, func()) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.backends) == 0 {
		return nil, func() {}
	}

	var selected *backendNode
	for _, b := range l.backends {
		if selected == nil || b.activeConns < selected.activeConns {
			selected = b
		}
	}

	// Registrar conexión in-flight
	selected.activeConns++
	target := selected.url

	// Callback capturado para reducir contador
	done := func() {
		l.mu.Lock()
		selected.activeConns--
		l.mu.Unlock()
	}

	return target, done
}

func (l *LeastConnections) SetBackends(urls []*url.URL) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var newBackends []*backendNode
	for _, u := range urls {
		newBackends = append(newBackends, &backendNode{url: u, activeConns: 0})
	}
	l.backends = newBackends
}
