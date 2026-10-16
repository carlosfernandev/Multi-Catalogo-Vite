# REQUIREMENTS_CHECKLIST — Checklist de Requerimientos

## Estado General

**Hito actual:** Práctica 03 — Unidad 1, Tema 5 (Implementación de interfaces interactivas a pantalla completa y flujos para catálogos de venta multinivel)

**Última tarea completada:** Fase 4 — Checkout transaccional (DI, `POST /api/checkout` ACID, Goroutine de comisiones, invalidación de caché)

**Siguiente tarea pendiente:** Paso 5 — Autenticación con rol (extender `AuthContext` con `email + rol` y crear `authService.ts`)

---

## Hito: Backend — Base de Autenticación

### Paso 01 — Rol en la API
- [x] `POST /api/login` con `admin@upse.edu.ec / 123456` retorna `rol: "admin"`, `email` y `token`
- [x] `POST /api/login` con `cliente@upse.edu.ec / 123456` retorna `rol: "cliente"`, `email` y `token`
- [x] Credenciales inválidas retornan HTTP 401 con estructura `APIError`
- [x] Validación de campos vacíos (`Email` y `Password`) con HTTP 400
- [x] `go build ./...` compila sin errores

### Salud del API
- [x] `GET /api/health` devuelve `{"status": "ok", "message": "Servidor Go/Fiber operativo"}`

### Productos API
- [x] `GET /api/productos` devuelve lista de productos (4 hardcodeados)
- [x] `GET /api/productos/:id` devuelve producto por ID con validación de numérico
- [x] Estructura `APIError` estandarizada en `models/models.go`

### Fase 3 — Caché con Redis (Cache-Aside)
- [x] Servicio `redis` (`redis:alpine`, puerto 6379) en `docker-compose.yml`; `backend` depende de `postgres` y `redis`
- [x] `config/redis.go` con `redis.NewClient` leyendo `REDIS_URL` / `REDIS_HOST` / `REDIS_PORT`
- [x] `GET /api/productos` cachea `productos:all` con TTL de 5 minutos
- [x] `GET /api/red` cachea `red:tree:{email}` (query `?email=`) con TTL de 10 minutos
- [x] Cache Hit responde sin consultar PostgreSQL; si Redis cae, la API sigue respondiendo desde PostgreSQL

### Fase 4 — Checkout transaccional, concurrencia e inyección de dependencias
- [x] Paso 1: controladores como `struct` con `DBRepository` y `CacheRepository` inyectados desde `main.go` / `routes.go` (sin variables globales `DB`/`RDB`)
- [x] Paso 2: `POST /api/checkout` recibe `usuario_email`, `cliente`, `metodo_pago` e `items` (`producto_id`, `cantidad`)
- [x] Paso 3: `CrearPedido` en transacción (`BEGIN` → `pedidos` → `pedido_items` + descuento de `productos.stock` → `COMMIT`); precio tomado de la DB
- [x] Paso 3: `ROLLBACK` ante producto inexistente (400), stock insuficiente (400), falla simulada (500) y PostgreSQL caído antes del `COMMIT` (500): sin pedido fantasma ni descuento de stock
- [x] Paso 4: Goroutine `calcularComisiones` tras el `COMMIT`: suma la venta al nodo del comprador y acredita 10/5/2 % a la cadena ascendente (tabla `comisiones`, `referidos.saldo_comision`)
- [x] Paso 5: la Goroutine invalida `red:tree:{email}` y `red:tree:*`; el siguiente `GET /api/red` es `MISS` con KPIs actualizados
- [x] Frontend: `Checkout.tsx` usa `services/checkoutService.ts`; error → alerta sin avanzar; éxito → `Confirmacion.tsx` con datos del backend
- [x] Dashboard y Mi Red envían `?email=` a `/api/red`
- [x] Corrección: `ConstruirArbolRed` incluye el nivel 3 de la red

---

## Hito: Frontend — Datos y Servicios

