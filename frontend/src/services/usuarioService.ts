// src/services/usuarioService.ts
// Servicio para administración de usuarios.

import { apiFetch } from "./http";

export interface UsuarioAdmin {
  id: number;
  email: string;
  nombre?: string;
  rol: string;
  activo: boolean;
  referido_id?: number;
  creado_en?: string;
}

export interface ActualizarUsuarioRequest {
  nombre: string;
  rol: string;
  activo: boolean;
}

/**
 * Obtiene la lista de usuarios para administración con filtro opcional.
 */
export const obtenerUsuariosAdmin = async (q?: string): Promise<UsuarioAdmin[]> => {
  const url = q ? `/api/admin/usuarios?q=${encodeURIComponent(q)}` : "/api/admin/usuarios";
  return apiFetch<UsuarioAdmin[]>(url);
};

/**
 * Actualiza un usuario (nombre, rol, activo).
 */
export const actualizarUsuarioAdmin = async (id: number, req: ActualizarUsuarioRequest): Promise<UsuarioAdmin> => {
  return apiFetch<UsuarioAdmin>(`/api/usuarios/${id}`, {
    method: "PATCH",
    body: JSON.stringify(req),
  });
};

/**
 * Alterna el estado activo/inactivo de un usuario.
 */
export const toggleActivoUsuario = async (id: number, activoActual: boolean): Promise<UsuarioAdmin> => {
  return apiFetch<UsuarioAdmin>(`/api/usuarios/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ activo: !activoActual }),
  });
};
