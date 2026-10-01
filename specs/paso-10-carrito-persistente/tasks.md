# Tareas — Paso 10

> Tareas pequeñas y atómicas: una responsabilidad por tarea, en orden de ejecución. Marcar `- [x]` al completar.

- [x] **T1**: Reemplazar `frontend/src/context/CartContext.tsx` con `CartItem`, `CART_STORAGE_PREFIX` y `storageKey` por email — *Completado cuando:* define tipos y clave por usuario.
- [x] **T2**: Implementar `useState` lazy init desde `localStorage` y `useEffect` de persistencia — *Completado cuando:* lee y escribe en `localStorage` en cada cambio.
- [x] **T3**: Agregar `incrementQuantity`, `decrementQuantity`, `clearCart` y derivados `totalItems`/`totalPrice` — *Completado cuando:* todas las funciones existen y `decrement` filtra en 0.
- [x] **T4**: Reemplazar `frontend/src/components/Carrito.tsx` con controles `− / +` y resumen — *Completado cuando:* muestra cantidad, subtotal y link a `/checkout`.
- [ ] **T5**: Probar persistencia y aislamiento por usuario — *Completado cuando:* recarga persiste y cuentas distintas no se mezclan.
