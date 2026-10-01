# Contexto: HTTP y Proxy Protocol

Reglas y conceptos de red para el desarrollo del Proxy:

- **Base del Proxy**: Usar siempre `net/http/httputil.ReverseProxy` como cimiento. No reinventar el enrutamiento base ni el manejo de conexiones TCP.
- **Headers Sensibles**: 
  - Al reenviar peticiones, asegúrate de preservar o inyectar correctamente `X-Forwarded-For`, `X-Real-IP`, y ajustar el header `Host` según el backend.
  - El `Director` (o el reverse proxy nativo) se encarga de reescribir la URL.
- **Timeouts**: Todo `http.Server` instanciado DEBE tener `ReadTimeout`, `WriteTimeout` y `IdleTimeout` explícitamente configurados para evitar file descriptor leaks (ataques tipo Slowloris).
- **Transportes Personalizados**: Si necesitas customizar el manejo de `Keep-Alive`, TLS, o limitar el pool de conexiones, se debe crear un custom `http.Transport` y asignarlo al ReverseProxy.
