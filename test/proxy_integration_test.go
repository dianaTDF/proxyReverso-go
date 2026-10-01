package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"proxyRev/internal/balancer"
	"proxyRev/internal/proxy"
)

func TestProxyRouting(t *testing.T) {
	// 1. Inicializar Fake Backends (Simulan servidores de aplicación)
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SERVER_1"))
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SERVER_2"))
	}))
	defer backend2.Close()

	// 2. Ensamblar Proxy
	u1, _ := url.Parse(backend1.URL)
	u2, _ := url.Parse(backend2.URL)
	urls := []*url.URL{u1, u2}

	rrBalancer := balancer.NewRoundRobin(urls)
	revProxy := proxy.New(rrBalancer)

	// Usamos httptest.NewServer para levantar el proxy sin ocupar puertos estáticos
	proxyServer := httptest.NewServer(revProxy)
	defer proxyServer.Close()

	client := proxyServer.Client()

	// 3. Ejecutar Asersiones HTTP
	// Request 1 -> Debería ir a SERVER_1
	resp1, err := client.Get(proxyServer.URL)
	if err != nil {
		t.Fatalf("Request 1 falló: %v", err)
	}
	defer resp1.Body.Close()
	
	body1, _ := io.ReadAll(resp1.Body)
	if string(body1) != "SERVER_1" {
		t.Errorf("Error Round-Robin. Esperaba SERVER_1, obtuve: %s", string(body1))
	}

	// Request 2 -> Debería ir a SERVER_2
	resp2, err := client.Get(proxyServer.URL)
	if err != nil {
		t.Fatalf("Request 2 falló: %v", err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	if string(body2) != "SERVER_2" {
		t.Errorf("Error Round-Robin. Esperaba SERVER_2, obtuve: %s", string(body2))
	}
}
