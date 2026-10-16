# Unidad 3 — Capa SaaS Multitenant (FastAPI)

Microservicio en Python con FastAPI para la gestión de inquilinos (tenants), catálogos personalizados por distribuidor y configuración multitenant de **Multi-Catálogo**.

## Arquitectura

- **Framework**: FastAPI (Python 3.10+)
- **Puerto**: 8000 (configurable vía `PORT`)
- **Prefijo base**: `/api/saas`
- **Contrato de error**: Unificado con Go/Fiber `{status: int, message: str, details?: any}`
- **Documentación Swagger**: `http://localhost:8000/docs`

## Endpoints

| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/` | Información del servicio y enlaces a docs. |
| `GET` | `/api/saas/health` | Health check con verificación de dependencias (Postgres/Redis). |
| `GET` | `/api/saas/tenants` | Listar todos los inquilinos/distribuidores. |
| `POST` | `/api/saas/tenants` | Registrar un nuevo inquilino SaaS con slug y tema. |
| `GET` | `/api/saas/tenants/{id}` | Detalle de un inquilino por ID. |
| `PUT` | `/api/saas/tenants/{id}` | Actualizar datos y configuración de un inquilino. |
| `GET` | `/api/saas/catalogos/{tenant_id}` | Obtener catálogo de productos personalizado de un tenant. |
| `POST` | `/api/saas/catalogos` | Crear o actualizar catálogo de productos personalizado. |

## Puesta en marcha local

```bash
cd Unidad3_SaaS

# Crear entorno virtual (opcional pero recomendado)
python3 -m venv venv
source venv/bin/activate

# Instalar dependencias
pip install -r requirements.txt

# Ejecutar el servicio
python3 main.py
# o con uvicorn directamente:
uvicorn main:app --reload --port 8000
```
