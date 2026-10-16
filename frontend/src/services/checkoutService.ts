// src/services/checkoutService.ts
// Servicio de checkout: envía el carrito a POST /api/checkout (API Go/Fiber).
// El backend guarda el pedido en una transacción ACID y calcula el total con
// los precios reales de PostgreSQL; aquí solo se envían IDs y cantidades.

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

import { getStoredToken } from "../context/AuthContext";

// Envía el pedido; si la API responde con error (400/500) lanza su mensaje
export const crearPedido = async (pedido: CheckoutRequest): Promise<PedidoResponse> => {
  let response: Response;
  const token = getStoredToken();
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  try {
    response = await fetch("/api/checkout", {
      method: "POST",
      headers,
      body: JSON.stringify(pedido),
    });
  } catch {
    throw new Error("No se pudo conectar con el servidor. Inténtalo de nuevo.");
  }

  const data = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(data?.message ?? "No se pudo procesar el pedido");
  }
  return data as PedidoResponse;
};
