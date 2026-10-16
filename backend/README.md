# Backend — MultiCatálogo API

API Go/Fiber de alta transaccionalidad con persistencia en PostgreSQL y caché en Redis para el proyecto MultiCatálogo.

---

## 🏛️ Arquitectura

- **Framework:** Fiber v2.52.15
- **Base de datos:** PostgreSQL 16 (en Docker, puerto 5432)
- **Driver:** pgx/v5 (pool de conexiones concurrente)
- **Caché:** Redis (`redis:alpine`, puerto 6379) con go-redis/v9 (patrón Cache-Aside)
- **Seguridad:** Hashing con bcrypt (`golang.org/x/crypto/bcrypt`) y tokens JWT con HMAC-SHA256 (`github.com/golang-jwt/jwt/v5`)
- **Concurrencia:** Goroutines asíncronas para cálculo y distribución de comisiones MLM en segundo plano
- **Transaccionalidad:** Propiedades ACID con bloqueo pesimista `SELECT ... FOR UPDATE`

---

## 📁 Estructura del Módulo

```
backend/
├── cache/
│   └── cache.go           # CacheRepository (Redis): Obtener, Guardar, Eliminar, EliminarPorPatron
├── config/
│   ├── db.go              # Pool de conexiones pgxpool
│   └── redis.go           # Cliente Redis con manejo de degradación
├── controllers/
│   ├── controllers.go     # Agregador Controllers + NewControllers (Inyección de dependencias)
│   ├── authController.go  # POST /api/login (verificación bcrypt + emisión JWT)
│   ├── prodController.go  # GET /api/productos, /api/productos/:id
│   ├── redController.go   # GET /api/red (árbol multinivel)
│   ├── checkoutController.go # POST /api/checkout (ACID + identidad por token + Goroutine)
│   └── healthController.go # GET /api/health (ping a PostgreSQL y Redis)
├── middleware/
│   └── auth.go            # RequireAuth (validador de token Bearer JWT) y RequireAdmin
├── models/
│   └── models.go          # Usuario, Producto, Referido, Pedido, Checkout, APIError
├── repository/
│   └── repository.go      # DBRepository / PostgresRepository (SQL puro + FOR UPDATE)
├── routes/
│   └── routes.go          # Registro de endpoints y protección con middlewares
├── database/
│   ├── init.sql           # Esquema base y seed de datos iniciales
│   ├── 002_checkout.sql   # Migración de tablas pedidos, pedido_items y stock
│   └── 003_seguridad.sql  # Migración pgcrypto, índices de rendimiento y bcrypt
├── scripts/
│   └── init-db.sh         # Script manual de inicialización de base de datos
├── docker-compose.yml     # Orquestación de infraestructura con perfiles
├── .env                  # Variables de entorno
├── go.mod, go.sum        # Dependencias del módulo Go
└── main.go               # Entry point, ensamblado de DI y configuración CORS
```

---

## 🚀 Arranque

### Opción 1: Infraestructura en Docker + Backend local (Recomendado)

Docker Compose está configurado con perfiles para levantar únicamente la infraestructura (`postgres` y `redis`), dejando el puerto `3000` libre para ejecutar la API en local:

```bash
cd backend

# 1. Levantar PostgreSQL y Redis
docker compose up -d

# 2. Iniciar el servidor Go
go run .
```

El servidor iniciará en `http://localhost:3000`.

### Opción 2: Toda la pila contenerizada

Si deseas levantar también el backend dentro de un contenedor Docker:

```bash
docker compose --profile all up -d --build
```

---

## ⚙️ Variables de Entorno (`.env`)

Crea o edita el archivo `.env` en la raíz de `backend/`:

```env
# Base de datos PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=multicatalogo_user
DB_PASSWORD=multicatalogo_pass
DB_NAME=multicatalogo

# Caché Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Servidor y Seguridad
PORT=3000
JWT_SECRET=multicatalogo-super-secret-jwt-key
APP_ENV=dev

# Perfil de Docker Compose (por defecto levanta solo PostgreSQL y Redis)
COMPOSE_PROFILES=infra
```

---

