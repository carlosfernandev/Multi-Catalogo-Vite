// src/context/AuthContext.tsx
import { createContext, useContext, useState, type ReactNode } from 'react';

// 1. Tipos de rol que maneja la aplicación
export type Rol = 'admin' | 'cliente';

// 2. Usuario autenticado: correo, rol y token JWT
export interface Usuario {
  email: string;
  rol: Rol;
  token?: string;
}

const STORAGE_KEY = "multicatalogo_sesion";

// Función utilitaria para obtener el token JWT almacenado
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

interface AuthContextType {
  isAuthenticated: boolean;
  user: Usuario | null;
  login: (usuario: Usuario) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe ser usado dentro de un AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<Usuario | null>(() => {
    try {
      const sesion = localStorage.getItem(STORAGE_KEY);
      if (!sesion) return null;
      const parsed = JSON.parse(sesion);
      return parsed?.email && parsed?.rol ? parsed : null;
    } catch {
      return null;
    }
  });

  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    try {
      const sesion = localStorage.getItem(STORAGE_KEY);
      return Boolean(sesion);
    } catch {
      return false;
    }
  });

  const login = (usuario: Usuario) => {
    setIsAuthenticated(true);
    setUser(usuario);
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(usuario));
    } catch (err) {
      console.error("Error al persistir sesión en localStorage:", err);
    }
  };

  const logout = () => {
    setIsAuthenticated(false);
    setUser(null);
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch (err) {
      console.error("Error al remover sesión de localStorage:", err);
    }
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
