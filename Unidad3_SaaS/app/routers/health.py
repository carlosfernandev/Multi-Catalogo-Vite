"""Router para health check y observabilidad de la capa SaaS."""

import socket
from fastapi import APIRouter
from ..config import settings
from ..schemas import HealthResponse

router = APIRouter(prefix="/health", tags=["Health"])


def verificar_puerto_tcp(host: str, port: int, timeout: float = 0.5) -> bool:
    try:
        with socket.create_connection((host, port), timeout=timeout):
            return True
    except OSError:
        return False


@router.get("", response_model=HealthResponse)
def get_health():
    """Retorna el estado operativo del microservicio FastAPI SaaS y sus dependencias."""
    postgres_ok = verificar_puerto_tcp(settings.DB_HOST, settings.DB_PORT)
    redis_ok = verificar_puerto_tcp(settings.REDIS_HOST, settings.REDIS_PORT)

    deps = {
        "postgres": "ok" if postgres_ok else "unreachable",
        "redis": "ok" if redis_ok else "unreachable",
    }

    # El servicio puede degradar si Redis no está disponible, pero reporta estado global
    status = "ok" if postgres_ok else ("degraded" if redis_ok else "unhealthy")

    return HealthResponse(
        status=status,
        service="fastapi-saas",
        version="1.0.0",
        dependencies=deps,
    )
