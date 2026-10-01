# Contexto: Concurrencia en Go

Reglas estrictas para manejar concurrencia en este proyecto:

- **Goroutines y Ciclo de Vida**: Toda goroutine persistente o de background (ej. health checks) DEBE recibir un `context.Context` y terminar limpiamente cuando este es cancelado.
- **Race Conditions**: Nunca compartir mapas o slices sin protección. Utilizar primitivas del paquete `sync`:
  - `sync.RWMutex` para estructuras de lectura frecuente (ej. lista de backends activos). Bloquear lectura con `RLock()` en el hot path.
  - Tipos `atomic` (ej. `atomic.Pointer`) si aplica para cambios de punteros libres de lock.
- **Validación Obligatoria**: Todo código de concurrencia debe ser ejecutado y validado localmente con `go run -race` o `go test -race` antes de considerarse funcional.
- **Channels**: Usarlos preferentemente para señalización y pipelines, no como simples mutexes disfrazados. Cerrar channels desde el sender.
