# Directivas de Ahorro de Tokens (MCP) .windsurfrules

## Uso de Playwright
- NUNCA tomes capturas de pantalla (screenshots) ni pidas el DOM completo/HTML de la página.
- Extrae únicamente el texto visible o el selector CSS puntual del elemento que se está testeando.
- Para verificar APIs o respuestas JSON, usa peticiones HTTP directas en lugar de abrir el navegador.

## Uso de Memory
- Guarda en memoria persistente ÚNICAMENTE:
  1. Causa raíz y solución final de errores de configuración o entorno recurrentes.
  2. Puertos locales y variables clave del proyecto.
- PROHIBIDO guardar código fuente completo, trazas de error (stack traces) o historial de chat en la memoria.

## Manejo de Errores y Pruebas
- Si una prueba falla, ejecuta SOLO ese test unitario individual, nunca la suite completa.
- Si el mismo error persiste tras 2 intentos de solución, DETENTE y explica la hipótesis conceptual antes de editar más archivos.