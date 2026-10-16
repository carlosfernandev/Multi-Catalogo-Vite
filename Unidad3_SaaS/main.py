"""Punto de entrada del microservicio FastAPI para la capa SaaS Multitenant (Unidad 3)."""

from contextlib import asynccontextmanager
from fastapi import FastAPI, Request, HTTPException
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.config import settings
from app.routers import health_router, tenants_router, catalogos_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Inicialización del microservicio SaaS
    print(f"🚀 Iniciando {settings.APP_NAME} en entorno [{settings.APP_ENV}]...")
    yield
    print(f"🛑 Deteniendo {settings.APP_NAME}...")


app = FastAPI(
    title=settings.APP_NAME,
    description="Microservicio SaaS Multitenant de Multi-Catálogo (Unidad 3). "
                "Gestiona inquilinos, personalización de catálogos y configuración de distribuidores.",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    lifespan=lifespan,
)

# Configuración de CORS para permitir consumo desde el SPA React Vite
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:5173",
        "http://127.0.0.1:5173",
        "http://localhost:3000",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# =====================================================================
# Manejadores de excepciones para contrato unificado {status, message, details}
# =====================================================================
@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException):
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "status": exc.status_code,
            "message": exc.detail if isinstance(exc.detail, str) else "Error HTTP",
            "details": exc.detail if not isinstance(exc.detail, str) else None,
        },
    )


@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    detalles = exc.errors() if settings.DEBUG else None
    return JSONResponse(
        status_code=400,
        content={
            "status": 400,
            "message": "Datos de petición inválidos",
            "details": detalles,
        },
    )


@app.exception_handler(Exception)
async def generic_exception_handler(request: Request, exc: Exception):
    detalles = str(exc) if settings.DEBUG else None
    return JSONResponse(
        status_code=500,
        content={
            "status": 500,
            "message": "Error interno del servidor SaaS",
            "details": detalles,
        },
    )


# =====================================================================
# Registro de Rutas bajo el prefijo /api/saas
# =====================================================================
app.include_router(health_router, prefix="/api/saas")
app.include_router(tenants_router, prefix="/api/saas")
app.include_router(catalogos_router, prefix="/api/saas")


@app.get("/", tags=["Root"])
def root():
    return {
        "servicio": settings.APP_NAME,
        "version": "1.0.0",
        "documentacion": "/docs",
        "endpoints": {
            "health": "/api/saas/health",
            "tenants": "/api/saas/tenants",
            "catalogos": "/api/saas/catalogos",
        },
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
    )
