# Paso 14 — Navegación según el rol

**Entregable:** F8
**Correspondencia con la planificación:** ACD 2
**Fuente:** `practica03.md` → Paso 14

## Objetivo

Filtrar la navegación por rol y reflejar el perfil en Sidebar y Navbar. El cliente no debe ver ni acceder a opciones de administración.

## Concepto

El Sidebar filtra las opciones según el rol y resalta la activa con `useLocation`. El Navbar muestra correo, insignia de rol y título acorde al perfil ("Panel de Administración" vs "Tienda MultiCatálogo").

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/components/Sidebar.tsx` |
| modificar | `frontend/src/components/Navbar.tsx` |

## Alcance técnico

- `Sidebar`: `NAV_ITEMS` con `soloAdmin` (Dashboard, Mi Red); `useAuth` + `useLocation`; `filter` por `user?.rol`; `esActivo` con `pathname.startsWith`; colapso con `collapsed`.
- `Navbar`: `user?.email` + badge `rol` (amber admin / indigo cliente); título por rol; `logout` + `navigate("/login")`.
- Ver código en `practica03.md`, Paso 14.

## Criterios de aceptación

- [x] Admin ve Dashboard, Tienda, Catálogo, Mi Red; cliente solo Tienda y Catálogo.
- [x] Opción activa resaltada con `bg-indigo-600`.
- [x] Navbar muestra badge de rol y título por perfil.
- [x] Cerrar sesión navega a `/login`.

## Verificación

Admin ve Dashboard, Tienda, Catálogo y Mi Red; cliente ve solo Tienda y Catálogo. El título del Navbar y la insignia de rol cambian según el perfil.
