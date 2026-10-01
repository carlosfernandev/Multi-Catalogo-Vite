# AGENTS.md — Multi-Catálogo

Proyecto académico (Framework Programación Web). Dos carpetas independientes en la raíz.

## Estructura

- `backend/` — API Go/Fiber (módulo Go: `backend`). Entrada: `backend/main.go`. Rutas en `backend/routes/routes.go`, lógica en `backend/controllers/` (authController, prodController), modelos en `backend/models/models.go`.
- `frontend/` — React 19 + Vite + Tailwind CSS 4 + React Router v7 + TypeScript. Código en `frontend/src/`.
- `practica03.md` — guía de laboratorio activa (Práctica 03). Todas las rutas allí ya reflejan esta estructura (`backend/...`, `frontend/src/...`).

## Arranque

```bash
# Backend (puerto 3000)
cd backend && go run .

# Frontend (puerto 5173)
cd frontend && npm install && npm run dev
```

API: `POST /api/login`, `GET /api/productos` (CORS permite `http://127.0.0.1:5173`).

## Comandos frontend

- `npm run dev` — servidor de desarrollo
- `npm run build` — `tsc -b && vite build`
- `npm run lint` — `eslint .`
- Sin test runner configurado.

## Frontend: ubicaciones

- `src/App.tsx` — rutas y guardas; `src/main.tsx` — entrada.
- `src/components/` — Login, Layout, Navbar, Sidebar, Dashboard, MiRed, Carrito; `Catalogo/Catalogo.tsx` + `Catalogo/interfaces/iproducto.ts`.
- `src/context/` — `AuthContext.tsx` (auth), `CartContext.tsx` (carrito en memoria).

## Convenciones

- Código, comentarios y textos de UI en español (proyecto académico).
- Respuestas al usuario en español.
- No existen tests: la verificación es `npm run lint` + `npm run build` (frontend) y `go build ./...` (backend).