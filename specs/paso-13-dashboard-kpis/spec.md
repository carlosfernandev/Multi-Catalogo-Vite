# Paso 13 — Dashboard con KPIs dinámicos

**Entregable:** F7
**Correspondencia con la planificación:** AA 2
**Fuente:** `practica03.md` → Paso 13

## Objetivo

Reemplazar los 3 valores estáticos del dashboard por KPIs derivados del estado real de la red. El dashboard comparte la misma fuente de verdad que `MiRed`.

## Concepto

Los valores ya no son estáticos; se calculan con las mismas funciones puras de `data/red.ts` (una sola fuente de verdad compartida con `MiRed`). Incluye el top 3 de referidos por ventas y una barra de progreso al siguiente nivel.

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/components/Dashboard.tsx` |

## Alcance técnico

- KPIs: `referidosDirectos`, `contarRed`, `sumarVentasRed`, `sumarComisiones`, `nivelAlcanzado`.
- Top 3: `[...hijos].sort(ventas desc).slice(0,3)` con `comisionDeReferido`.
- Barra de progreso `width = min(referidosDirectos/6*100, 100)%`; leyenda de niveles.
- Ver código en `practica03.md`, Paso 13.

## Criterios de aceptación

- [x] KPIs muestran $4 090 ventas, 7 referidos, $319.50 comisiones, nivel Plata.
- [x] Top 3 ordenado por ventas descendente con comisión.
- [x] Barra de progreso refleja `referidosDirectos / 6`.
- [x] Solo admin ve `/` (guarda `AdminRoute`).

## Verificación

Entrar al Dashboard como admin y verificar los KPIs: $4,090 ventas, 7 referidos, $319.50 comisiones, nivel Plata.
