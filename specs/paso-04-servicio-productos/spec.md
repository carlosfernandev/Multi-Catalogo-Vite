# Paso 04 — Capa de servicios del catálogo

**Entregable:** base de F1, F3
**Correspondencia con la planificación:** ACD 1
**Fuente:** `practica03.md` → Paso 4

## Objetivo

Aislar la obtención de datos del catálogo en una capa de servicios. Los componentes consumen `Promise` igual que con `fetch` real, por lo que en la Unidad 2 solo se cambiará el cuerpo de estas funciones hacia `GET /api/productos`.

## Concepto

Separar la obtención de datos de los componentes. Las funciones devuelven `Promise` (como haría `fetch`), por lo que los componentes se escriben igual que contra una API real.

## Archivos

| Acción | Ruta |
|---|---|
| crear | `frontend/src/components/Catalogo/services/productosService.ts` |

## Alcance técnico

- `getProductos(): Promise<IProducto[]>` con `setTimeout 400 ms` sobre `productosMock`.
- `getProductoById(id): Promise<IProducto | undefined>` con `setTimeout 300 ms` + `find`.
- Imports: `productosMock` y `IProducto` por feature; sin `fetch` aún.
- Ver código en `practica03.md`, Paso 4.

## Criterios de aceptación

- [x] Ambas funciones retornan `Promise` y simulan latencia.
- [x] `getProductos` resuelve con la lista completa (8 productos).
- [x] `getProductoById` resuelve con el producto o `undefined` si no existe.
- [x] Import paths usan `../data/productos` y `../interfaces/iproducto`.

## Verificación

```bash
cd frontend && npm run build
# Importar getProductos en un componente y verificar que resuelve tras ~400 ms.
```
