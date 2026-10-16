"""Servicio para la gestión de inquilinos (tenants) y catálogos en la capa SaaS."""

from datetime import datetime
from typing import Dict, List, Optional
from ..schemas import TenantCreate, TenantOut, TenantUpdate, CatalogoPersonalizadoOut, CatalogoPersonalizadoCreate, CatalogoPersonalizadoItem


class TenantService:
    def __init__(self):
        # Repositorio en memoria inicial con tenants semilla para desarrollo local
        self._tenants: Dict[int, Dict] = {
            1: {
                "id": 1,
                "nombre": "Distribuidora Guayas Cosmetics",
                "slug": "guayas-cosmetics",
                "dominio_personalizado": "guayas.multicatalogo.local",
                "contacto_email": "contacto@guayascosmetics.ec",
                "activo": True,
                "configuracion": {
                    "moneda": "USD",
                    "tema_color": "#0ea5e9",
                    "comision_extra_pct": 2.5,
                },
                "creado_en": datetime(2026, 1, 15, 10, 0, 0),
            },
            2: {
                "id": 2,
                "nombre": "Santa Elena Beauty Red",
                "slug": "santa-elena-beauty",
                "dominio_personalizado": None,
                "contacto_email": "admin@upse.edu.ec",
                "activo": True,
                "configuracion": {
                    "moneda": "USD",
                    "tema_color": "#4f46e5",
                    "comision_extra_pct": 5.0,
                },
                "creado_en": datetime(2026, 2, 1, 12, 0, 0),
            },
        }

        self._catalogos: Dict[int, Dict] = {
            1: {
                "id": 1,
                "tenant_id": 1,
                "nombre_catalogo": "Catálogo Exclusivo Costa",
                "activo": True,
                "items": [
                    {"producto_id": 1, "precio_personalizado": 42.00, "destacado": True, "orden": 1},
                    {"producto_id": 2, "precio_personalizado": 30.00, "destacado": True, "orden": 2},
                    {"producto_id": 8, "precio_personalizado": 120.00, "destacado": True, "orden": 3},
                ],
                "actualizado_en": datetime(2026, 2, 15, 14, 0, 0),
            }
        }
        self._next_tenant_id = 3
        self._next_catalogo_id = 2

    def listar_tenants(self, solo_activos: bool = True) -> List[TenantOut]:
        resultado = []
        for t in self._tenants.values():
            if not solo_activos or t["activo"]:
                resultado.append(TenantOut(**t))
        return resultado

    def obtener_tenant(self, tenant_id: int) -> Optional[TenantOut]:
        data = self._tenants.get(tenant_id)
        if not data:
            return None
        return TenantOut(**data)

    def obtener_tenant_por_slug(self, slug: str) -> Optional[TenantOut]:
        for t in self._tenants.values():
            if t["slug"] == slug:
                return TenantOut(**t)
        return None

    def crear_tenant(self, tenant: TenantCreate) -> TenantOut:
        # Validar slug único
        if self.obtener_tenant_por_slug(tenant.slug):
            raise ValueError(f"Ya existe un tenant con el slug '{tenant.slug}'")

        tenant_id = self._next_tenant_id
        self._next_tenant_id += 1

        nuevo = {
            "id": tenant_id,
            **tenant.model_dump(),
            "creado_en": datetime.now(),
        }
        self._tenants[tenant_id] = nuevo
        return TenantOut(**nuevo)

    def actualizar_tenant(self, tenant_id: int, datos: TenantUpdate) -> Optional[TenantOut]:
        if tenant_id not in self._tenants:
            return None

        actual = self._tenants[tenant_id]
        campos = datos.model_dump(exclude_unset=True)
        actual.update(campos)
        return TenantOut(**actual)

    def obtener_catalogo_por_tenant(self, tenant_id: int) -> Optional[CatalogoPersonalizadoOut]:
        for cat in self._catalogos.values():
            if cat["tenant_id"] == tenant_id:
                return CatalogoPersonalizadoOut(**cat)
        return None

    def guardar_catalogo_tenant(self, datos: CatalogoPersonalizadoCreate) -> CatalogoPersonalizadoOut:
        # Verificar tenant
        if datos.tenant_id not in self._tenants:
            raise KeyError(f"El tenant {datos.tenant_id} no existe")

        for cat_id, cat in self._catalogos.items():
            if cat["tenant_id"] == datos.tenant_id:
                cat["nombre_catalogo"] = datos.nombre_catalogo
                cat["items"] = [i.model_dump() for i in datos.items]
                cat["actualizado_en"] = datetime.now()
                return CatalogoPersonalizadoOut(**cat)

        cat_id = self._next_catalogo_id
        self._next_catalogo_id += 1
        nuevo = {
            "id": cat_id,
            "tenant_id": datos.tenant_id,
            "nombre_catalogo": datos.nombre_catalogo,
            "activo": True,
            "items": [i.model_dump() for i in datos.items],
            "actualizado_en": datetime.now(),
        }
        self._catalogos[cat_id] = nuevo
        return CatalogoPersonalizadoOut(**nuevo)


tenant_service = TenantService()
