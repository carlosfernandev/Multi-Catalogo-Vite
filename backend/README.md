# Backend — MultiCatálogo API

API Go/Fiber con persistencia en PostgreSQL para el proyecto MultiCatálogo.

## Arquitectura

- **Framework:** Fiber v2.52.15
- **Base de datos:** PostgreSQL 16 (en Docker)
- **Driver:** pgx/v5 (pool de conexiones)
- **Configuración:** godotenv (variables de entorno)

## Estructura

```
backend/
├── config/
│   └── db.go              # Pool de conexiones PostgreSQL
├── controllers/
│   ├── authController.go  # POST /api/login
│   ├── prodController.go  # GET /api/productos, /api/productos/:id
│   ├── redController.go   # GET /api/red (red multinivel)
│   └── healthController.go # GET /api/health
├── models/
│   └── models.go          # Usuario, Producto, Referido, APIError
├── repository/
│   └── repository.go      # Funciones de acceso a datos
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

### 1. Levantar PostgreSQL (Docker)

```bash
cd backend
docker-compose up -d
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
```

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/login` | Login (devuelve email, rol, token) |
| GET | `/api/productos` | Lista completa de productos |
| GET | `/api/productos/:id` | Producto por ID |
| GET | `/api/red` | Red multinivel de referidos (árbol jerárquico) |

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
