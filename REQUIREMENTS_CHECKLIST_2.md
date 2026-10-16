# REQUIREMENTS CHECKLIST — MultiCatálogo

> Fuente de verdad del progreso. Leer al inicio de cada sesión; marcar `- [x]`
> y actualizar "Estado General" al completar cualquier tarea.
>
> **Estructura:** el documento se divide en fases por **Unidad/Tema del sílabo**
> (`Silabo_Framework_Programacion_Web.md` de referencia). Cada fase trae un comentario con lo
> que se hará en esa clase. El producto integrador es UNO: MultiCatálogo —
> una tienda cosmética con red de ventas multinivel (MLM) que crece como
> **sistema políglota**: React (front) + Go/Fiber (núcleo transaccional, U2)
> + FastAPI (capa SaaS/multitenant, U3), desplegado contenerizado (U4).
>
> **Principios rectores (aplican a toda fase):**
> - Seguridad por defecto: ningún endpoint nuevo nace sin auth ni validación.
> - Inyección de dependencias: interfaces en los bordes; globals solo en el entry point.
> - Funciones puras para el dominio (comisiones, totales, validadores) — testeables sin I/O.
> - Mobile-first: UI a 360px, touch ≥44px, un campo por fila en formularios.

## Estado General

| Ítem | Valor |
|---|---|
| **Hito actual** | **Unidades 1 y 2 completadas** (SPA con sesión JWT persistida + backend seguro: bcrypt, JWT real, rutas protegidas, health con dependencias) |
| **Última tarea completada** | Seguridad U2-T4 (DI NewControllers + bcrypt + JWT real + RequireAuth + identidad por token + FOR UPDATE ACID + health check con ping) (2026-10-16) |
| **Siguiente tarea pendiente** | Commit del trabajo actual; luego U3-T1: scaffold FastAPI `Unidad3_SaaS/` |
| **Bloqueadores** | Ninguno. Servicios en `backend/docker-compose.yml`: PostgreSQL **5432**, Redis **6379**, API **3000** |

## Arquitectura objetivo del semestre

```
React SPA (:5173, proxy /api) ──► Caddy/Nginx (U4: proxy inverso + SSL)
                                    ├─ /api/tx/*   → Go/Fiber  :3000  — núcleo transaccional (U2)
                                    └─ /api/saas/* → FastAPI   :8000  — capa SaaS/multitenant (U3)
                                                        │
                                      PostgreSQL 16 (:5432) + Redis (:6379)
```

---

# UNIDAD 1 — FRAMEWORKS FRONTEND (SPA React + Vite + Tailwind)

## U1-T1 · Introducción a frameworks web y arquitectura de soluciones

> Se presenta la pila tecnológica del semestre, la inversión de control y los
> lineamientos del producto integrador MultiCatálogo (tienda + red MLM).

- [x] Pila definida y workspace estructurado (`frontend/`, `backend/`)

## U1-T2 · Entorno SPA con Vite + Tailwind CSS

> Scaffold del frontend: empaquetador Vite, Tailwind v4, estructura de carpetas
> por responsabilidad (`components/`, `context/`, `services/`, `data/`).

- [x] Entorno SPA con Vite + Tailwind CSS
- [ ] `frontend/package.json` renombrado a `multicatalogo-frontend`
- [ ] `VITE_API_URL` en `frontend/.env` + `import.meta.env` (cliente único en `frontend/src/services/http.ts`)

## U1-T3 · Componentes y gestión de estado global

> Estado global con React Context: sesión de usuario y carrito por usuario.

- [x] `AuthContext` (sesión `{email, rol}`) y `CartContext` (carrito en localStorage por email) — `frontend/src/context/`
- [x] Persistencia de sesión en `frontend/src/context/AuthContext.tsx` (`multicatalogo_sesion`: usuario + JWT en localStorage)
- [ ] Warnings eslint resueltos: split `frontend/src/context/auth.ts` + `frontend/src/context/cart.ts` (hooks/tipos separados del Provider); `npm run lint` → 0 errores, 0 warnings

## U1-T4 · Consumo de API RESTful y enrutamiento

> Promesas con fetch, capa de servicios desacoplada y guardas de ruta por rol.

- [x] `frontend/src/components/Catalogo/services/productosService.ts`, `frontend/src/components/MiRed/services/redService.ts` + guardas `ProtectedRoute`/`AdminRoute` (`frontend/src/App.tsx`)
- [ ] `frontend/src/services/http.ts`: `apiFetch` único con `Authorization: Bearer` automático, `ApiError` tipado con status y helpers de sesión

## U1-T5 · Interfaces full-screen y flujos MLM (Práctica 03)

> Vistas inmersivas del catálogo y flujo completo de compra y red multinivel.

