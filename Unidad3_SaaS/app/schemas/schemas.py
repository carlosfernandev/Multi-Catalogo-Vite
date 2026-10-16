"""Esquemas Pydantic para validación y serialización de datos de la capa SaaS."""

from datetime import datetime
from typing import Any, Dict, List, Optional
from pydantic import BaseModel, EmailStr, Field


class APIErrorResponse(BaseModel):
    """Contrato estándar de error unificado con Go/Fiber: {status, message, details}."""
    status: int
    message: str
    details: Optional[Any] = None


class HealthResponse(BaseModel):
    status: str
    service: str = "fastapi-saas"
    version: str = "1.0.0"
    dependencies: Dict[str, str] = Field(default_factory=dict)


class TenantBase(BaseModel):
    nombre: str = Field(..., min_length=2, max_length=100, description="Nombre de la marca o distribuidor")
    slug: str = Field(..., min_length=2, max_length=50, description="Identificador único para subdominio o ruta")
    dominio_personalizado: Optional[str] = None
    contacto_email: EmailStr
    activo: bool = True
    configuracion: Dict[str, Any] = Field(
        default_factory=lambda: {
            "moneda": "USD",
            "tema_color": "#4f46e5",
            "comision_extra_pct": 0.0,
        }
    )


class TenantCreate(TenantBase):
    pass


class TenantUpdate(BaseModel):
    nombre: Optional[str] = None
    dominio_personalizado: Optional[str] = None
    contacto_email: Optional[EmailStr] = None
    activo: Optional[bool] = None
    configuracion: Optional[Dict[str, Any]] = None


class TenantOut(TenantBase):
    id: int
    creado_en: datetime

    model_config = {"from_attributes": True}


class CatalogoPersonalizadoItem(BaseModel):
    producto_id: int
    precio_personalizado: Optional[float] = None
    destacado: bool = False
    orden: int = 0


class CatalogoPersonalizadoCreate(BaseModel):
    tenant_id: int
    nombre_catalogo: str
    items: List[CatalogoPersonalizadoItem]


class CatalogoPersonalizadoOut(BaseModel):
    id: int
    tenant_id: int
    nombre_catalogo: str
    activo: bool
    items: List[CatalogoPersonalizadoItem]
    actualizado_en: datetime

    model_config = {"from_attributes": True}
