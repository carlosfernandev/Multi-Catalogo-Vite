"""Router para gestión de Tenants (inquilinos SaaS / distribuidores)."""

from typing import List, Optional
from fastapi import APIRouter, HTTPException, Query, status
from ..schemas import TenantCreate, TenantOut, TenantUpdate, APIErrorResponse
from ..services import tenant_service

router = APIRouter(prefix="/tenants", tags=["Tenants"])


@router.get(
    "",
    response_model=List[TenantOut],
    summary="Listar inquilinos SaaS",
    description="Devuelve todos los tenants o distribuidores registrados.",
)
def listar_tenants(solo_activos: bool = Query(True, description="Filtrar solo inquilinos activos")):
    return tenant_service.listar_tenants(solo_activos=solo_activos)


@router.post(
    "",
    response_model=TenantOut,
    status_code=status.HTTP_201_CREATED,
    summary="Registrar nuevo inquilino",
    responses={
        400: {"model": APIErrorResponse, "description": "Slug duplicado o datos inválidos"},
    },
)
def crear_tenant(tenant: TenantCreate):
    try:
        return tenant_service.crear_tenant(tenant)
    except ValueError as err:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(err),
        )


@router.get(
    "/{tenant_id}",
    response_model=TenantOut,
    summary="Obtener inquilino por ID",
    responses={
        404: {"model": APIErrorResponse, "description": "Inquilino no encontrado"},
    },
)
def obtener_tenant(tenant_id: int):
    t = tenant_service.obtener_tenant(tenant_id)
    if not t:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Inquilino con ID {tenant_id} no encontrado",
        )
    return t


@router.put(
    "/{tenant_id}",
    response_model=TenantOut,
    summary="Actualizar inquilino",
    responses={
        404: {"model": APIErrorResponse, "description": "Inquilino no encontrado"},
    },
)
def actualizar_tenant(tenant_id: int, datos: TenantUpdate):
    t = tenant_service.actualizar_tenant(tenant_id, datos)
    if not t:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Inquilino con ID {tenant_id} no encontrado",
        )
    return t
