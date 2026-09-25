# Paso 09 — Detalle de producto con galería y lightbox

**Entregable:** F2, F4
**Correspondencia con la planificación:** ACD 2
**Fuente:** `practica03.md` → Paso 9

## Objetivo

Crear la ruta dinámica `/producto/:id` con galería navegable y lightbox a pantalla completa. Debe manejar carga, "no encontrado" y accesibilidad del overlay (Escape, navegación).

## Concepto

Ruta dinámica con `useParams`; galería con imagen activa y miniaturas; lightbox como overlay fijo (`fixed inset-0`) con navegación anterior/siguiente y cierre con `Escape` (listener en `window` con limpieza en `useEffect`).

## Archivos

| Acción | Ruta |
|---|---|
| crear | `frontend/src/components/DetalleProducto.tsx` |

## Alcance técnico

- `useParams<{id}>` + `getProductoById(Number(id))` con bandera `activo` para evitar setState tras desmontaje.
- Galería `imagenes = [img, ...galeria]`; miniaturas con `border-indigo-600` en activa.
- Lightbox `fixed inset-0 z-50` con `‹ ›`, `×`, contador `n / total` y `keydown Escape`.
- `addToCart` con feedback `agregado` por 1.5 s.
- Ver código en `practica03.md`, Paso 9.

## Criterios de aceptación

- [ ] `/producto/1` muestra datos correctos y 4 imágenes navegables.
- [ ] Lightbox abre al click, navega con `‹ ›` y cierra con `×` o `Escape`.
- [ ] `/producto/999` muestra "Producto no encontrado" con link a `/catalogo`.
- [ ] Botón "Añadir al Carrito" da feedback visual temporal.

## Verificación

Abrir `/producto/1`, cambiar de imagen con miniaturas, abrir lightbox, navegar con ‹ › y cerrar con Escape o ×; probar `/producto/999`.
