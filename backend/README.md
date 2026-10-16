# Backend — MultiCatálogo API

API Go/Fiber con persistencia en PostgreSQL para el proyecto MultiCatálogo.

## Arquitectura

- **Framework:** Fiber v2.52.15
- **Base de datos:** PostgreSQL 16 (en Docker)
- **Driver:** pgx/v5 (pool de conexiones)
- **Configuración:** godotenv (variables de entorno)
- **Caché:** Redis (`redis:alpine`) con go-redis/v9, patrón Cache-Aside

## Estructura

```
backend/
├── cache/
│   └── cache.go           # CacheRepository (Redis): Obtener, Guardar, Eliminar
├── config/
│   ├── db.go              # Pool de conexiones PostgreSQL
│   └── redis.go           # Cliente Redis
├── controllers/
│   ├── authController.go  # POST /api/login
│   ├── prodController.go  # GET /api/productos, /api/productos/:id
│   ├── redController.go   # GET /api/red (red multinivel)
│   ├── checkoutController.go # POST /api/checkout (ACID + Goroutine)
│   └── healthController.go # GET /api/health
├── models/
│   └── models.go          # Usuario, Producto, Referido, APIError
├── repository/
│   └── repository.go      # DBRepository / PostgresRepository
├── routes/
│   └── routes.go          # Registro de endpoints
├── database/
│   └── init.sql           # Script de inicialización de DB
├── scripts/
│   └── init-db.sh         # Script manual de inicialización
├── docker-compose.yml     # Configuración PostgreSQL
├── .env                  # Variables de entorno (no commiteado)
├── go.mod, go.sum        # Dependencias Go
└── main.go               # Punto de entrada
```

## Arranque

### Opción rápida: todo en Docker

```bash
cd backend
docker-compose up -d --build
```

Levanta `postgres` (5432), `redis` (6379) y `backend` (3000). En un volumen nuevo, `database/init.sql` se ejecuta automáticamente.

### 1. Levantar PostgreSQL y Redis (Docker)

```bash
cd backend
docker-compose up -d postgres redis
```

Esto levanta PostgreSQL 16 en el puerto 5432 con:
- Usuario: `multicatalogo_user`
- Contraseña: `multicatalogo_pass`
- Base de datos: `multicatalogo`

### 2. Inicializar la base de datos

El script `init.sql` crea las tablas y pobla los datos semilla. Ejecutar manualmente:

```bash
./scripts/init-db.sh
```

Este script:
- Espera a que PostgreSQL esté listo
- Ejecuta `database/init.sql` (tablas: usuarios, productos, referidos)
- Inserta datos de prueba (2 usuarios, 8 productos, 7 referidos)

### 3. Iniciar el servidor Go

```bash
go run .
```

El servidor escuchará en `http://localhost:3000`.

## Variables de entorno

Archivo `.env` (crear si no existe):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=multicatalogo_user
DB_PASSWORD=multicatalogo_pass
DB_NAME=multicatalogo
REDIS_HOST=localhost
REDIS_PORT=6379
# Alternativa: REDIS_URL=redis://localhost:6379/0
```

## Caché (Redis, Cache-Aside)

| Endpoint | Clave | TTL |
|----------|-------|-----|
| `GET /api/productos` | `productos:all` | 5 min |
| `GET /api/red?email=...` | `red:tree:{email}` (sin `email` → `anonimo`) | 10 min |

La cabecera `X-Cache` indica `MISS` (consultó PostgreSQL) o `HIT` (respondió desde Redis). Si Redis no está disponible, la API sigue funcionando contra PostgreSQL.

```bash
curl -i http://localhost:3000/api/productos | grep X-Cache   # MISS, luego HIT
docker exec multicatalogo_redis redis-cli KEYS '*'
docker exec multicatalogo_redis redis-cli TTL productos:all
docker exec multicatalogo_redis redis-cli FLUSHALL           # invalidar caché
```

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/login` | Login (devuelve email, rol, token) |
| GET | `/api/productos` | Lista completa de productos |
| GET | `/api/productos/:id` | Producto por ID |
| GET | `/api/red` | Red multinivel de referidos (árbol jerárquico) |
| POST | `/api/checkout` | Registra un pedido (transacción ACID) y calcula comisiones en segundo plano |

## Checkout (Fase 4)

Los controladores son `struct` que reciben sus dependencias (`DBRepository`, `CacheRepository`) desde `main.go` → `routes.SetupRoutes`.

`POST /api/checkout`:

```json
{
  "usuario_email": "cliente@upse.edu.ec",
  "cliente": { "nombre": "Cliente Demo", "email": "cliente@upse.edu.ec", "direccion": "Av. 1", "ciudad": "La Libertad", "telefono": "099" },
  "metodo_pago": "Tarjeta de crédito",
  "items": [{ "producto_id": 1, "cantidad": 2 }]
}
```

1. **Transacción:** `BEGIN` → `INSERT pedidos` → `INSERT pedido_items` (precio real de la DB) + `UPDATE productos SET stock = stock - n` → total → `COMMIT`. Cualquier error ejecuta `ROLLBACK` (400 si son datos inválidos, 500 si es un fallo interno). Respuesta `201` con `numero`, `total`, `items`.
2. **Goroutine `calcularComisiones`:** tras el `COMMIT`, suma la venta al nodo del comprador (`referidos.usuario_id`) y acredita 10 % / 5 % / 2 % a los referidos de nivel 1 / 2 / 3 de su cadena ascendente (tabla `comisiones`, columna `saldo_comision`). El cliente recibe el 201 sin esperar este cálculo.
3. **Invalidación de caché:** al aplicar las comisiones se borran `red:tree:{email}` y `red:tree:*` (el árbol es compartido), así el Dashboard obtiene los KPIs nuevos.

Simulación de fallas para demostrar el `ROLLBACK`:

```bash
# Falla tras insertar el pedido y antes del detalle (HTTP 500, no queda el pedido)
CHECKOUT_FALLA_SIMULADA=true docker compose up -d backend
# Pausa de N segundos antes del COMMIT para detener PostgreSQL a mitad
CHECKOUT_PAUSA_COMMIT=10 docker compose up -d backend   # luego: docker stop multicatalogo_db
# Volver al modo normal
docker compose up -d backend
```

Revisar en PostgreSQL:

```bash
docker exec -it multicatalogo_db psql -U multicatalogo_user -d multicatalogo \
  -c "SELECT count(*) FROM pedidos" -c "SELECT id, stock FROM productos" -c "SELECT * FROM comisiones"
```

> En una base existente, aplicar el esquema nuevo con `./scripts/init-db.sh` (sentencias idempotentes).

## Cuentas de prueba

| Rol | Email | Password |
|-----|-------|----------|
| Admin | `admin@upse.edu.ec` | `123456` |
| Cliente | `cliente@upse.edu.ec` | `123456` |

## Verificación

```bash
# Build
go build ./...

# Prueba de login
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}'

# Prueba de productos
curl http://localhost:3000/api/productos

# Prueba de red
curl http://localhost:3000/api/red
```

## Detener servicios

```bash
# Detener servidor Go (Ctrl+C)

# Detener PostgreSQL
docker-compose down

# Detener PostgreSQL y borrar datos (volumen)
docker-compose down -v
```

## Notas

- El password está en texto plano (temporal, para producción usar bcrypt)
- El token JWT es ficticio (temporal, para producción usar JWT real)
- Las ventas del usuario raíz están hardcodeadas en `ConstruirArbolRed`
- El árbol de referidos se construye en memoria desde la lista plana de la DB
