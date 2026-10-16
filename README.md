# Multi-Catálogo Vite

Proyecto didáctico de comercio electrónico y ventas multinivel (MLM) para la asignatura **Framework Programación Web** (UPSE, 2026-02).
Combina una SPA en React con un núcleo transaccional en Go/Fiber, persistencia en PostgreSQL y caché en Redis.

---

## 🏛️ Arquitectura

```
React SPA (:5173, proxy /api) ──fetch──► Go/Fiber API (:3000, PORT) ──pgx────► PostgreSQL 16 (:5432, Docker)
                                               │
                                               └──go-redis──► Redis (:6379, Docker) — Cache-Aside
```

- **Frontend (`frontend/`):** React 19, TypeScript, Vite 8, Tailwind CSS 4, React Router 7.
- **Backend (`backend/`):** Go 1.26, Fiber v2, pgx v5, Redis (go-redis v9), JWT (HS256) y bcrypt.
- **Infraestructura (`backend/docker-compose.yml`):** PostgreSQL 16 y Redis Alpine administrados con perfiles (`infra`).

---

## 🚀 Guía de Arranque Rápido

Para ejecutar la solución completa en tu entorno local, sigue estos 3 pasos en orden:

### 1. Levantar la Infraestructura (PostgreSQL + Redis)

La configuración de Docker Compose utiliza el perfil `infra` por defecto para levantar únicamente la base de datos y la caché, dejando libre el puerto 3000 para ejecutar el backend en local:

```bash
cd backend
docker compose up -d
# O explícitamente: docker compose --profile infra up -d
```

> **Verificación:** Puedes confirmar que los contenedores están listos con `docker ps`. Verás `multicatalogo_db` (:5432) y `multicatalogo_redis` (:6379) en estado `healthy`.

### 2. Iniciar el Backend (Go / Fiber)

En una terminal separada, inicia el servidor Go:

```bash
cd backend
go run .
```

- El servidor iniciará en `http://localhost:3000`.
- Puedes configurar el puerto con la variable de entorno `PORT` si necesitas otro puerto: `PORT=3001 go run .`.
- Endpoint de verificación: `GET http://localhost:3000/api/health` retorna `{"status":"ok","database":"up","cache":"up"}`.

### 3. Iniciar el Frontend (React + Vite)

En otra terminal separada, instala las dependencias y arranca el servidor de desarrollo:

```bash
cd frontend
npm install
npm run dev
```

- Vite iniciará en `http://localhost:5173`.
- Las peticiones `/api/*` se redirigen automáticamente al backend local mediante el proxy configurado en `vite.config.ts`.

---

## 🔐 Cuentas y Credenciales de Prueba

| Rol | Email | Password | Acceso / Funciones |
|---|---|---|---|
| **Administrador** | `admin@upse.edu.ec` | `123456` | Acceso a Dashboard (`/`), Red multinivel (`/mi-red`), Catálogo y Carrito |
| **Cliente** | `cliente@upse.edu.ec` | `123456` | Acceso a Tienda (`/tienda`), Catálogo (`/catalogo`), Carrito y Checkout |

> Las contraseñas están almacenadas con hash **bcrypt** en PostgreSQL y la autenticación emite tokens **JWT reales (HS256)** con expiración de 2 horas.

---

## 🛠️ Comandos de Verificación

```bash
# Backend (dentro de backend/)
go vet ./...
go build ./...

# Frontend (dentro de frontend/)
npm run lint
npm run build
```

---

## 🛑 Detener los Servicios

```bash
# 1. Detener procesos locales (Backend y Frontend)
# Presiona Ctrl + C en las respectivas terminales

# 2. Detener contenedores de Docker
cd backend
docker compose down

# (Opcional) Detener y eliminar volúmenes de datos persistentes
docker compose down -v
```
