// src/components/MiRed/services/redService.ts
// Capa de servicios para la red multinivel (feature MiRed).
// Consume el endpoint GET /api/red expuesto por la API Go/Fiber.

import type { Referido } from "../data/red";
import { apiFetch } from "../../../services/http";

// Obtiene la estructura jerárquica de la red de referidos desde GET /api/red.
// El email identifica la clave de caché en Redis (red:tree:{email}).
export const getRed = async (email?: string): Promise<Referido> => {
  const endpoint = email ? `/api/red?email=${encodeURIComponent(email)}` : "/api/red";
  return apiFetch<Referido>(endpoint);
};