- [x] F8 Login con rol + redirección por rol (`frontend/src/components/Login/Login.tsx`)
- [x] F1 Storefront `/tienda` full-screen (`frontend/src/components/Storefront.tsx`) | F3 Catálogo con búsqueda/filtros (`frontend/src/components/Catalogo/Catalogo.tsx`) | F2 Detalle `/producto/:id` (`frontend/src/components/DetalleProducto.tsx`) | F4 Lightbox
- [x] F5 Carrito + checkout + confirmación (`frontend/src/components/Carrito.tsx`, `Checkout.tsx`, `Confirmacion.tsx`) | F6 Árbol multinivel comisiones 10/5/2 % (`frontend/src/components/MiRed.tsx`) | F7 Dashboard KPIs (`frontend/src/components/Dashboard.tsx`)
- [x] `frontend/src/components/Checkout.tsx` → `POST /api/checkout` real (sin simulación) + manejo de `APIError`
- [ ] Audit mobile-first: 375px sin overflow en Dashboard/MiRed/Catálogo; botones de acción ≥44px (`Añadir al Carrito`, steppers de `Carrito`, `Cerrar Sesión`)

---

# UNIDAD 2 — BACKEND ALTA TRANSACCIONALIDAD (Go/Fiber · "entornos bancarios")

## U2-T1 · Rendimiento, memoria y concurrencia

> Goroutines y contextos: procesos de misión crítica que no bloquean la respuesta HTTP.

- [x] Goroutine `procesarComisiones` tras el COMMIT del pedido (contexto propio, no el del request) — `backend/controllers/checkoutController.go`

## U2-T2 · Microservicios y APIs robustas (Fiber)

> API REST en capas con contrato de error estándar y configuración por entorno.

- [x] Capas `backend/routes` → `backend/controllers` → `backend/repository`; endpoints login/productos/red/health
- [x] `APIError{status,message,details}` (`backend/models/models.go`) en todos los errores + validación de entrada (400/404)
- [x] Puerto de escucha configurable (`PORT`, defecto 3000) en `backend/main.go`
- [x] `/api/health` real: `HealthController` (`backend/controllers/healthController.go`) con ping a Postgres (503 si cae) y Redis (`"degraded"` si cae)

## U2-T3 · Persistencia PostgreSQL + caché Redis

> Dockeriza la BD, SQL puro con pgx, reconstrucción del árbol y Cache-Aside.

- [x] `backend/docker-compose.yml` Postgres 16 + Redis alpine + volumen + healthchecks
- [x] `backend/database/init.sql` (usuarios/productos/referidos + seed), `backend/repository/repository.go` a SQL real, árbol por `parent_id`
- [x] Cache-Aside (`backend/cache/cache.go`): `productos:all` TTL 5 min, `red:tree:{email}` TTL 10 min, degradación sin caché
- [ ] Índices: `referidos(parent_id)`, `pedidos(cliente_email)`, `pedido_items(pedido_id)`, `productos(categoria)` (migración `backend/database/003_seguridad.sql` + `init.sql`)
- [ ] Errores propagados: firmas del repo devuelven `error`; `ErrNoEncontrado`→404/401, resto→500 (nunca lista vacía silenciosa)

## U2-T4 · Inyección de dependencias y microservicio seguro

> DI por interfaces en toda la cadena + transacción ACID del checkout + hardening.

- [x] Interfaces `DBRepository`/`CacheRepository`; `PostgresRepository`/`RedisCache`; structs de controladores; `NewControllers` en `backend/main.go`
- [x] `POST /api/checkout` ACID: `FOR UPDATE` + descuento de stock + `pedidos`/`pedido_items`, ROLLBACK → 400/500, total recalculado en servidor
- [x] Invalidación de caché tras compra (`productos:all`, `red:tree:*`) + migración `backend/database/002_checkout.sql`
- [x] bcrypt de passwords (pgcrypto `crypt()/gen_salt('bf')` en migración `backend/database/003_seguridad.sql` + seed; login con `CompareHashAndPassword`)
- [x] JWT real (`JWT_SECRET` en env, claims `{sub,email,rol,iat,exp}` 2h, HS256) + `middleware.RequireAuth(secret)`/`RequireAdmin()` en `backend/middleware/`
- [x] Protegidos `/api/red` y `/api/checkout` con el middleware (401 sin token / token adulterado)
- [x] Checkout toma la identidad del TOKEN (`c.Locals` desde claims) — verificado: `email` falso en body NO se persiste ni atribuye comisiones
- [x] Errores sanitizados: `Details` solo con `APP_ENV=dev`; `ErrCantidadInvalida` propio (400 "cantidad debe ser mayor que cero", distinto de stock)

---
