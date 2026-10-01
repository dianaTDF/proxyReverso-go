# Contexto: Arquitectura y Diseño

Principios de acoplamiento y estructura:

- **Separación de Responsabilidades (SoC)**: 
  - El **Balanceador de Carga** (`Balancer`) solo devuelve la siguiente URL o Host. NO sabe de peticiones HTTP, no lee body, no escribe responses.
  - El **Health Checker** verifica conectividad. Su única tarea es marcar si una URL está viva o muerta y notificárselo al Balanceador.
  - El **Manejador de Proxy** orquesta: Toma request -> Pide URL al Balanceador -> Ejecuta proxy -> Retorna respuesta.
- **Inyección de Dependencias**: Pasa dependencias a través de structs/constructores (`NewProxy(b Balancer, c Config)`), no uses variables globales.
- **Configuración**: El Hot-Reload de la configuración debe re-construir el pool del balanceador de manera atómica o bajo un Lock exclusivo (`Lock()`), mientras el hot path usa el Lock de lectura (`RLock()`).
