// src/components/MiRed/services/redService.ts
// Capa de servicios para la red multinivel (feature MiRed).
// Consume el endpoint GET /api/red expuesto por la API Go/Fiber.

import type { Referido } from "../data/red";
import { getStoredToken } from "../../../context/AuthContext";

// Obtiene la estructura jerárquica de la red de referidos desde GET /api/red.
// El email identifica la clave de caché en Redis (red:tree:{email}).
export const getRed = async (email?: string): Promise<Referido> => {
  const url = email ? `/api/red?email=${encodeURIComponent(email)}` : "/api/red";
  const headers: Record<string, string> = {};
  const token = getStoredToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(url, { headers });
  if (!response.ok) {
    if (response.status === 401) {
      throw new Error("Sesión no autorizada o expirada. Por favor, inicia sesión nuevamente.");
    }
    throw new Error("Error al obtener la red de referidos desde la API");
  }
  return response.json();
};
