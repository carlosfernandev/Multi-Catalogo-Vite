# Reglas Globales del Agente

## Navegación: NUNCA hagas escaneos recursivos ni leas carpetas completas

**OBLIGATORIO:** Antes de leer cualquier archivo de código, consulta primero `PROJECT_MAP.md` para ubicar la ruta exacta del archivo que necesitas modificar.

- Prohibido usar `find`, `grep` o `code_search` sin contexto específico.
- Prohibido leer carpetas completas (`ls` recursivo, lectura masiva de archivos).
- Usa `PROJECT_MAP.md` como índice: lee el mapa, identifica el archivo, lee solo ese archivo.
- Si no encuentras el archivo en el mapa, verifica primero que el mapa esté actualizado.

## Seguimiento de tareas: REQUIREMENTS_CHECKLIST.md es la fuente de verdad

**OBLIGATORIO:** Al iniciar cualquier sesión o tarea, lee PRIMERO `REQUIREMENTS_CHECKLIST.md` para conocer:

- Hito actual en el que estamos trabajando
- Última tarea completada
- Siguiente tarea pendiente

- Prohibido explorar el código para "adivinar" el progreso del proyecto.
- Prohibido asumir qué está implementado basándote en el contenido de archivos.
- El estado del proyecto vive en `REQUIREMENTS_CHECKLIST.md`, no en inferencias del código.

## Actualización continua: marca tareas al completarlas

**OBLIGATORIO:** Cada vez que completes una tarea, prueba o endpoint con éxito, debes actualizar inmediatamente `REQUIREMENTS_CHECKLIST.md`:

- Marca la casilla correspondiente: `- [ ]` → `- [x]`
- Actualiza la sección "Estado General" al inicio del archivo:
  - "Última tarea completada": nombre de la tarea que acabas de terminar
  - "Siguiente tarea pendiente": nombre de la siguiente tarea según el orden del checklist

- Esta actualización es síncrona: no dejes la tarea marcada sin actualizar el estado general.
- Si una tarea tiene múltiples sub-items, marca cada uno a medida que los completes.

## Ahorro de tokens: lecturas puntuales y ediciones mínimas

**OBLIGATORIO:** Lee únicamente los archivos que vas a modificar de forma puntual.

- Prohibido reescribir archivos enteros si solo cambia una función.
- Usa `edit` con `old_string` específico (contexto suficiente para ser único).
- Si vas a leer múltiples archivos para contexto, hazlo en paralelo (un solo tool call).
- No leas archivos que no vas a tocar: infiere la estructura desde `PROJECT_MAP.md`.

## Prioridad de lectura al iniciar sesión

Orden obligatorio de lectura al comienzo de cualquier sesión:

1. `REQUIREMENTS_CHECKLIST.md` — Estado del proyecto
2. `PROJECT_MAP.md` — Ubicación de archivos (solo si necesitas modificar algo)
3. Archivo específico a modificar (identificado desde el mapa)

## Excepciones permitidas

- Solo puedes usar `grep` o `code_search` cuando:
  - Necesitas encontrar dónde se usa una función específica
  - Estás debuggeando un error y necesitas rastrear una llamada
  - El usuario explícitamente te pide buscar algo en el código

- En estos casos, usa patrones muy específicos (no wildcards amplios como `.*`).

## Violaciones de estas reglas

Si violas estas reglas (escaneo global, lectura sin mapa, no actualizar checklist), estarás desperdiciando tokens y haciendo el trabajo menos eficiente. El objetivo es máxima precisión con mínimo contexto.
