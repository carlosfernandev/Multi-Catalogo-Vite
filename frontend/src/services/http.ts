// src/services/http.ts
// Cliente HTTP centralizado para el frontend de Multi-Catálogo.
// Gestiona baseUrl (VITE_API_URL), Authorization Bearer automático y ApiError tipado.

import { getStoredToken } from '../context/auth';

/**
 * Error personalizado para respuestas no exitosas de la API REST.
 * Preserva el código de estado HTTP y los detalles técnicos devueltos por el backend.
 */
export class ApiError extends Error {
  status: number;
  details?: string;

  constructor(message: string, status: number, details?: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.details = details;
  }
}

const RAW_API_URL = import.meta.env.VITE_API_URL || '';
export const BASE_API_URL = RAW_API_URL.endsWith('/')
  ? RAW_API_URL.slice(0, -1)
  : RAW_API_URL;

/**
 * Construye la URL completa anteponiendo BASE_API_URL si está configurada.
 */
export const buildApiUrl = (endpoint: string): string => {
  if (!BASE_API_URL) return endpoint;
  const path = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  return `${BASE_API_URL}${path}`;
};

/**
 * Función fetch unificada que adjunta automáticamente el token JWT y procesa errores tipados.
 */
export async function apiFetch<T = unknown>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = buildApiUrl(endpoint);

  const headers = new Headers(options.headers || {});

  // Adjuntar token de sesión si existe en localStorage y no se especificó un Authorization manual
  const token = getStoredToken();
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  // Establecer Content-Type JSON si hay cuerpo y no es un form-data nativo
  if (
    options.body &&
    !(options.body instanceof FormData) &&
    !(options.body instanceof Blob) &&
    !headers.has('Content-Type')
  ) {
    headers.set('Content-Type', 'application/json');
  }

  let response: Response;
  try {
    response = await fetch(url, {
      ...options,
      headers,
    });
  } catch (err) {
    throw new ApiError(
      'No se pudo establecer conexión con el servidor. Verifica tu conexión o intenta más tarde.',
      0,
      err instanceof Error ? err.message : String(err)
    );
  }

  if (!response.ok) {
    let errorMessage = response.statusText || 'Error en la petición';
    let errorDetails: string | undefined;

    try {
      const data = await response.json();
      if (data && typeof data === 'object') {
        if ('message' in data && typeof data.message === 'string') {
          errorMessage = data.message;
        }
        if ('details' in data && typeof data.details === 'string') {
          errorDetails = data.details;
        }
      }
    } catch {
      // La respuesta no era JSON (ej. 502/504 de proxy o texto plano)
    }

    if (response.status === 401 && !errorMessage) {
      errorMessage = 'Sesión no autorizada o expirada. Por favor, inicia sesión nuevamente.';
    }

    throw new ApiError(errorMessage, response.status, errorDetails);
  }

  // Si la respuesta no tiene contenido (204 No Content)
  if (response.status === 204) {
    return undefined as unknown as T;
  }

  return response.json() as Promise<T>;
}

export const http = {
  get: <T = unknown>(url: string, options?: RequestInit) =>
    apiFetch<T>(url, { ...options, method: 'GET' }),

  post: <T = unknown>(url: string, body?: unknown, options?: RequestInit) =>
    apiFetch<T>(url, {
      ...options,
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),

  put: <T = unknown>(url: string, body?: unknown, options?: RequestInit) =>
    apiFetch<T>(url, {
      ...options,
      method: 'PUT',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),

  delete: <T = unknown>(url: string, options?: RequestInit) =>
    apiFetch<T>(url, { ...options, method: 'DELETE' }),
};
