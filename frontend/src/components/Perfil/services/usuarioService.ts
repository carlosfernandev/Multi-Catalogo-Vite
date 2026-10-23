// src/components/Perfil/services/usuarioService.ts
// Servicio para obtener el perfil del usuario autenticado.

import { apiFetch } from "../../../services/http";

export interface UsuarioPerfil {
  id: number;
  email: string;
  nombre: string;
  rol: string;
  activo: boolean;
  referido_id?: number;
  creado_en: string;
  referido_codigo?: string;
  referido_nivel?: number;
  referido_ventas?: number;
}

/**
 * Obtiene el perfil completo del usuario autenticado con su código de referido.
 */
export const obtenerPerfil = async (): Promise<UsuarioPerfil> => {
  return apiFetch<UsuarioPerfil>("/api/me");
};
