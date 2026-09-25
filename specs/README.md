# Specs — Práctica 03 Multi-Catálogo

Specs atómicas derivadas de `practica03.md`. Cada paso tiene un `spec.md` (qué y por qué) y un `tasks.md` (cómo, en orden).

## Índice de pasos

| Paso | Carpeta | Entregable | Archivos principales |
|---|---|---|---|
| 1 — Rol en la API | [paso-01-rol-api](paso-01-rol-api/) | base de F8 | `backend/controllers/authController.go` |
| 2 — Datos de productos | [paso-02-datos-productos](paso-02-datos-productos/) | base de F1, F3, F4 | `components/Catalogo/interfaces/iproducto.ts`, `components/Catalogo/data/productos.ts` |
| 3 — Datos de la red MLM | [paso-03-datos-red-mlm](paso-03-datos-red-mlm/) | base de F6, F7 | `components/MiRed/data/red.ts` |
| 4 — Servicio de productos | [paso-04-servicio-productos](paso-04-servicio-productos/) | base de F1, F3 | `components/Catalogo/services/productosService.ts` |
| 5 — Autenticación con rol | [paso-05-autenticacion-rol](paso-05-autenticacion-rol/) | F8 | `context/AuthContext.tsx`, `components/Login/services/authService.ts`, `components/Login/Login.tsx` |
| 6 — Rutas y guardas | [paso-06-rutas-guardas](paso-06-rutas-guardas/) | base | `App.tsx` (`ProtectedRoute`, `AdminRoute`, `CartBoundary`) |
| 7 — Storefront | [paso-07-storefront](paso-07-storefront/) | F1 | `components/Storefront.tsx` |
| 8 — Catálogo con filtros | [paso-08-catalogo-filtros](paso-08-catalogo-filtros/) | F3 | `components/Catalogo/Catalogo.tsx` |
| 9 — Detalle y lightbox | [paso-09-detalle-lightbox](paso-09-detalle-lightbox/) | F2, F4 | `components/DetalleProducto.tsx` |
| 10 — Carrito persistente | [paso-10-carrito-persistente](paso-10-carrito-persistente/) | F5 | `context/CartContext.tsx`, `components/Carrito.tsx` |
| 11 — Checkout y confirmación | [paso-11-checkout-confirmacion](paso-11-checkout-confirmacion/) | F5 | `components/Checkout.tsx`, `components/Confirmacion.tsx` |
| 12 — Red multinivel | [paso-12-red-multinivel](paso-12-red-multinivel/) | F6 | `components/MiRed.tsx` |
| 13 — Dashboard con KPIs | [paso-13-dashboard-kpis](paso-13-dashboard-kpis/) | F7 | `components/Dashboard.tsx` |
| 14 — Navegación por rol | [paso-14-navegacion-rol](paso-14-navegacion-rol/) | F8 | `components/Sidebar.tsx`, `components/Navbar.tsx` |

## Cómo usar estas specs

1. Seguí el orden 1 → 14; cada paso asume que el anterior está completo.
2. Leé `spec.md` como referencia (objetivo, alcance, criterios).
3. Ejecutá las tareas de `tasks.md` en orden y marcá `- [x]` al completar.
4. El código completo vive en `practica03.md` (fuente de verdad); estas specs solo lo resumen y atomizan.
5. Verificá cada paso con el bloque **Verificación** antes de avanzar.

> Convención: docs en español; identificadores y rutas en inglés tal cual en el código.