## ⚡ Caché (Redis · Cache-Aside)

| Endpoint | Clave en Redis | TTL | Invalidación |
|---|---|---|---|
| `GET /api/productos` | `productos:all` | 5 min | Tras compra en `/api/checkout` |
| `GET /api/red?email=...` | `red:tree:{email}` | 10 min | Tras compra en `/api/checkout` (`red:tree:*`) |

- Si Redis cae o no responde, la API degrada suavemente y consulta directamente PostgreSQL.
- El endpoint `/api/health` reportará `"status": "degraded"` en caso de fallo en Redis, manteniendo la API operativa.

---

## 🛣️ Endpoints de la API

| Método | Endpoint | Autenticación | Descripción |
|---|---|---|---|
| `GET` | `/api/health` | Pública | Estado del servidor, ping a PostgreSQL y Redis (`200 OK` / `503 Error`) |
| `POST` | `/api/login` | Pública | Login con validación bcrypt; retorna token JWT real firmado |
| `GET` | `/api/productos` | Pública | Catálogo completo (Cache-Aside `productos:all`) |
| `GET` | `/api/productos/:id` | Pública | Detalle de un producto por su ID |
| `GET` | `/api/red` | `Bearer <JWT>` | Árbol jerárquico de la red multinivel (Cache-Aside `red:tree:{email}`) |
| `POST` | `/api/checkout` | `Bearer <JWT>` | Procesamiento de compra transaccional ACID con bloqueo pesimista |

---

## 🛒 Checkout Transaccional (U2-T4)

`POST /api/checkout` requiere cabecera `Authorization: Bearer <token>`:

```json
{
  "usuario_email": "cliente@upse.edu.ec",
  "cliente": {
    "nombre": "Diana Paz",
    "email": "cliente@upse.edu.ec",
    "direccion": "Av. Principal 123",
    "ciudad": "Santa Elena",
    "telefono": "0991234567"
  },
  "metodo_pago": "Tarjeta de crédito",
  "items": [
    { "producto_id": 1, "cantidad": 2 }
  ]
}
```

### Garantías Transaccionales
1. **Identidad validada por Token:** El comprador se extrae de los claims verificados del JWT (`c.Locals`), impidiendo la suplantación de cuenta.
2. **Bloqueo Pesimista (`FOR UPDATE`):** Se bloquea la fila del producto en PostgreSQL durante la transacción para prevenir condiciones de carrera de stock.
3. **Validación de cantidades:** Cantidades $\le 0$ son rechazadas con error `400` (`ErrCantidadInvalida`).
4. **Cálculo de Total en Servidor:** El total se calcula sumando el detalle con los precios oficiales de la base de datos.
5. **Atomicidad:** Si cualquier paso falla, se ejecuta `ROLLBACK` total. Tras el `COMMIT` exitoso, responde `201 Created`.
6. **Goroutine Asíncrona:** Una Goroutine independiente procesa las comisiones de la red multinivel (10%, 5%, 2%) y elimina las claves de caché en Redis (`productos:all` y `red:tree:*`).

---

## 👥 Cuentas de Prueba

| Rol | Email | Password | Acceso / Capacidades |
|---|---|---|---|
| **Admin** | `admin@upse.edu.ec` | `123456` | Acceso a métricas de red, árbol MLM, catálogo y panel administrativo |
| **Cliente** | `cliente@upse.edu.ec` | `123456` | Compra en tienda, catálogo, carrito y checkout transaccional |

> Las contraseñas se almacenan mediante hash **bcrypt (`$2a$`)** generado por la migración `003_seguridad.sql`.

---

## 🧪 Comprobación y Verificación

```bash
# Análisis estático y compilación
go vet ./...
go build ./...

# Probar Health Check
curl -s http://localhost:3000/api/health

# Probar Login y obtener JWT
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Probar endpoint protegido con el token obtenido
curl -s http://localhost:3000/api/red -H "Authorization: Bearer $TOKEN"
```

---

## 🛑 Detener Servicios

```bash
# Detener contenedores de base de datos y caché
docker compose down

# Detener contenedores y eliminar volúmenes persistentes
docker compose down -v
```
