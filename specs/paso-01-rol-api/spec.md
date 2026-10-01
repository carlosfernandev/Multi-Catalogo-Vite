# Paso 01 — Rol en la API

**Entregable:** base de F8
**Correspondencia con la planificación:** ACD 2
**Fuente:** `practica03.md` → Paso 1

## Objetivo

Habilitar que la API distinga dos perfiles al autenticar: `admin` y `cliente`. Este paso es la base del flujo F8 (login con rol y navegación por rol); sin el campo `rol` en la respuesta, el frontend no puede decidir a dónde redirigir ni qué rutas proteger.

## Concepto

La autenticación distingue dos perfiles: `admin` (gestiona el negocio) y `cliente` (compra en línea). El backend sigue con datos hardcodeados; la persistencia real con PostgreSQL se abordará en la Unidad 2.

## Archivos

| Acción | Ruta |
|---|---|
| modificar | `backend/controllers/authController.go` |

## Alcance técnico

- `switch` sobre `req.Email` + `req.Password` con dos casos válidos; `default` retorna 401.
- Respuesta exitosa incluye `token` (ficticio), `email` y `rol` (`"admin"` / `"cliente"`).
- Sin cambios en `models` ni rutas: solo el cuerpo de `Login`.

## Criterios de aceptación

- [x] `POST /api/login` con `admin@upse.edu.ec / 123456` retorna `rol: "admin"`.
- [x] `POST /api/login` con `cliente@upse.edu.ec / 123456` retorna `rol: "cliente"`.
- [x] Credenciales inválidas retornan 401 con `error`.
- [x] `go build ./...` compila sin errores.

## Verificación

```bash
cd backend && go build ./...
curl -X POST http://localhost:3000/api/login -H "Content-Type: application/json" \
  -d '{"email":"admin@upse.edu.ec","password":"123456"}'
# → {"email":"admin@upse.edu.ec","rol":"admin","token":"fake-jwt-token-123"}
```
