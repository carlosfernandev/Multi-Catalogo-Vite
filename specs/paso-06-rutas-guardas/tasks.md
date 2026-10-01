# Tareas — Paso 06

> Tareas pequeñas y atómicas: una responsabilidad por tarea, en orden de ejecución. Marcar `- [x]` al completar.

- [x] **T1**: Crear `ProtectedRoute` en `frontend/src/App.tsx` con `Outlet` y `Navigate` — *Completado cuando:* exige `isAuthenticated`.
- [x] **T2**: Crear `AdminRoute` que verifica `user?.rol === 'admin'` — *Completado cuando:* no-admin es redirigido a `/tienda`.
- [x] **T3**: Crear `CartBoundary` con `key={user?.email}` que envuelve `CartProvider` — *Completado cuando:* cambia `key` al cambiar usuario.
- [x] **T4**: Reconfigurar `Routes`: `AuthProvider > CartBoundary > BrowserRouter` y anidar `AdminRoute` vs flujo de compra — *Completado cuando:* el mapa de rutas coincide con la tabla de la guía.
- [ ] **T5**: Verificar guardas y compilar — *Completado cuando:* `npm run build` y pruebas manuales de redirección pasan.
