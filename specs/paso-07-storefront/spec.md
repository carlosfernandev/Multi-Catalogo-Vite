# Paso 07 — Storefront a pantalla completa

**Entregable:** F1
**Correspondencia con la planificación:** ACD 1
**Fuente:** `practica03.md` → Paso 7

## Objetivo

Construir la tienda pública inmersiva (`/tienda`): hero full-viewport, categorías, destacados y explicación del plan multinivel. Es la primera vista que ve un cliente tras el login.

## Concepto

Interfaces full-screen: el hero usa `min-h-screen` con degradados y el resto se compone de secciones. El CTA "Conocer el Plan Multinivel" solo se renderiza para `admin`; el botón final se etiqueta según su destino real.

## Archivos

| Acción | Ruta |
|---|---|
| crear | `frontend/src/components/Storefront.tsx` |

## Alcance técnico

- `useEffect` + `getProductos().slice(0,4)` para destacados; estado `loading`.
- Hero `min-h-screen` con gradiente y círculos `blur-3xl`; categorías mapeadas desde `categorias` con `Link ?categoria=`.
- Condicional `user?.rol === "admin"` para el CTA a `/mi-red`; CTA final a `/catalogo`.
- Ver código en `practica03.md`, Paso 7.

## Criterios de aceptación

- [ ] Hero ocupa el viewport completo (`min-h-screen`) con CTA a `/catalogo`.
- [ ] Sección categorías muestra 5 cards enlazadas con `?categoria=`.
- [ ] Destacados muestran 4 productos con imagen, categoría y precio.
- [ ] Cliente no ve "Conocer el Plan Multinivel"; admin sí.

## Verificación

Navegar a `/tienda` y verificar que el hero ocupa toda la pantalla; con rol cliente el botón del plan multinivel no aparece; con rol admin sí.
