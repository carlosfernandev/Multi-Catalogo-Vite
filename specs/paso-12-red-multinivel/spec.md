# Paso 12 — Red multinivel visual

**Entregable:** F6
**Correspondencia con la planificación:** AA 2
**Fuente:** `practica03.md` → Paso 12

## Objetivo

Visualizar la red de referidos como un árbol jerárquico con KPIs derivados. El admin ve su red completa, comisiones por nivel y nivel alcanzado.

## Concepto

Un componente recursivo renderiza el árbol: cada `NodoReferido` se dibuja a sí mismo y, si tiene hijos, los dibuja debajo con indentación y borde lateral. Las comisiones se calculan con las funciones puras del Paso 3. La prop se llama `nodo` porque `ref` es reservada en React 19.

## Archivos

| Acción | Ruta |
|---|---|
| reemplazar | `frontend/src/components/MiRed.tsx` |

## Alcance técnico

- `NIVEL_ESTILO` por nivel (badge y border); `NodoReferido({ nodo })` recursivo con `<ul ml-6 border-l-2>`.
- KPIs: `contarRed`, `sumarVentasRed`, `sumarComisiones`, `nivelAlcanzado`; raíz "Tú" con comisiones destacadas.
- Imports desde `./MiRed/data/red.ts` (funciones puras compartidas).
- Ver código en `practica03.md`, Paso 12.

## Criterios de aceptación

- [x] Resumen muestra 7 referidos, $4 090 ventas, $319.50 comisiones, nivel Plata.
- [x] Árbol renderiza 3 niveles con indentación y badge de comisión por nodo.
- [x] Cada nodo muestra `+ $comision` calculada por `comisionDeReferido`.
- [x] Solo admin accede a `/mi-red` (guardas del Paso 6).

## Verificación

Entrar a `/mi-red` como admin y verificar el resumen (7 referidos, $4,090 ventas, $319.50 comisiones, nivel Plata) y el árbol con los 3 niveles.
