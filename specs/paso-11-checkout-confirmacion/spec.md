# Paso 11 — Checkout y confirmación

**Entregable:** F5
**Correspondencia con la planificación:** AA 1
**Fuente:** `practica03.md` → Paso 11

## Objetivo

Cerrar el flujo transaccional: formulario de envío con validación, resumen en vivo y confirmación con número de pedido. Al confirmar se vacía el carrito y se pasa el pedido por `location.state`.

## Concepto

El checkout valida datos de envío y resume el pedido en vivo. Al confirmar simula procesamiento (1,2 s), ejecuta `clearCart` y navega a `/confirmacion` con `location.state`. La confirmación lee ese estado o muestra un mensaje genérico si se accede directo.

## Archivos

| Acción | Ruta |
|---|---|
| crear | `frontend/src/components/Checkout.tsx` |
| crear | `frontend/src/components/Confirmacion.tsx` |

## Alcance técnico

- `Checkout`: `DatosEnvio` (nombre, email, dirección, ciudad, teléfono, metodoPago), `METODOS_PAGO` (3), `clearCart` + `navigate("/confirmacion", { state: { pedido } })` con `numero: MC-xxxxxx`.
- `Confirmacion`: `useLocation` tipado `PedidoConfirmado`; fallback si no hay `location.state.pedido`.
- Carrito vacío en checkout muestra link a `/catalogo`.
- Ver código en `practica03.md`, Paso 11.

## Criterios de aceptación

- [x] Checkout muestra carrito vacío si no hay ítems.
- [x] Formulario valida campos `required` y método de pago.
- [x] Confirmar genera `MC-xxxxxx`, vacía el carrito y muestra detalle con total.
- [x] Acceso directo a `/confirmacion` muestra mensaje genérico con link a `/tienda`.

## Verificación

Flujo completo: catálogo → añadir al carrito → carrito → checkout → llenar formulario → confirmar → ver número de pedido y carrito vacío.
