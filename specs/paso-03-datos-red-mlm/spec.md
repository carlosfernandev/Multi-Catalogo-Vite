# Paso 03 — Datos de ejemplo de la red multinivel

**Entregable:** base de F6, F7
**Correspondencia con la planificación:** AA 2
**Fuente:** `practica03.md` → Paso 3

## Objetivo

Modelar la red de referidos como un árbol jerárquico y proveer funciones puras para comisiones y KPIs. Estas funciones serán la única fuente de verdad que compartirán `MiRed` y `Dashboard`, evitando duplicar cálculos.

## Concepto

La red de referidos es un árbol jerárquico. Definimos la estructura de datos y funciones puras (sin efectos) que calculan comisiones, totales y el nivel alcanzado.

## Archivos

| Acción | Ruta |
|---|---|
| crear | `frontend/src/components/MiRed/data/red.ts` |

## Alcance técnico

- `Referido` con `id`, `nombre`, `nivel` (1–3), `ventas`, `hijos?`.
- `TASA_COMISION` (10 % / 5 % / 2 %), `redInicial` (7 referidos, 3 niveles, ventas $4 090).
- Funciones puras: `contarRed`, `sumarVentasRed`, `comisionDeReferido`, `sumarComisiones`, `nivelAlcanzado`.
- Ver código en `practica03.md`, Paso 3.

## Criterios de aceptación

- [ ] `redInicial` tiene estructura de 3 niveles y 7 referidos (sin contar raíz).
- [ ] `sumarComisiones(redInicial)` retorna 319.50.
- [ ] `contarRed` y `sumarVentasRed` coinciden con los valores esperados (7 / $4 090).
- [ ] `nivelAlcanzado` mapea 0–1→Bronce, 2–3→Plata, 4–5→Oro, 6+→Diamante.

## Verificación

Verificación del cálculo esperado: Ana (nivel 1, $1200) → $120.00; Carlos (nivel 2, $500) → $25.00; Diana (nivel 3, $300) → $6.00. Total red: $319.50 y 7 referidos. Validar con `npm run build`.
