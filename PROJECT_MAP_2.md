# PROJECT MAP — Multi-Catalogo-Vite

> Mapa arquitectónico del workspace. Leer ESTE archivo antes de explorar código.
> Proyecto didáctico "MultiCatálogo": tienda de cosméticos con red de ventas
> multinivel (MLM). Asignatura: Framework Programación Web (UPSE, 2026-02).

## Visión general

```
React SPA (:5173, proxy /api) ──fetch──► Go/Fiber API (:3000, env PORT) ──pgx────► PostgreSQL 16 (:5432, Docker)
                                               │
                                               └──go-redis──► Redis (:6379, Docker) — caché de lecturas
```

## Documentos raíz (meta, no código)

| Archivo | Contenido |
|---|---|
| `PROJECT_MAP.md` / `PROJECT_MAP_2.md` | Este archivo. Índice de navegación y mapa de arquitectura del proyecto. |
| `REQUIREMENTS_CHECKLIST.md` / `REQUIREMENTS_CHECKLIST_2.md` | Estado del proyecto en **fases por Unidad/Tema del sílabo** (U1 SPA → U2 Go/Fiber → U3 FastAPI SaaS → U4 despliegue) + roadmap post-sílabo. Principios: seguridad por defecto, DI, funciones puras de dominio, mobile-first. |
| `AGENTS.md` | Guía para agentes: estructura `backend/` + `frontend/` y comandos de verificación. |
| `README.md` | Descripción general del proyecto. |
| `practica03.md` | Guía de la Práctica 03 (interfaces full-screen + flujos MLM). |
| `specs/` | Specs y tareas por paso (`paso-01-rol-api` … `paso-14-navegacion-rol`). |
| `odd/tasks/` | Notas de tareas (`conectar-controladores-backend.md`). |
| `.windsurf/rules/` | Reglas del agente IA (`globales.md`, `contexto.md`, `ahorro-tokens.md`). |

## frontend/ — SPA

Stack: React 19 + TypeScript + Vite 8 + Tailwind CSS 4 + React Router 7.
Comandos (en `frontend/`): `npm run dev` (:5173, proxy `/api` → `:3000` en `vite.config.ts`) | `npm run build` | `npm run lint`.

| Ruta (`frontend/src/`) | Responsabilidad |
|---|---|
| `App.tsx` | Router, guardas `ProtectedRoute`/`AdminRoute`, `CartBoundary` (carrito por usuario). |
| `main.tsx` | Bootstrap de React. |
| `index.css` | Tailwind v4 (import único). |
| `context/auth.ts` | Tipos `Usuario`, `Rol`, helper `getStoredToken()`, hook `useAuth()`. |
| `context/AuthContext.tsx` | `AuthProvider` con persistencia JWT en localStorage (`multicatalogo_sesion`). |
| `context/cart.ts` | Tipos `CartItem`, constante `CART_STORAGE_PREFIX`, hook `useCart()`. |
| `context/CartContext.tsx` | `CartProvider` aislado para Fast Refresh de ESLint. |
| `services/http.ts` | Cliente HTTP centralizado (`apiFetch`, `http`), Bearer automático y `ApiError` tipado. |
| `components/Catalogo/interfaces/iproducto.ts` | Tipo `IProducto`. |
| `components/Catalogo/data/productos.ts` | Tipos `Producto` / datos mock + lista `categorias`. |
| `components/Catalogo/services/productosService.ts` | `getProductos()`, `getProductoById()` → `GET /api/productos[/:id]`. |
| `components/MiRed/data/red.ts` | Tipos `Referido` + funciones puras: comisiones (10/5/2 %), totales, nivel. |
| `components/MiRed/services/redService.ts` | `getRed(email?)` → `GET /api/red[?email=]`. |
| `components/Login/Login.tsx` | Login; lee `APIError.message` del backend; redirige por rol. |
| `components/Login/services/authService.ts` | `login(email, password)` → `POST /api/login`. |
| `services/checkoutService.ts` | `crearPedido()` → `POST /api/checkout`; lanza `APIError.message`. |
| `services/healthService.ts` | `checkHealth()` → `GET /api/health`. |
| `components/Storefront.tsx` | `/tienda` — hero full-screen, categorías, destacados, plan MLM. |
| `components/Catalogo/Catalogo.tsx` | `/catalogo` — búsqueda + filtro por `?categoria=` (useMemo, useSearchParams). |
| `components/DetalleProducto.tsx` | `/producto/:id` — detalle + galería/lightbox. |
| `components/Carrito.tsx` | `/carrito` — ítems, cantidades, total. |
| `components/Checkout.tsx` | `/checkout` — formulario → `POST /api/checkout` real → resumen; muestra `APIError`. |
| `components/Confirmacion.tsx` | `/confirmacion` — pedido confirmado, limpia carrito. |
| `components/MiRed.tsx` | `/mi-red` (admin) — árbol multinivel + comisiones. |
| `components/Dashboard.tsx` | `/` (admin) — KPIs derivados de la red. |
| `components/Layout.tsx` `Navbar.tsx` `Sidebar.tsx` | Shell: sidebar colapsable, navbar con carrito/avatar/logout. |

## backend/ — API REST

Stack: Go 1.26 + Fiber v2 + pgx v5 (pool) + godotenv + PostgreSQL 16 + Redis (go-redis v9).
Comandos (en `backend/`): `docker compose up -d` | `go run .` (:3000, configurable con `PORT`) | `go vet ./...` | `go build ./...`.
Capas: `routes` → `controllers` → `repository` / `cache` (dependencias inyectadas desde `main.go`).

