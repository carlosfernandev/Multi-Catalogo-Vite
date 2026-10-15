// src/components/Catalogo/services/productosService.ts
// Capa de servicios del catálogo, organizada por feature (Catalogo).
// Consume los endpoints reales de la API Go/Fiber (/api/productos).

import type { IProducto } from "../interfaces/iproducto";

// Obtiene la lista completa de productos desde GET /api/productos
export const getProductos = async (): Promise<IProducto[]> => {
  const response = await fetch("/api/productos");
  if (!response.ok) {
    throw new Error("Error al obtener productos desde la API");
  }
  return response.json();
};

// Obtiene un producto por ID desde GET /api/productos/:id
export const getProductoById = async (id: number): Promise<IProducto | undefined> => {
  const response = await fetch(`/api/productos/${id}`);
  if (response.status === 404) {
    return undefined;
  }
  if (!response.ok) {
    throw new Error("Error al buscar producto en la API");
  }
  return response.json();
};