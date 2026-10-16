// src/components/Catalogo/services/productosService.ts
// Capa de servicios del catálogo, organizada por feature (Catalogo).
// Consume los endpoints reales de la API Go/Fiber (/api/productos).

import type { IProducto } from "../interfaces/iproducto";
import { apiFetch, ApiError } from "../../../services/http";

// Obtiene la lista completa de productos desde GET /api/productos
export const getProductos = async (): Promise<IProducto[]> => {
  return apiFetch<IProducto[]>("/api/productos");
};

// Obtiene un producto por ID desde GET /api/productos/:id
export const getProductoById = async (id: number): Promise<IProducto | undefined> => {
  try {
    return await apiFetch<IProducto>(`/api/productos/${id}`);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) {
      return undefined;
    }
    throw err;
  }
};