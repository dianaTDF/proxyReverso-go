# Contexto: Testing en Go

Estrategia para las pruebas del repositorio:

- **Table-Driven Tests**: Usar arrays anónimos de structs para definir múltiples casos de prueba. Obligatorio iterar usando `t.Run(name, func(t *testing.T))`.
- **Mocks y Fakes HTTP**: 
  - Para testear el proxy, instanciar backends de prueba utilizando `httptest.NewServer()`.
  - Capturar respuestas del proxy para tests utilizando `httptest.NewRecorder()`.
- **Asserción**: Usar paquetes nativos o `testify/assert` si el repositorio lo adopta, para verificar StatusCode y Headers de respuesta esperados.
- **Test de Concurrencia**: Ejecutar pruebas de carga ligeras directamente dentro del paquete `testing` utilizando goroutines masivas (`sync.WaitGroup`) sobre los handlers + flag `-race`.
