# Tareas — Paso 01

> Tareas pequeñas y atómicas: una responsabilidad por tarea, en orden de ejecución. Marcar `- [x]` al completar.

- [ ] **T1**: Actualizar `backend/controllers/authController.go` con el `switch` de dos roles y respuesta `rol` — *Completado cuando:* el archivo contiene ambos casos y el `default` con 401.
- [ ] **T2**: Verificar imports (`fiber`, `backend/models`) sin cambios innecesarios — *Completado cuando:* `go vet` / `go build` no reporta imports no usados.
- [ ] **T3**: Compilar el backend — *Completado cuando:* `cd backend && go build ./...` termina sin errores.
- [ ] **T4**: Probar login admin con curl — *Completado cuando:* `curl POST /api/login` con admin retorna `rol: "admin"` y token.
- [ ] **T5**: Probar login cliente y credenciales inválidas — *Completado cuando:* cliente retorna `rol: "cliente"` e inválidas retornan 401.
