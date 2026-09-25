# Paso 05 — Autenticación con rol

**Entregable:** F8
**Correspondencia con la planificación:** ACD 2
**Fuente:** `practica03.md` → Paso 5

## Objetivo

Extender el estado global de autenticación para guardar `email + rol` y aislar el `POST /api/login` en un servicio por feature. El login normaliza el rol y redirige: admin → `/`, cliente → `/tienda`.

## Concepto

El estado global guarda el usuario completo (`email` + `rol`). El login normaliza el rol recibido de la API y redirige según el perfil: admin → Dashboard, cliente → Tienda.

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/context/AuthContext.tsx` |
| crear | `frontend/src/components/Login/services/authService.ts` |
| reemplazar | `frontend/src/components/Login/Login.tsx` |

## Alcance técnico

- `AuthContext`: `Rol = 'admin' | 'cliente'`, `Usuario { email, rol }`, `login(usuario: Usuario)`.
- `authService.ts`: `login(email, password): Promise<LoginResponse>` con `fetch /api/login` y normalización `data.rol === "admin" ? "admin" : "cliente"`.
- `Login.tsx` movido a `components/Login/`, usa alias `login as loginRequest` y `navigate` por rol.
- Ver código en `practica03.md`, Paso 5.

## Criterios de aceptación

- [ ] `useAuth` expone `user: Usuario | null` y `isAuthenticated`.
- [ ] `authService.login` tipa `Rol` y lanza error si `!response.ok`.
- [ ] Login como admin redirige a `/`; como cliente a `/tienda`.
- [ ] `npm run build` sin errores de tipos.

## Verificación

Iniciar sesión con cada cuenta → admin cae en `/`, cliente en `/tienda` (si `/tienda` aún no existe, el comodín `*` redirige a `/` hasta el Paso 7).
