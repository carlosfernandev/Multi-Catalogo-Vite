# Paso 08 — Catálogo con búsqueda y filtros

**Entregable:** F3
**Correspondencia con la planificación:** APD
**Fuente:** `practica03.md` → Paso 8

## Objetivo

Convertir el catálogo en una vista interactiva con búsqueda por texto y filtro por categoría. La URL debe reflejar el filtro para ser compartible y el filtrado debe ser eficiente.

## Concepto

`useMemo` recalcula el filtrado solo cuando cambian `productos`, `busqueda` o `categoria`. Se usa `useSearchParams` para leer el filtro enviado desde el Storefront (`?categoria=`) y mantener la URL compartible.

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/components/Catalogo/Catalogo.tsx` |

## Alcance técnico

- `useSearchParams` para leer/escribir `?categoria`; estado `busqueda` + `categoria`.
- `useMemo` filtra por `categoria === "Todas"` y `nombre/descripcion includes texto`.
- `getProductos` en `useEffect`; estados `loading` y vacío con botón "Limpiar filtros".
- Cards con `Link to=/producto/:id` y `addToCart`.
- Ver código en `practica03.md`, Paso 8.

## Criterios de aceptación

- [ ] Búsqueda filtra por nombre o descripción (case-insensitive).
- [ ] Selector de categoría sincroniza `?categoria=` en la URL.
- [ ] `useMemo` lista dependencias `[productos, busqueda, categoria]`.
- [ ] Estado vacío muestra mensaje y botón de limpiar.

## Verificación

Buscar "serum", filtrar por categoría y verificar que la URL refleja `?categoria=...`; probar el estado vacío con una búsqueda sin resultados.
