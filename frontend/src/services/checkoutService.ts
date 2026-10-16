// src/services/checkoutService.ts
// Servicio de checkout: envía el carrito a POST /api/checkout (API Go/Fiber).
// El backend guarda el pedido en una transacción ACID y calcula el total con
// los precios reales de PostgreSQL; aquí solo se envían IDs y cantidades.

import { http } from "./http";

export interface ClienteCheckout {
  nombre: string;
  email: string;
  direccion: string;
  ciudad: string;
  telefono: string;
}

export interface CheckoutRequest {
  usuario_email: string;
  cliente: ClienteCheckout;
  metodo_pago: string;
  items: { producto_id: number; cantidad: number }[];
}

export interface PedidoResponse {
  pedido_id: number;
  numero: string;
  total: number;
  estado: string;
  fecha: string;
  items: {
    producto_id: number;
    nombre: string;
    cantidad: number;
    precio_unitario: number;
    subtotal: number;
  }[];
}

// Envía el pedido; si la API responde con error (400/500) lanza su mensaje
export const crearPedido = async (pedido: CheckoutRequest): Promise<PedidoResponse> => {
  return http.post<PedidoResponse>("/api/checkout", pedido);
};
