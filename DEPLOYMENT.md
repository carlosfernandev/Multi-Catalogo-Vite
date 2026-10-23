# Despliegue en EC2 con Docker Compose

Guía de despliegue y puesta en marcha en servidores Linux / AWS EC2 para la plataforma MultiCatálogo.

---

## 1. Requisitos previos en el servidor

### 1.1 Instalar Docker y Docker Compose
```bash
sudo apt update
sudo apt install -y docker.io docker-compose-v2
# O en distribuciones más antiguas: sudo apt install -y docker.io docker-compose
```

### 1.2 Configurar permisos de usuario (Evitar "permission denied unix:///var/run/docker.sock")
Si al ejecutar comandos de docker recibes el error `permission denied while trying to connect to the docker API at unix:///var/run/docker.sock`, añade tu usuario al grupo `docker`:
```bash
sudo usermod -aG docker $USER
newgrp docker
```
> [!TIP]
> Si la terminal te sigue pidiendo contraseña o no reconoce el grupo, cierra la sesión SSH y vuelve a conectarte (`exit` y reconectar) para que los permisos se recarguen de inmediato.

---

## 2. Clonar y Configurar el Proyecto

1. **Clonar el repositorio**:
   ```bash
   git clone https://github.com/carlosfernandev/Multi-Catalogo-Vite.git
   cd Multi-Catalogo-Vite
   ```

2. **Asegurar que estás en la rama principal**:
   ```bash
   git checkout main
   git pull origin main
   ```

3. **Crear archivo de variables de entorno de producción**:
   ```bash
   cp .env.prod.example .env
   # Editar con contraseñas seguras y secret de JWT
   nano .env
   ```

---

## 3. Despliegue con Docker Compose

> [!NOTE]
> En Docker moderno (v2), el comando oficial es `docker compose` (con espacio). Si tu servidor utiliza la versión 1 de Compose, usa `docker-compose` (con guión).

1. **Construir y levantar todos los servicios en segundo plano**:
   ```bash
   docker compose -f docker-compose.prod.yml up -d --build
   # O con Compose v1:
   # docker-compose -f docker-compose.prod.yml up -d --build
   ```

2. **Verificar el estado de los 5 contenedores**:
   ```bash
   docker compose -f docker-compose.prod.yml ps
   ```
   Debes ver los 5 contenedores en estado `Up` / `healthy`:
   - `multicatalogo-postgres` (healthy)
   - `multicatalogo-redis` (healthy)
   - `multicatalogo-backend` (healthy)
   - `multicatalogo-saas` (healthy)
   - `multicatalogo-frontend` (healthy)

3. **Ver logs de ejecución**:
   ```bash
   # Logs generales
   docker compose -f docker-compose.prod.yml logs -f

   # Logs específicos del backend o base de datos
   docker compose -f docker-compose.prod.yml logs -f backend
   docker compose -f docker-compose.prod.yml logs -f postgres
   ```

---

## 4. Arquitectura de Red y Puertos

El contenedor `multicatalogo-frontend` incluye un servidor **Nginx** que actúa como punto de entrada unificado y reverse-proxy hacia los microservicios:

| Ruta Externa | Destino Interno | Descripción |
|---|---|---|
| `http://<ip-servidor>:4173/` | Frontend Nginx (puerto 80) | SPA React 19 + Vite PWA |
| `http://<ip-servidor>:4173/api/` | `http://backend:3000/api/` | API Go/Fiber (Auth, Catálogo, Pedidos, Liquidaciones) |
| `http://<ip-servidor>:4173/saas/` | `http://saas:8000/` | Microservicio FastAPI SaaS Multi-Tenant |

- **Ventaja de esta arquitectura:** Al usar el proxy inverso integrado en el Frontend, **se elimina cualquier problema de CORS** y no es necesario exponer los puertos de PostgreSQL (5432), Redis (6379), Backend (3000) ni SaaS (8000) a Internet en el Security Group de AWS.

> [!NOTE]
> En entornos Linux rootless (como Podman/Fedora), el frontend se expone por defecto en el puerto `4173` para no requerir privilegios de superusuario en puertos `< 1024`.
> En EC2, si deseas exponer directamente en el puerto HTTP estándar `80`, edita `docker-compose.prod.yml` en la sección `frontend` cambiando `- "4173:80"` por `- "80:80"`.

---

## 5. Configuración de Dominio y Certificado SSL / HTTPS (Recomendado)

### Opción A: Caddy (Automático con Let's Encrypt)
1. Instalar Caddy:
   ```bash
   sudo apt install -y caddy
   ```
2. Editar `/etc/caddy/Caddyfile`:
   ```caddy
   tu-dominio.com {
       reverse_proxy localhost:4173
   }
   ```
3. Reiniciar:
   ```bash
   sudo systemctl restart caddy
   ```

### Opción B: Nginx + Certbot
```bash
sudo apt install -y nginx certbot python3-certbot-nginx
```
Crear `/etc/nginx/sites-available/multicatalogo`:
```nginx
server {
    listen 80;
    server_name tu-dominio.com;

    location / {
        proxy_pass http://localhost:4173;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
Habilitar y certificar:
```bash
sudo ln -s /etc/nginx/sites-available/multicatalogo /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
sudo certbot --nginx -d tu-dominio.com
```

---

## 6. Mantenimiento y Troubleshooting

- **Actualizar código y reconstruir:**
  ```bash
  git pull origin main
  docker compose -f docker-compose.prod.yml up -d --build
  ```

- **Detener todos los servicios:**
  ```bash
  docker compose -f docker-compose.prod.yml down
  ```

- **Reiniciar un servicio específico:**
  ```bash
  docker compose -f docker-compose.prod.yml restart backend
  ```

- **Aplicar migraciones SQL sobre un volumen existente:**
  PostgreSQL ejecuta los scripts de `/docker-entrypoint-initdb.d/` solo la primera vez que se crea el volumen. Si ya tenías un volumen antiguo y necesitas aplicar la migración de Backoffice (`005_backoffice.sql`), ejecuta:
  ```bash
  docker exec -i multicatalogo-postgres psql -U multicatalogo -d multicatalogo < backend/database/005_backoffice.sql
  ```

- **Verificar salud de la base de datos y Redis:**
  ```bash
  curl -s http://localhost:3000/api/health
  # O vía proxy del frontend:
  curl -s http://localhost:4173/api/health
  ```
