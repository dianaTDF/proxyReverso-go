# Reverse Proxy in Go

Un proxy inverso de alto rendimiento implementado en Go, diseñado para servir como un proyecto educativo/portfolio que demuestra el uso avanzado de primitivas de concurrencia, interfaces y diseño idiomático en Go.

## Características (Ver `ROADMAP.md` para estado)
- Balanceo de carga (Round-Robin, Least-Connections).
- Middlewares (Logging, Rate Limit, Auth).
- Zero-Downtime Reload de configuración.
- Active Health Checks mediante goroutines.
- Soporte para exportar métricas a Prometheus.
- Graceful Shutdown.

## Requisitos
- [Go](https://go.dev/) 1.21+ (recomendado).

## Inicio Rápido (Próximamente)

El servidor arrancará leyendo la configuración inicial y quedará a la espera de requests.

---

---
# Roadmap - Reverse Proxy Go

Este roadmap define las fases de desarrollo iterativas para el proyecto. Cada etapa debe estar respaldada por tests funcionales y mantener la coherencia del diseño en Go.

## Etapa 1: MVP Funcional
- [ ] Inicializar módulo de Go y estructura de carpetas.
- [ ] Parseo básico de configuración (YAML o JSON) para una ruta a un único backend.
- [ ] Implementación de `httputil.ReverseProxy` para reenviar peticiones.
- [ ] Integración del log básico.
- [ ] Implementar Graceful Shutdown.

## Etapa 2: Balanceo de Carga y Multi-Backend
- [ ] Ampliar configuración para aceptar múltiples backends por ruta.
- [ ] Diseñar interfaz `Balancer`.
- [ ] Implementar estrategia "Round-Robin".
- [ ] Implementar estrategia "Least-Connections" y/o "Weighted".

## Etapa 3: Active Health Checks
- [ ] Goroutines periódicas (pings HTTP) a los backends.
- [ ] Lógica para sacar/introducir backends del rotador si fallan.
- [ ] Uso de primitives de concurrencia (`sync`, channels) para evitar Data Races en la lista de backends del balanceador.

## Etapa 4: Cadena de Middlewares
- [ ] Implementar mecanismo estándar para enlazar middlewares (Middleware Chain).
- [ ] Middleware: Logging estructurado (tiempos de respuesta, status codes).
- [ ] Middleware: Rate Limiting por IP o Header.
- [ ] Middleware: Inyección de Auth Headers o Validación simple.

## Etapa 5: Configuración Dinámica (Hot Reload)
- [ ] Endpoint de control (ej. `POST /_/reload`) o señal del SO (`SIGHUP`).
- [ ] Recarga de archivo de configuración en memoria utilizando `sync.RWMutex` o punteros atómicos para evitar bloqueo o caída de peticiones concurrentes.

## Etapa 6: Observabilidad (Métricas)
- [ ] Exponer endpoint `/metrics` en formato Prometheus.
- [ ] Contadores (ej. total de peticiones por ruta/backend).
- [ ] Histogramas (ej. latencia de los request al backend).