### Paso 02 — Datos de ejemplo del catálogo
- [x] `IProducto` contiene 7 campos: `id`, `nombre`, `descripcion`, `precio`, `categoria`, `img`, `galeria`
- [x] `productosMock` tiene 8 productos con galería de 3 imágenes cada uno
- [x] `categorias` exportada con 5 valores: Serum, Crema, Tónico, Mascarilla, Kit
- [x] `npm run build` y `npm run lint` sin errores de tipos

### Paso 03 — Datos de la red multinivel
- [ ] Interfaz `Referido` con `id`, `nombre`, `nivel`, `ventas`, `hijos?`
- [ ] `redInicial` con estructura jerárquica de 3 niveles (7 referidos)
- [ ] Función `contarRed(raiz)` cuenta referidos totales
- [ ] Función `sumarVentasRed(raiz)` suma ventas de toda la red
- [ ] Función `comisionDeReferido(ref)` calcula comisión por nodo
- [ ] Función `sumarComisiones(raiz)` suma comisiones totales
- [ ] Función `nivelAlcanzado(directos)` devuelve Bronce/Plata/Oro/Diamante
- [ ] `TASA_COMISION` por nivel (10%, 5%, 2%)

### Paso 04 — Servicio de productos
- [ ] `getProductos(): Promise<IProducto[]>` simula GET /api/productos
- [ ] `getProductoById(id): Promise<IProducto | undefined>` simula GET /api/productos/:id
- [ ] Timeout simulado de 400ms y 300ms respectivamente
- [ ] Import desde `productosMock` y tipado correcto

---

## Hito: Frontend — Autenticación y Rutas

### Paso 05 — Autenticación con rol
- [ ] `AuthContext` expone `Rol = 'admin' | 'cliente'`
- [ ] `Usuario` interface con `email` y `rol`
- [ ] `login(usuario: Usuario)` guarda email y rol en estado global
- [ ] `authService.ts` con `login(email, password): Promise<LoginResponse>`
- [ ] Normalización de rol: `data.rol === "admin" ? "admin" : "cliente"`
- [ ] `Login.tsx` redirige: admin → `/`, cliente → `/tienda`
- [ ] `npm run build` sin errores de tipos

### Paso 06 — Rutas y guardas por rol
- [ ] `ProtectedRoute` redirige a `/login` si no autenticado
- [ ] `AdminRoute` redirige a `/tienda` si rol != admin
- [ ] `CartBoundary` con `key={user?.email ?? "anonimo"}`
- [ ] Rutas públicas: `/login`
- [ ] Rutas admin-only: `/`, `/mi-red`
- [ ] Rutas ambos roles: `/tienda`, `/catalogo`, `/producto/:id`, `/carrito`, `/checkout`, `/confirmacion`
- [ ] `npm run build` sin errores

---

## Hito: Frontend — Storefront y Catálogo

### Paso 07 — Storefront a pantalla completa
- [ ] `Storefront.tsx` creado con hero `min-h-screen`
- [ ] Hero con gradiente y círculos `blur-3xl`
- [ ] Sección categorías con 5 cards enlazadas a `?categoria=`
- [ ] Destacados con 4 productos desde `getProductos().slice(0,4)`
- [ ] Condicional `user?.rol === "admin"` para CTA "Conocer el Plan Multinivel"
- [ ] Estado `loading` durante carga de productos
- [ ] CTA final a `/catalogo`

### Paso 08 — Catálogo con búsqueda y filtros
- [ ] `useSearchParams` para leer/escribir `?categoria`
- [ ] Estado `busqueda` y `categoria`
- [ ] `useMemo` filtra por categoría (case-insensitive) y texto en nombre/descripción
- [ ] Dependencias: `[productos, busqueda, categoria]`
- [ ] Selector de categoría sincroniza URL
- [ ] Estado vacío con mensaje y botón "Limpiar filtros"
- [ ] Cards con `Link to=/producto/:id` y botón `addToCart`

