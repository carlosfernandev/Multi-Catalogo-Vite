# Despliegue en EC2 con Docker Compose

## Requisitos previos en EC2

1. **Instalar Docker y Docker Compose**
```bash
sudo apt update
sudo apt install -y docker.io docker-compose
sudo usermod -aG docker $USER
newgrp docker
```

2. **Clonar el repositorio**
```bash
git clone <tu-repo-url>
cd Multi-Catalogo-Vite
```

3. **Configurar variables de entorno**
```bash
cp .env.prod.example .env
# Editar .env con valores seguros
nano .env
```

## Despliegue

1. **Construir y levantar contenedores**
```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

2. **Verificar estado**
```bash
docker-compose -f docker-compose.prod.yml ps
```

3. **Ver logs**
```bash
# Todos los servicios
docker-compose -f docker-compose.prod.yml logs -f

# Servicio específico
docker-compose -f docker-compose.prod.yml logs -f backend
```

## Arquitectura de puertos

- **Frontend (nginx)**: puerto 80 (HTTP)
- **Backend Go**: puerto 3000 (interno)
- **FastAPI SaaS**: puerto 8000 (interno)
- **PostgreSQL**: puerto 5432 (interno)
- **Redis**: puerto 6379 (interno)

## Configuración HTTPS (Opcional - Recomendado)

Para producción, usa Caddy o nginx con Certbot:

### Opción A: Caddy (más simple)

Instalar Caddy:
```bash
sudo apt install -y caddy
```

Crear `/etc/caddy/Caddyfile`:
```
tu-dominio.com {
    reverse_proxy localhost:80
}
```

Reiniciar Caddy:
```bash
sudo systemctl restart caddy
```

### Opción B: nginx + Certbot

1. Instalar nginx y Certbot:
```bash
sudo apt install -y nginx certbot python3-certbot-nginx
```

2. Configurar nginx:
```bash
sudo nano /etc/nginx/sites-available/multicatalogo
```

Contenido:
```nginx
server {
    listen 80;
    server_name tu-dominio.com;

    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

3. Habilitar y obtener certificado SSL:
```bash
sudo ln -s /etc/nginx/sites-available/multicatalogo /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
sudo certbot --nginx -d tu-dominio.com
```

## Comandos útiles

- **Detener servicios**: `docker-compose -f docker-compose.prod.yml down`
- **Reiniciar servicio**: `docker-compose -f docker-compose.prod.yml restart backend`
- **Ver logs de un contenedor**: `docker logs multicatalogo-backend`
- **Entrar a un contenedor**: `docker exec -it multicatalogo-backend sh`
- **Actualizar código**: `git pull && docker-compose -f docker-compose.prod.yml up -d --build`

## Seguridad adicional

1. **Configurar firewall UFW**
```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

2. **Solo exponer puertos necesarios**
   - Solo puerto 80/443 accesibles desde internet
   - Puertos 3000, 8000, 5432, 6379 solo en la red Docker interna

## Troubleshooting

- **Contenedor no inicia**: `docker logs <container-name>`
- **Base de datos no conecta**: Verificar que postgres esté healthy: `docker-compose ps`
- **Frontend no carga API**: Verificar CORS en backend/main.go (origins permitidos)
