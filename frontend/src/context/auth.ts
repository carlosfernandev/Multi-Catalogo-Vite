// src/context/auth.ts
// Tipos, constantes y hook de autenticación separados del Provider para cumplir con Fast Refresh de ESLint.
import { createContext, useContext, type ReactNode } from 'react';

// 1. Tipos de rol que maneja la aplicación
export type Rol = 'admin' | 'cliente';

// 2. Usuario autenticado: correo, rol y token JWT
export interface Usuario {
  email: string;
  rol: Rol;
  token?: string;
}

export const STORAGE_KEY = "multicatalogo_sesion";

// Función utilitaria para obtener el token JWT almacenado en localStorage
export const getStoredToken = (): string | null => {
  try {
    const sesion = localStorage.getItem(STORAGE_KEY);
    if (!sesion) return null;
    const parsed = JSON.parse(sesion);
    return parsed?.token ?? null;
  } catch {
    return null;
  }
};

export interface AuthContextType {
  isAuthenticated: boolean;
  user: Usuario | null;
  login: (usuario: Usuario) => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe ser usado dentro de un AuthProvider");
  }
  return context;
};

export interface AuthProviderProps {
  children: ReactNode;
}
