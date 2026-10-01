# Arquitectura - Reverse Proxy Go

## Visión General
Este documento describe las decisiones arquitectónicas para la implementación del Reverse Proxy en Go. Se optó por utilizar `net/http/httputil.ReverseProxy` como núcleo, lo que permite enfocar el esfuerzo en el ecosistema alrededor del proxy (balanceo, middlewares, observabilidad, etc.) sin reinventar el manejo de conexiones HTTP.

## Componentes y Patrones

### 1. Núcleo del Proxy (`internal/proxy`)
- Base: `httputil.ReverseProxy`.
- Manejo de requests y forwarding a backends.

### 2. Balanceador de Carga (`internal/balancer`)
- Interfaz común: `type Balancer interface { Next() *url.URL }`.
- Estrategias intercambiables: Round-Robin (inicial), Least-Connections (futuro).
- Uso de primitivas de sincronización (`sync.Mutex`, `atomic`) para proteger el estado (backends disponibles).

### 3. Health Checks (`internal/health`)
- Patrón: Worker goroutines no bloqueantes.
- Ejecución periódica (tickers) comprobando la salud de los backends.
- Comunicación de cambios de estado a través de channels o punteros atómicos para modificar el pool del balanceador en tiempo real y de forma segura.

### 4. Middlewares (`internal/middleware`)
- Implementación de un Middleware Chain `type Middleware func(http.Handler) http.Handler`.
- Capas previstas: Logging, Rate Limiting y Auth.

### 5. Configuración y Hot-Reload (`internal/config`)
- Formato: YAML o JSON.
- Recarga en caliente sin downtime ("Zero Downtime Reload").
- Sincronización mediante `sync.RWMutex` para lectura concurrente rápida (en cada request) y escritura segura exclusiva (durante reload).

### 6. Terminación Graceful (`cmd/proxy`)
- Uso de `signal.Notify` para capturar `SIGINT` y `SIGTERM`.
- Cierre ordenado con `server.Shutdown(ctx)` apoyado en `context.Context` con timeout, asegurando el vaciado de las conexiones activas.
