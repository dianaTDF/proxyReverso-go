# High-Performance Reverse Proxy (Go)

Un proxy inverso de alto rendimiento implementado en Go. Diseñado como un proyecto educativo/portfolio que demuestra el uso avanzado de primitivas de concurrencia, interfaces idiomáticas y seguridad en ambientes altamente concurrentes.

## Características Principales

- **Balanceo de Carga Multi-Backend**: Estrategia Round-Robin asíncrona implementada con `sync/atomic` para lograr operaciones *Lock-Free* y maximizar el throughput.
- **Active Health Checks**: Trabajadores (Workers) evaluando periódicamente la salud de los servidores en background. Evitan latencia inyectando mutaciones seguras al pool de tráfico mediante `sync.RWMutex`.
- **Hot Reloading (Zero Downtime)**: Recarga y aplicación en caliente del archivo de configuración sin tirar el listener HTTP, accionado mediante señales del sistema operativo (`SIGHUP`).
- **Cadena de Middlewares Transversales**:
  - **Rate Limiting**: Protección L7 automatizada basada en IP, con mitigador de memory leaks (Goroutine GC).
  - **Logging y Telemetría**: Observabilidad estricta interceptando las peticiones nativas para extraer códigos de estado de forma agnóstica.
- **Exportación de Métricas**: Integración nativa con Prometheus. `CounterVec` y `HistogramVec` expuestos en un puerto privado/auxiliar (`:9090`) por seguridad.
- **Graceful Shutdown**: Drenado seguro de peticiones activas mediado por canalizaciones temporizadas (`context.Context`), interceptando `SIGTERM` y `SIGINT`.

## Requisitos

- **Go**: 1.21 o superior.

## Instrucciones de Uso

### 1. Configuración del Entorno

Modifica o crea el archivo `config.json` en la raíz del proyecto. Aquí se definen el puerto de escucha y el array de backends destino:

```json
{
  "server": {
    "port": "8080"
  },
  "proxy": {
    "backends": [
      "http://127.0.0.1:8081",
      "http://127.0.0.1:8082"
    ]
  }
}
```

### 2. Arranque del Servidor

Para compilar y correr el proxy:

```bash
# Compilar binario de producción
go build -o reverse_proxy ./cmd/proxy

# Ejecutar
./reverse_proxy
```

O para desarrollo rápido:
```bash
go run ./cmd/proxy/main.go
```

**Puertos Activos**:
- `:8080` - Tráfico de usuarios (Enrutado al backend correspondiente).
- `:9090` - Tráfico de Observabilidad (Ruta `/metrics` para que el scraper de Prometheus lo consuma).

### 3. Operativa Avanzada (Zero-Downtime)

#### Hot-Reloading de la Configuración
Si necesitas añadir o quitar un backend de la infraestructura:
1. Modifica `config.json` y guarda el archivo.
2. Descubre el PID del proceso y envíale la señal `SIGHUP`.
   ```bash
   # Buscar PID (ej. si usaste go run)
   pgrep -f "proxy"
   
   # Mandar señal de recarga
   kill -HUP <PID>
   ```
El proxy recargará el rotador de balanceo de inmediato sin rechazar ninguna de las conexiones HTTP en vuelo.

#### Apagado Seguro
Para detener el proceso sin corromper respuestas que estén en procesamiento:
Presiona `Ctrl+C` en la terminal, o envía un `kill -TERM <PID>`. El proceso rechazará nuevas conexiones inmediatamente y esperará un máximo de 5 segundos a que todos los backends terminen de contestar.
