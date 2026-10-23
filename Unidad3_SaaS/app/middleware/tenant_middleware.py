"""Middleware para resolución de inquilino (Tenant) por subdominio, header o parámetro."""

from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import Response
from ..services import tenant_service


class TenantSubdomainMiddleware(BaseHTTPMiddleware):
    """
    Middleware multitenant:
    1. Revisa header 'X-Tenant-Slug' o 'X-Tenant'
    2. Revisa parámetro query '?tenant=slug'
    3. Revisa subdominio en Host (ej: guayas-cosmetics.multicatalogo.local)
    4. Revisa dominio completo en 'dominio_personalizado'
    Si se detecta un tenant válido, lo asocia a request.state.tenant y agrega el header en la respuesta.
    """

    async def dispatch(self, request: Request, call_next) -> Response:
        tenant = None
        slug_candidato = None

        # 1. Header explícito
        header_slug = request.headers.get("x-tenant-slug") or request.headers.get("x-tenant")
        if header_slug:
            slug_candidato = header_slug.strip().lower()

        # 2. Query param ?tenant=
        if not slug_candidato and "tenant" in request.query_params:
            slug_candidato = request.query_params["tenant"].strip().lower()

        # 3. Subdominio o Host
        host = request.headers.get("host", "").split(":")[0].lower()
        if not slug_candidato and host:
            # Buscar por dominio personalizado exacto
            tenant_por_dominio = tenant_service.obtener_tenant_por_dominio(host)
            if tenant_por_dominio:
                tenant = tenant_por_dominio
            else:
                # Extraer subdominio si aplica (ej. {sub}.multicatalogo.local o {sub}.localhost)
                partes = host.split(".")
                if len(partes) > 1 and partes[0] not in ("www", "api", "localhost", "127"):
                    slug_candidato = partes[0]

        if not tenant and slug_candidato:
            tenant = tenant_service.obtener_tenant_por_slug(slug_candidato)

        # Inyectar en el estado de la petición
        request.state.tenant = tenant
        request.state.tenant_slug = tenant.slug if tenant else slug_candidato

        response: Response = await call_next(request)

        # Agregar encabezado en la respuesta para rastreo
        if tenant:
            response.headers["X-Resolved-Tenant"] = tenant.slug
            response.headers["X-Tenant-Id"] = str(tenant.id)

        return response
