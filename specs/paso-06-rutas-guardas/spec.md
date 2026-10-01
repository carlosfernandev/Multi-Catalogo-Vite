# Paso 06 — Rutas y guardas por rol

**Entregable:** base
**Correspondencia con la planificación:** ACD 2
**Fuente:** `practica03.md` → Paso 6

## Objetivo

Proteger el ruteo con dos guardas anidadas y asegurar aislamiento del carrito por usuario. Este paso define el mapa de accesos que el resto de la práctica asume.

## Concepto

Dos guardas anidadas: `ProtectedRoute` exige autenticación; `AdminRoute` exige rol `admin` (cliente es redirigido a `/tienda`). `CartBoundary` remonta `CartProvider` con `key={user?.email}` para que cada cuenta tenga su propio carrito.

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/App.tsx` |

## Alcance técnico

- `ProtectedRoute` → `Navigate to="/login"` si no autenticado; `AdminRoute` → `Navigate to="/tienda"` si no es admin.
- `CartBoundary` envuelve `CartProvider` con `key={user?.email ?? "anonimo"}`.
- Rutas: `/login` pública; `/` y `/mi-red` solo admin; `/tienda`, `/catalogo`, `/producto/:id`, `/carrito`, `/checkout`, `/confirmacion` para ambos roles.
- Ver código en `practica03.md`, Paso 6.

## Criterios de aceptación

- [x] Usuario no autenticado es redirigido a `/login` en cualquier ruta protegida.
- [x] Cliente que navega a `/` o `/mi-red` es redirigido a `/tienda`.
- [x] Admin accede a todas las rutas.
- [x] `CartBoundary` cambia `key` al cambiar de usuario.

## Verificación

Probar manualmente: login admin → acceso a `/mi-red`; login cliente → redirección desde `/mi-red` a `/tienda`; sin login → redirección a `/login`. `npm run build`.
