import { apiFetch } from './http';
import type { IProducto } from '../components/Catalogo/interfaces/iproducto';

export interface ProductoInput {
  nombre: string;
  descripcion: string;
  precio: number;
  categoria: string;
  img: string;
  galeria: string[];
  stock: number;
  activo?: boolean;
}

/**
 * Obtiene la lista completa de productos para administración (incluyendo inactivos).
 * Soporta búsqueda opcional por texto.
 */
export async function obtenerProductosAdmin(q?: string): Promise<IProducto[]> {
  const query = q ? `?q=${encodeURIComponent(q)}` : '';
  return apiFetch<IProducto[]>(`/api/admin/productos${query}`);
}

/**
 * Registra un nuevo producto en el catálogo. Requiere rol de administrador.
 */
export async function crearProducto(data: ProductoInput): Promise<IProducto> {
  return apiFetch<IProducto>('/api/productos', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

/**
 * Actualiza los datos de un producto existente.
 */
export async function actualizarProducto(id: number, data: Partial<ProductoInput>): Promise<IProducto> {
  return apiFetch<IProducto>(`/api/productos/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

/**
 * Desactiva un producto aplicando borrado lógico.
 */
export async function eliminarProducto(id: number): Promise<{ mensaje: string; id: number }> {
  return apiFetch<{ mensaje: string; id: number }>(`/api/productos/${id}`, {
    method: 'DELETE',
  });
}

/**
 * Alterna el estado activo/inactivo de un producto.
 */
export async function toggleActivoProducto(id: number, estadoActual: boolean): Promise<IProducto> {
  return actualizarProducto(id, { activo: !estadoActual });
}
