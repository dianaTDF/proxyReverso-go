package balancer

import (
	"net/url"
	"sync"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	u1, _ := url.Parse("http://back1")
	u2, _ := url.Parse("http://back2")
	u3, _ := url.Parse("http://back3")

	backends := []*url.URL{u1, u2, u3}
	rr := NewRoundRobin(backends)

	// Validar secuencia Round-Robin exacta y callbacks mock
	u, d := rr.Next()
	d()
	if u != u1 {
		t.Errorf("Esperaba %v, obtuve %v", u1, u)
	}
	u, d = rr.Next()
	d()
	if u != u2 {
		t.Errorf("Esperaba %v, obtuve %v", u2, u)
	}
	u, d = rr.Next()
	d()
	if u != u3 {
		t.Errorf("Esperaba %v, obtuve %v", u3, u)
	}
	u, d = rr.Next()
	d()
	if u != u1 {
		t.Errorf("Esperaba reinicio en %v, obtuve %v", u1, u)
	}
}

func TestRoundRobinConcurrency(t *testing.T) {
	urls := []*url.URL{
		{Host: "b1"}, {Host: "b2"}, {Host: "b3"},
	}
	rr := NewRoundRobin(urls)

	var wg sync.WaitGroup
	workers := 100
	requestsPerWorker := 1000

	// Simular lectura concurrente masiva (100,000 requests)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerWorker; j++ {
				_, d := rr.Next()
				d()
			}
		}()
	}

	// Simular escritura concurrente (Hot-Reload de backends simulado)
	go func() {
		for i := 0; i < 50; i++ {
			rr.SetBackends(urls)
		}
	}()

	wg.Wait()
	// Si finaliza sin panic ni data race reportado por -race, es seguro.
}
