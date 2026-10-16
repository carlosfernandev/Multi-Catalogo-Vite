"""Módulo de esquemas Pydantic."""
from .schemas import (
    APIErrorResponse,
    HealthResponse,
    TenantBase,
    TenantCreate,
    TenantUpdate,
    TenantOut,
    CatalogoPersonalizadoItem,
    CatalogoPersonalizadoCreate,
    CatalogoPersonalizadoOut,
)

__all__ = [
    "APIErrorResponse",
    "HealthResponse",
    "TenantBase",
    "TenantCreate",
    "TenantUpdate",
    "TenantOut",
    "CatalogoPersonalizadoItem",
    "CatalogoPersonalizadoCreate",
    "CatalogoPersonalizadoOut",
]
