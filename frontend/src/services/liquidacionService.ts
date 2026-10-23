// src/services/liquidacionService.ts
// Servicio para consultar y procesar liquidaciones y comisiones del MLM.

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

export interface DistribuidorPendiente {
  referido_id: number;
  referido_nombre: string;
  total_pendiente: number;
  cantidad_comisiones: number;
}

export interface LiquidacionHistorial {
  id: number;
  referido_id: number;
  referido_nombre?: string;
  total: number;
  estado: string;
  creado_en: string;
  pagada_en?: string;
}

/**
 * Consulta el historial general de comisiones (usado en MiRed).
 */
export const getLiquidaciones = async (): Promise<ResumenLiquidaciones> => {
  return apiFetch<ResumenLiquidaciones>("/api/liquidaciones");
};

/**
 * Obtiene los saldos de comisiones pendientes acumulados por cada distribuidor.
 */
export const obtenerPendientesLiquidacion = async (): Promise<DistribuidorPendiente[]> => {
  return apiFetch<DistribuidorPendiente[]>("/api/liquidaciones/pendientes");
};

/**
 * Obtiene el historial de comprobantes de pago emitidos a distribuidores.
 */
export const obtenerHistorialLiquidaciones = async (): Promise<LiquidacionHistorial[]> => {
  return apiFetch<LiquidacionHistorial[]>("/api/liquidaciones/historial");
};

/**
 * Liquida (paga) formalmente el total de comisiones pendientes a un distribuidor.
 */
export const procesarLiquidacion = async (referidoId: number): Promise<LiquidacionHistorial> => {
  return apiFetch<LiquidacionHistorial>("/api/liquidaciones", {
    method: "POST",
    body: JSON.stringify({ referido_id: referidoId }),
  });
};
