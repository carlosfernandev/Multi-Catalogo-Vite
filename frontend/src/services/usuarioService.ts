import { apiFetch } from './http';

export interface UsuarioAdmin {
  id: number;
  email: string;
  nombre?: string;
  rol: 'admin' | 'cliente';
  activo: boolean;
  referido_id?: number | null;
  creado_en?: string;
}

export interface ActualizarUsuarioInput {
  nombre?: string;
  rol?: 'admin' | 'cliente';
  activo?: boolean;
}

/**
 * Obtiene la lista completa de usuarios para el panel de administración.
 * Permite filtrar por nombre o correo electrónico mediante el parámetro q.
 */
export async function obtenerUsuariosAdmin(q?: string): Promise<UsuarioAdmin[]> {
  const query = q ? `?q=${encodeURIComponent(q)}` : '';
  return apiFetch<UsuarioAdmin[]>(`/api/usuarios${query}`);
}

/**
 * Actualiza los datos de un usuario (rol, nombre o estado activo/inactivo).
 */
export async function actualizarUsuarioAdmin(
  id: number,
  data: ActualizarUsuarioInput
): Promise<UsuarioAdmin> {
  return apiFetch<UsuarioAdmin>(`/api/usuarios/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

/**
 * Cambia el rol de un usuario entre 'admin' y 'cliente'.
 */
export async function cambiarRolUsuario(
  id: number,
  nuevoRol: 'admin' | 'cliente'
): Promise<UsuarioAdmin> {
  return actualizarUsuarioAdmin(id, { rol: nuevoRol });
}

/**
 * Activa o desactiva la cuenta de un usuario.
 */
export async function toggleActivoUsuario(
  id: number,
  estadoActual: boolean
): Promise<UsuarioAdmin> {
  return actualizarUsuarioAdmin(id, { activo: !estadoActual });
}
