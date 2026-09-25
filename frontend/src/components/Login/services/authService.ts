// src/components/Login/services/authService.ts
// Servicio de autenticación (feature Login). Aísla el proceso POST /api/login
// para que el componente solo llame a esta función, igual que Catalogo/services.
// En la Unidad 2 este cuerpo se reemplazará por JWT real sin tocar Login.tsx.

import type { Rol } from "../../../context/AuthContext";

export interface LoginResponse {
  email: string;
  rol: Rol;
}

// Envía las credenciales a POST /api/login y normaliza la respuesta con el rol
export const login = (email: string, password: string): Promise<LoginResponse> => {
  return fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
    .then((response) => {
      if (!response.ok) throw new Error("Credenciales incorrectas");
      return response.json();
    })
    .then((data) => {
      // La API devuelve el rol (admin | cliente) junto al correo (Tema 5)
      const rol: Rol = data.rol === "admin" ? "admin" : "cliente";
      return { email: data.email, rol };
    });
};
