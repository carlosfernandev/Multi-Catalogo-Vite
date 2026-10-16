// src/services/healthService.ts
// Servicio de verificación de estado (health check) para la API Go/Fiber.
// Consume el endpoint GET /api/health.

import { apiFetch } from './http';

export interface HealthResponse {
  status: string;
  message: string;
}

// Consulta el estado operativo del backend Go/Fiber
export const checkHealth = async (): Promise<HealthResponse> => {
  return apiFetch<HealthResponse>("/api/health");
};
