# PROJECT_MAP — Mapa Arquitectónico del Proyecto

## Estructura General

- `backend/` — API Go/Fiber (módulo Go: `backend`). Puerto 3000.
- `frontend/` — React 19 + Vite + Tailwind CSS 4 + React Router v7 + TypeScript. Puerto 5173.
- `specs/` — Specs atómicas de la Práctica 03 (14 pasos atomizados).
- `practica03.md` — Guía de laboratorio completa (fuente de verdad del código).

## Backend (Go/Fiber)

**Raíz del módulo:** `backend/main.go` — Entrada, configuración Fiber, CORS, setup de rutas.

**Controladores** (`backend/controllers/`):
- `authController.go` — `POST /api/login` (login con rol: admin/cliente).
- `prodController.go` — `GET /api/productos` (lista), `GET /api/productos/:id` (detalle por ID).
- `healthController.go` — `GET /api/health` (health check).
- `redController.go` — `GET /api/red` (árbol multinivel, caché `red:tree:{email}` TTL 10 min).
- `prodController.go` usa caché Redis `productos:all` (TTL 5 min) en `GET /api/productos`.

**Caché** (`backend/cache/`, `backend/config/redis.go`):
- `config/redis.go` — Cliente Redis (`REDIS_URL` o `REDIS_HOST`/`REDIS_PORT`/`REDIS_PASSWORD`).
- `cache/cache.go` — `Obtener` / `Guardar` (Cache-Aside JSON + TTL; cabecera `X-Cache: HIT|MISS`).

**Modelos** (`backend/models/`):
- `models.go` — `LoginRequest`, `Producto`, `APIError` (estructuras de datos).

**Rutas** (`backend/routes/`):
- `routes.go` — Registro de endpoints: `/api/health`, `/api/login`, `/api/productos`, `/api/productos/:id`.

**Configuración:**
- `go.mod`, `go.sum` — Dependencias Go.
- `docker-compose.yml` + `Dockerfile` — servicios `postgres`, `redis` y `backend` (depende de ambos).
- `go build ./...` — Comando de verificación (no hay tests).

## Frontend (React 19 + Vite)

**Entrada:** `frontend/src/main.tsx` — Montaje de `App.tsx`.

**Rutas y Layout** (`frontend/src/`):
- `App.tsx` — Rutas React Router v7, guardas (`ProtectedRoute`, `AdminRoute`), layout principal.
- `main.tsx` — Punto de entrada.

**Contextos Globales** (`frontend/src/context/`):
- `AuthContext.tsx` — Estado de autenticación (`isAuthenticated`, `user: {email, rol}`), `login()`, `logout()`.
- `CartContext.tsx` — Carrito persistente por usuario (`addToCart`, `removeFromCart`, `totalItems`, `totalPrice`).

**Componentes** (`frontend/src/components/`):
- `Login/` — Login con servicio `authService.ts`, redirección por rol.
- `Layout.tsx` — Layout principal con Navbar + Sidebar.
- `Navbar.tsx` — Navegación superior, carrito, avatar, logout.
- `Sidebar.tsx` — Navegación lateral colapsable, filtros por rol.
- `Dashboard.tsx` — Dashboard admin con KPIs dinámicos.
- `MiRed/` — Red multinivel de referidos (`data/red.ts` con funciones puras).
- `Carrito.tsx` — Lista de ítems, total, checkout.
- `Catalogo/` — Catálogo con filtros, detalle, lightbox, servicios, datos mock.
- `Storefront.tsx` — Storefront full-screen con hero + secciones.
- `Checkout.tsx`, `Confirmacion.tsx` — Flujo de compra.

**Catalogo Feature** (`frontend/src/components/Catalogo/`):
- `Catalogo.tsx` — Grilla de productos con filtros por categoría/búsqueda.
- `interfaces/iproducto.ts` — `IProducto` (id, nombre, descripcion, precio, categoria, img, galeria).
- `data/productos.ts` — Datos mock de 8 productos con categorías.
- `services/productosService.ts` — `getProductos()`, `getProductoById(id)` (simulan API).

**MiRed Feature** (`frontend/src/components/MiRed/`):
- `MiRed.tsx` — Visualización de árbol de referidos.
- `data/red.ts` — Estructura `Referido` (id, nombre, nivel, ventas, hijos), funciones puras de cálculo.

**Configuración:**
- `package.json` — Dependencias: React 19, Vite, Tailwind 4, React Router v7.
- `tsconfig.json`, `vite.config.ts` — Configuración TypeScript y Vite.
- `npm run dev` — Servidor desarrollo; `npm run build` — `tsc -b && vite build`; `npm run lint` — ESLint.

## Specs y Documentación

**Specs atómicas** (`specs/`):
- `paso-01-rol-api/` a `paso-14-navegacion-rol/` — Cada paso tiene `spec.md` (qué) y `tasks.md` (cómo).
- `README.md` — Índice de los 14 pasos con correspondencia a entregables F1-F8.

**Documentos de referencia:**
- `practica03.md` — Guía de laboratorio (código completo, pasos 1-14).
- `AGENTS.md` — Convenciones del proyecto (código en español, estructura, comandos).
- `README.md` — Descripción general del proyecto.

## Patrones de Organización

- **Backend:** Capas clásicas (models → controllers → routes). Hardcodeado sin persistencia aún.
- **Frontend:** Organización por features (`Catalogo/`, `MiRed/`, `Login/`) con subcarpetas `interfaces/`, `data/`, `services/`.
- **Contextos:** Estado global en `context/` (Auth, Cart).
- **Rutas:** Guardas por rol en `App.tsx` (`ProtectedRoute`, `AdminRoute`).

## Comandos de Verificación

**Backend:**
```bash
cd backend && go build ./...
```

**Frontend:**
```bash
cd frontend && npm run lint && npm run build
```
