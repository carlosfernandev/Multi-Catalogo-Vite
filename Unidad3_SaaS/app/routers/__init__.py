"""Módulo de routers de la capa SaaS."""
from .health import router as health_router
from .tenants import router as tenants_router
from .catalogos import router as catalogos_router

__all__ = ["health_router", "tenants_router", "catalogos_router"]