### Paso 09 — Detalle y lightbox
- [ ] Ruta dinámica `/producto/:id` con `useParams`
- [ ] `DetalleProducto.tsx` carga producto con `getProductoById(id)`
- [ ] Galería de imágenes con miniaturas
- [ ] Lightbox overlay al hacer clic en imagen principal
- [ ] Botón cerrar lightbox (X o ESC)
- [ ] Accesibilidad: focus trap, ARIA labels
- [ ] Botón "Añadir al Carrito" en detalle

---

## Hito: Frontend — Carrito Persistente

### Paso 10 — Carrito persistente por usuario
- [ ] `CartItem extends IProducto { cantidad }`
- [ ] `CART_STORAGE_PREFIX = "multicatalogo_carrito_"`
- [ ] `storageKey` por `user?.email`
- [ ] `useState` lazy init desde `localStorage`
- [ ] `useEffect` persiste en cada cambio
- [ ] `addToCart` incrementa cantidad si existe
- [ ] `decrementQuantity` elimina si cantidad llega a 0
- [ ] `incrementQuantity` suma 1
- [ ] `clearCart` vacía el carrito
- [ ] `Carrito.tsx` con controles `− / +`, subtotales y link a `/checkout`
- [ ] Recargar página conserva carrito
- [ ] Cambiar de usuario muestra carrito propio (no mezclado)

---

## Hito: Frontend — Checkout y Confirmación

### Paso 11 — Checkout y confirmación
- [ ] `Checkout.tsx` con formulario de datos de envío
- [ ] Validación de campos obligatorios
- [ ] Resumen del pedido (ítems, total)
- [ ] Botón "Confirmar Pedido" llama `clearCart()`
- [ ] Navegación a `/confirmacion`
- [ ] `Confirmacion.tsx` con mensaje de éxito
- [ ] Botón "Volver al Catálogo" redirige a `/catalogo`

---

## Hito: Frontend — Red Multinivel

### Paso 12 — Red multinivel visual
- [ ] `MiRed.tsx` reemplazado con árbol jerárquico
- [ ] `NodoReferido({ nodo })` recursivo
- [ ] `NIVEL_ESTILO` por nivel (badge y border)
- [ ] Indentación y borde lateral `ml-6 border-l-2`
- [ ] KPIs: 7 referidos, $4,090 ventas, $319.50 comisiones, nivel Plata
- [ ] Comisiones por nodo con `+ $comision`
- [ ] Raíz "Tú" con comisiones destacadas
- [ ] Solo admin accede a `/mi-red` (guardas)

---

## Hito: Frontend — Dashboard y Navegación

### Paso 13 — Dashboard con KPIs dinámicos
- [ ] `Dashboard.tsx` con 3 tarjetas de KPIs
- [ ] KPIs derivados desde `red.ts`: ventas totales, comisiones, nivel
- [ ] Gráficos o indicadores visuales
- [ ] Solo admin accede a `/` (guardas)

### Paso 14 — Navegación por rol
- [ ] `Sidebar.tsx` con menú condicional por rol
- [ ] Admin: Dashboard, Mi Red, Catálogo, Carrito
- [ ] Cliente: Tienda, Catálogo, Carrito
- [ ] `Navbar.tsx` con avatar y logout
- [ ] Colapsable en móvil
- [ ] Active route highlighting

---

## Hito: Verificación Final

### Pruebas end-to-end
- [ ] Login admin → acceso a `/` y `/mi-red`
- [ ] Login cliente → acceso a `/tienda` y `/catalogo`, redirección desde `/mi-red`
- [ ] Sin login → redirección a `/login` en rutas protegidas
- [ ] Carrito persiste al recargar
- [ ] Carrito aislado por usuario
- [ ] Catálogo con filtros funciona y URL refleja `?categoria=`
- [ ] Detalle con lightbox funciona
- [ ] Checkout limpia carrito y muestra confirmación
- [ ] Red multinivel muestra KPIs correctos

### Build y Lint
- [ ] `backend: go build ./...` sin errores
- [ ] `frontend: npm run lint` sin errores
- [ ] `frontend: npm run build` sin errores
