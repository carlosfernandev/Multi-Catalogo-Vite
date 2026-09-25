# Paso 10 — Carrito persistente por usuario

**Entregable:** F5
**Correspondencia con la planificación:** APD
**Fuente:** `practica03.md` → Paso 10

## Objetivo

Hacer el carrito persistente en `localStorage` y aislado por usuario, con controles de cantidad. Cada cuenta debe recuperar su propio carrito al iniciar sesión.

## Concepto

El carrito se inicializa leyendo `localStorage` y se guarda en cada cambio con `useEffect`. La clave incluye el correo del usuario; `CartBoundary` (Paso 6) remonta el provider con `key` al cambiar de cuenta. Se agregan `incrementQuantity`, `decrementQuantity` y `clearCart`.

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/context/CartContext.tsx` |
| reemplazar | `frontend/src/components/Carrito.tsx` |

## Alcance técnico

- `CartItem extends IProducto { cantidad }`; `CART_STORAGE_PREFIX = "multicatalogo_carrito_"`.
- `storageKey` por `user?.email`; `useState` lazy init desde `localStorage`; `useEffect` persiste en cada cambio.
- `addToCart` incrementa si existe; `decrementQuantity` filtra si cantidad llega a 0.
- `Carrito.tsx` con controles `− / +`, subtotales y link a `/checkout`.
- Ver código en `practica03.md`, Paso 10.

## Criterios de aceptación

- [ ] Recargar la página conserva el carrito.
- [ ] Cambiar de usuario muestra carrito vacío o el suyo propio (no mezclado).
- [ ] `−` a 0 elimina el ítem; `+` incrementa cantidad.
- [ ] `clearCart` vacía el carrito (usado por checkout).

## Verificación

Agregar productos, recargar (persiste), cambiar de usuario (cada cuenta tiene el suyo), y usar los controles − / +.
