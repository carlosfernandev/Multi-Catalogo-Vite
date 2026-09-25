# Tareas — Paso 05

> Tareas pequeñas y atómicas: una responsabilidad por tarea, en orden de ejecución. Marcar `- [x]` al completar.

- [ ] **T1**: Reemplazar `frontend/src/context/AuthContext.tsx` con `Rol`, `Usuario` y nuevo `login` — *Completado cuando:* el contexto expone `user` y `isAuthenticated`.
- [ ] **T2**: Crear `frontend/src/components/Login/services/authService.ts` con `login` + normalización de rol — *Completado cuando:* hace `fetch /api/login` y retorna `{ email, rol }`.
- [ ] **T3**: Mover y reemplazar `frontend/src/components/Login/Login.tsx` (feature Login) — *Completado cuando:* importa `login as loginRequest` y redirige por rol.
- [ ] **T4**: Eliminar `frontend/src/components/Login.tsx` anterior si existe — *Completado cuando:* solo queda `Login/Login.tsx`.
- [ ] **T5**: Verificar tipado y compilar — *Completado cuando:* `npm run build` sin errores.
- [ ] **T6**: Probar login admin y cliente — *Completado cuando:* admin navega a `/` y cliente a `/tienda`.
