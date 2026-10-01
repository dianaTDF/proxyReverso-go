# Directivas para el Agente (AI)

Al editar código en este repositorio, se deben cumplir las siguientes reglas estrictas:

0. **[CRÍTICO] Mecanismo de Detección de Skills**:
   - Al inicio de cada tarea, evalúa las palabras clave de la petición del usuario.
   - **OBLIGATORIO**: Revisa y lee el archivo `docs/skills/INDEX.md` para determinar si existe un archivo de contexto específico (skill) que debas procesar antes de planificar o escribir código.

1. **Estilo y Modismos Go**:
   - Respetar `gofmt` y `goimports`.
   - Preferir interfaces pequeñas y definidas donde se consumen (Consumer-driven Interfaces).
   - Manejar errores explícitamente. No usar `panic` en el flujo de control normal, devolver error a la función llamadora.
   - Utilizar el patrón "Return early" para reducir la anidación (indentación).

2. **Concurrencia**:
   - Evitar data races. El acceso a recursos compartidos (como la lista de backends) debe estar protegido con `sync.RWMutex`, `sync.Mutex` o `atomic`.
   - Goroutines que realizan tareas periódicas deben aceptar un `context.Context` para ser canceladas limpiamente durante el shutdown.
   - Ninguna goroutine debe quedar "huérfana" o provocar memory leaks.

3. **Arquitectura**:
   - Respetar los módulos indicados (`internal/balancer`, `internal/health`, etc.).
   - No exponer (hacer public) structs o variables a menos que sea estrictamente necesario. Usar minúsculas.
   - Construir dependencias explícitas (e.g. Inyección de dependencias pasando interfaces por los constructores).

4. **Testing**:
   - Todo cambio de lógica central debe estar acompañado de un test (unitario o integrativo).
   - Favorecer la inyección de "Fake Backends" mediante `httptest.Server` para emular requests en tests en lugar de levantar puertos reales a menos que sea inevitable.
