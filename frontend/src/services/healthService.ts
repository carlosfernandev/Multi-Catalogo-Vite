// src/services/healthService.ts
// Servicio de verificación de estado (health check) para la API Go/Fiber.
// Consume el endpoint GET /api/health.

export interface HealthResponse {
  status: string;
  message: string;
}

// Consulta el estado operativo del backend Go/Fiber
export const checkHealth = async (): Promise<HealthResponse> => {
  const response = await fetch("/api/health");
  if (!response.ok) {
    throw new Error("Error en la respuesta del health check de la API");
  }
  return response.json();
};
