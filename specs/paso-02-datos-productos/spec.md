# Paso 02 — Datos de ejemplo del catálogo

**Entregable:** base de F1, F3, F4
**Correspondencia con la planificación:** ACD 1
**Fuente:** `practica03.md` → Paso 2

## Objetivo

Enriquecer el modelo de producto y proveer datos mock que permitan construir el storefront, los filtros y la galería antes de tener PostgreSQL. Sin estos campos el detalle y el lightbox no tendrían qué mostrar.

## Concepto

El detalle con galería y los filtros necesitan más campos que los 4 productos actuales de la API (`descripción`, `categoría`, varias imágenes). Como esta unidad es de frontend, definimos un mock que simula la API; en la Unidad 2 estos datos saldrán de PostgreSQL.

## Archivos

| Acción | Ruta |
|---|---|
| modificar | `frontend/src/components/Catalogo/interfaces/iproducto.ts` |
| crear | `frontend/src/components/Catalogo/data/productos.ts` |

## Alcance técnico

- `IProducto` agrega `descripcion`, `categoria`, `img` (portada) y `galeria: string[]`.
- `data/productos.ts` exporta `categorias` (5 valores) y `productosMock` (8 productos con picsum).
- Import de `IProducto` con `type` en el mock; patrón por feature `interfaces/` + `data/`.
- Ver código en `practica03.md`, Paso 2.

## Criterios de aceptación

- [x] `IProducto` contiene los 7 campos (`id`, `nombre`, `descripcion`, `precio`, `categoria`, `img`, `galeria`) con tipos correctos.
- [x] `productosMock` tiene 8 productos y cada uno tiene `galeria` de 3 imágenes.
- [x] `categorias` coincide con las categorías usadas en los productos.
- [x] `npm run build` y `npm run lint` sin errores de tipos.

## Verificación

```bash
cd frontend && npm run build
# Verificar manualmente que productosMock se importa sin errores de tipo.
```
