# Router de Skills (Índice Maestro)

Este archivo actúa como un mecanismo de detección de contexto automático para el agente.
**Agente**: Antes de comenzar a escribir o modificar código, evalúa la intención de la petición del usuario y busca coincidencias con las palabras clave a continuación. Si hay coincidencia, DEBES leer el archivo referenciado antes de ejecutar.

| Intención / Palabras Clave | Archivo de Contexto a Leer |
|---|---|
| `goroutine`, `channel`, `race`, `mutex`, `sync`, `atomic`, concurrencia, hilos | [concurrency.md](file:///var/www/html/otros/proxyRev/docs/skills/concurrency.md) |
| `http`, `proxy`, `headers`, `keep-alive`, `status`, reescritura, ruteo | [http_proxy.md](file:///var/www/html/otros/proxyRev/docs/skills/http_proxy.md) |
| arquitectura, diseño, estructura, balanceador, responsabilidades, desacoplamiento | [architecture.md](file:///var/www/html/otros/proxyRev/docs/skills/architecture.md) |
| `test`, `mock`, `httptest`, unitario, integración, coverage, tdd | [testing.md](file:///var/www/html/otros/proxyRev/docs/skills/testing.md) |
