// src/components/MiRed/services/redService.ts
// Capa de servicios para la red multinivel (feature MiRed).
// Consume el endpoint GET /api/red expuesto por la API Go/Fiber.

import type { Referido } from "../data/red";

// Obtiene la estructura jerárquica de la red de referidos desde GET /api/red
export const getRed = async (): Promise<Referido> => {
  const response = await fetch("/api/red");
  if (!response.ok) {
    throw new Error("Error al obtener la red de referidos desde la API");
  }
  return response.json();
};
