# Multi-Catálogo Vite

Este proyecto está modularizado en dos partes principales: el backend y el frontend, ubicados en sus respectivas carpetas.

## Estructura del Proyecto

- `backend/`: Contiene la API y la lógica del servidor, desarrollada en Go.
- `frontend/`: Contiene la interfaz de usuario, desarrollada con React, TypeScript y Vite.

## Cómo iniciar el proyecto

Para ejecutar la aplicación completa, necesitarás levantar tanto el servidor del backend como el servidor de desarrollo del frontend en terminales separadas.

### 1. Iniciar el Backend (Go)

Abre una terminal, navega a la carpeta `backend` y ejecuta la aplicación de Go:

```bash
cd backend
go run main.go
```

Esto iniciará el servidor del backend, el cual quedará a la escucha de peticiones HTTP.

### 2. Iniciar el Frontend (React + Vite)

Abre una nueva terminal, navega a la carpeta `frontend`, instala las dependencias y arranca el entorno de desarrollo:

```bash
cd frontend
npm install      # O puedes usar: pnpm install
npm run dev      # O puedes usar: pnpm dev
```

El servidor de Vite se iniciará y te mostrará una URL local en la terminal (por defecto, suele ser `http://localhost:5173`). Abre esa URL en tu navegador web para interactuar con la aplicación.
