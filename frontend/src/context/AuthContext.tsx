// src/context/AuthContext.tsx
import { useState } from 'react';
import {
  AuthContext,
  STORAGE_KEY,
  type Usuario,
  type AuthProviderProps,
  type Rol,
  type AuthContextType,
} from './auth';

export type { Rol, Usuario, AuthContextType, AuthProviderProps };

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
