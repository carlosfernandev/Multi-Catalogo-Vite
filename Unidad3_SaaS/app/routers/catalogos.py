"""Router para gestión de catálogos personalizados por Tenant."""

from fastapi import APIRouter, HTTPException, status
from ..schemas import CatalogoPersonalizadoCreate, CatalogoPersonalizadoOut, APIErrorResponse
from ..services import tenant_service

router = APIRouter(prefix="/catalogos", tags=["Catálogos SaaS"])


@router.get(
    "/{tenant_id}",
    response_model=CatalogoPersonalizadoOut,
    summary="Obtener catálogo personalizado por inquilino",
    responses={
        404: {"model": APIErrorResponse, "description": "Catálogo no encontrado para este tenant"},
    },
)
def obtener_catalogo_tenant(tenant_id: int):
    cat = tenant_service.obtener_catalogo_por_tenant(tenant_id)
    if not cat:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"No se encontró catálogo personalizado para el inquilino {tenant_id}",
        )
    return cat


@router.post(
    "",
    response_model=CatalogoPersonalizadoOut,
    status_code=status.HTTP_200_OK,
    summary="Crear o actualizar catálogo personalizado para inquilino",
    responses={
        404: {"model": APIErrorResponse, "description": "Inquilino no encontrado"},
    },
)
def guardar_catalogo_tenant(datos: CatalogoPersonalizadoCreate):
    try:
        return tenant_service.guardar_catalogo_tenant(datos)
    except KeyError as err:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(err),
        )