| Ruta (`backend/`) | Responsabilidad |
|---|---|
| `main.go` | Entry point: `InitDB()` + `InitRedis()`, ensambla DI (`NewPostgresRepository` + `NewRedisCache` → `SetupRoutes`), CORS (`127.0.0.1:5173`), `app.Listen(PORT‖":3000")`. |
| `config/db.go` | Pool pgx; carga `.env`; `InitDB()`/`Close()`. |
| `config/redis.go` | Cliente go-redis (`InitRedis()`); si Redis cae, la API sigue sin caché. |
| `cache/cache.go` | Interfaz `CacheRepository` + `RedisCache` (Get/Set/Del/EliminarPorPatron nil-safe). |
| `models/models.go` | `LoginRequest`, `Usuario`, `Producto`(+stock), `Referido`, `CheckoutRequest`/`ItemPedido`, `Pedido`/`PedidoItem`, `APIError{status,message,details}`. |
| `routes/routes.go` | `SetupRoutes(app, repo, cache)` — registra endpoints con métodos de structs inyectados. |
| `controllers/controllers.go` | Agregador opcional `Controllers` + `NewControllers(repo, cache)` — punto de ensamblado de la DI. |
| `controllers/authController.go` | `AuthController{repo}` → `POST /api/login` — valida body + campos vacíos (400), 401 con APIError. |
| `controllers/checkoutController.go` | `CheckoutController{repo,cache}` → `POST /api/checkout`: transacción ACID, Goroutine comisiones, invalidación caché. |
| `controllers/prodController.go` | `ProductoController` → `GET /api/productos` (Cache-Aside `productos:all`, TTL 5 min), `GET /api/productos/:id` (400 si id no numérico, 404). |
| `controllers/redController.go` | `RedController` → `GET /api/red[?email=]` — árbol multinivel; Cache-Aside `red:tree:{email}`, TTL 10 min. |
| `controllers/healthController.go` | `HealthController` → `GET /api/health` — `{"status":"ok"}` con verificación de PostgreSQL y Redis. |
| `repository/repository.go` | Interfaz `DBRepository` + `PostgresRepository{db}`: SQL puro; árbol `referidos` por `parent_id`; `CrearPedido` transaccional (`FOR UPDATE`); `AplicarComisiones` (upline 3 niveles). |
| `database/init.sql` | DDL + seed: `usuarios`, `productos` (`galeria TEXT[]`, `stock`), `referidos` (`parent_id` FK), `pedidos`, `pedido_items`. |
| `database/002_checkout.sql` | Migración checkout para volúmenes existentes (`ALTER stock` + tablas pedidos). |
| `scripts/init-db.sh` | Aplica `init.sql` sobre el contenedor existente (`docker exec -i multicatalogo_db psql`). |
| `docker-compose.yml` + `Dockerfile` | `postgres:16` (:5432, volumen `postgres_data`, autoejecuta `init.sql`), `redis:alpine` (:6379) y `backend` (:3000), con healthchecks. |
| `.env` (gitignored) / `.env.example` | `DB_HOST/PORT/USER/PASSWORD/NAME` + `REDIS_HOST/PORT`. Puerto host por defecto **5432** (o configurable con `DB_PORT`). |
| `README.md` | Arranque, variables, caché, checkout y simulación de fallas. |

## Datos y contratos clave

- Errores API: siempre `{"status":N,"message":"...","details"?}` (struct `APIError`).
- Login OK: `{"token","email","rol"}`. Cuentas: `admin@upse.edu.ec`/`cliente@upse.edu.ec` — `123456`.
- Checkout: `POST /api/checkout` `{usuario_email, cliente:{nombre,email,direccion,ciudad,telefono}, metodo_pago, items:[{producto_id,cantidad}]}` → `201 {"pedido_id","numero","total","estado","fecha","items"}`; el total se recalcula en servidor y el stock se descuenta en la misma transacción.
- Árbol de referidos: raíz id 0 "Tú", 7 nodos, comisiones por nivel (front: `frontend/src/components/MiRed/data/red.ts`).
- Comisión esperada de la red semilla: **$319.50** total, 7 referidos.

## Notas operativas

- Puertos: PostgreSQL **5432** (en `backend/docker-compose.yml`, o configurable si hay otro Postgres local), Redis **6379**, API **3000**, Vite **5173** (`frontend/vite.config.ts`).
- Puerto API configurable con `PORT` (defecto 3000): si otro servicio ocupa 3000, usar `PORT=3001 go run .`.
- Caché Cache-Aside: `productos:all` (TTL 5 min) y `red:tree:{email}` (TTL 10 min); si Redis no responde, la API consulta PostgreSQL directamente. El checkout invalida ambas claves (`red:tree:*`).
- Comisiones MLM: Goroutine tras COMMIT suma `total` a `ventas` del nodo del comprador (`referidos.usuario_id`) y acredita 10/5/2 % a su cadena ascendente (máx. 3 niveles por `parent_id`). Con el seed actual, `cliente` → Diana Paz (nivel 3) acredita a Diana, Carlos Ruiz y Ana García.
- Passwords en texto plano (temporal didáctico); token JWT ficticio; endpoints sin auth real (fase U2-T4 introduce bcrypt y JWT real).
- Estructura del repositorio: `frontend/`, `backend/` y los `.md` raíz.
