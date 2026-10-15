# Feature: Conectar todos los controladores del backend con el frontend

## Objetivo
Conectar todos los endpoints y controladores expuestos en el backend Go/Fiber (`/api/login`, `/api/productos`, `/api/productos/:id`, `/api/red`, `/api/health`) con sus respectivos servicios y componentes en el frontend React/Vite.

## Estado actual de los controladores
- `authController.go` (`POST /api/login`): Conectado en `frontend/src/components/Login/services/authService.ts`.
- `prodController.go` (`GET /api/productos`, `GET /api/productos/:id`): Conectado en `frontend/src/components/Catalogo/services/productosService.ts`.
- `redController.go` (`GET /api/red`): Pendiente de conectar. El frontend consume datos estáticos de `MiRed/data/red.ts` en `MiRed.tsx` y `Dashboard.tsx`.
- `healthController.go` (`GET /api/health`): Pendiente de conectar. No existía servicio ni consumo en frontend.

## Tareas

- [x] TASK-1: Crear servicio `redService.ts` (`GET /api/red`) y actualizar `MiRed.tsx` y `Dashboard.tsx` para consumir datos reales de la API. (Ruta: Delegated direct, Trigger: Writer trigger)
- [x] TASK-2: Crear servicio `healthService.ts` (`GET /api/health`) y consumirlo en `Sidebar.tsx` para indicar el estado operativo del backend. (Ruta: Delegated direct, Trigger: Writer trigger)
- [x] TASK-3: Ejecutar suite de verificación (`npm run lint`, `npm run build`, `go build ./...`).

## Verificación
- `frontend`: `npm run lint` ejecutado sin errores (0 errors, 2 warnings preexistentes de react-refresh).
- `frontend`: `npm run build` (`tsc -b && vite build`) completado exitosamente.
- `backend`: `go build ./...` compilado sin errores.

## Evidencia de entrega
- Commit: `fafa7b4` (`feat(frontend): conectar servicios de productos, red y health con la API de Go`)

## Espejo Engram
- Estado: Pendiente (herramienta de memoria engram no disponible en esta sesión).

