// src/services/liquidacionService.ts
// Servicio para consultar el historial de liquidaciones y comisiones registradas en backend.

import { apiFetch } from "./http";

export interface LiquidacionRegistro {
  id: number;
  pedido_id: number;
  pedido_numero: string;
  referido_id: number;
  referido_nombre: string;
  nivel: number;
  monto: number;
  fecha: string;
}

export interface ResumenLiquidaciones {
  total_comisiones: number;
  cantidad: number;
  liquidaciones: LiquidacionRegistro[];
}

export const getLiquidaciones = async (): Promise<ResumenLiquidaciones> => {
  return apiFetch<ResumenLiquidaciones>("/api/liquidaciones");
};
