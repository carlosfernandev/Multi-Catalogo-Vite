// src/components/Login/services/authService.ts
// Servicio de autenticación (feature Login). Aísla el proceso POST /api/login
// para que el componente solo llame a esta función, igual que Catalogo/services.

import type { Rol } from "../../../context/auth";
import { http } from "../../../services/http";

export interface LoginResponse {
  email: string;
  rol: Rol;
  token: string;
}

export interface RegisterRequest {
  nombre: string;
  email: string;
  password: string;
  ref?: string;
}

export interface RegisterResponse {
  email: string;
  nombre: string;
  rol: Rol;
  token: string;
  referido_id?: number;
}

// Envía las credenciales a POST /api/login y normaliza la respuesta con el rol y token JWT
export const login = async (email: string, password: string): Promise<LoginResponse> => {
  const data = await http.post<{ email: string; rol: string; token: string }>("/api/login", {
    email,
    password,
  });
  const rol: Rol = data.rol === "admin" ? "admin" : "cliente";
  return { email: data.email, rol, token: data.token };
};

// Registra un nuevo usuario con código de referido opcional en POST /api/register
export const register = async (req: RegisterRequest): Promise<RegisterResponse> => {
  const data = await http.post<RegisterResponse>("/api/register", req);
  const rol: Rol = data.rol === "admin" ? "admin" : "cliente";
  return { ...data, rol };
};
