# Exposición Temporal con ngrok (Opción de Emergencia)

En caso de no poder usar VPS/EC2, puedes exponer temporalmente el proyecto usando ngrok + nginx local.

## Requisitos

1. **nginx** (proxy local):
```bash
sudo dnf install nginx
```

2. **ngrok** (tunnel HTTPS):
```bash
sudo dnf install ngrok
# O descargar desde https://ngrok.com/download
```

3. **docker** (para exponer-ngrok.sh):
```bash
sudo dnf install docker
sudo systemctl start docker
```

## Opción 1: Exponer con Docker Compose (Automático)

Inicia todos los servicios (frontend, backend, saas, postgres, redis) en Docker y expone con ngrok:

```bash
cd /home/carlosfpatino/Workspaces/Multi-Catalogo-Vite
./exponer-ngrok.sh
```

Este script:
- Verifica dependencias
- Inicia docker-compose con todos los servicios
- Inicia nginx local en puerto 8080
- Inicia ngrok para exponer 8080
- Hace cleanup al terminar (Ctrl+C)

## Opción 2: Exponer con Servicios Manuales (Más Control)

Inicia los servicios manualmente y luego expone con ngrok:

```bash
# Terminal 1: Frontend
cd frontend
pnpm run build
pnpm run preview

# Terminal 2: Backend Go
cd backend
go run .

# Terminal 3: FastAPI SaaS
cd Unidad3_SaaS
python main.py

# Terminal 4: Exponer con ngrok
cd /home/carlosfpatino/Workspaces/Multi-Catalogo-Vite
./exponer-ngrok-simple.sh
```

## Arquitectura

```
ngrok (HTTPS público)
    ↓
nginx local (puerto 8080)
    ↓
├─ / → frontend (localhost:4173)
├─ /api/tx/* → backend Go (localhost:3000)
├─ /api/saas/* → FastAPI (localhost:8000)
└─ /api/* → backend Go (localhost:3000)
```

## Configuración nginx

El archivo `nginx-local.conf` define:
- Upstreams para cada servicio
- Reglas de proxy por ruta
- Headers X-Real-IP y X-Forwarded-For

## Limitaciones ngrok Free

- URL aleatoria (ej: `https://a1b2c3d4.ngrok-free.app`)
- Cambia cada sesión
- Sin dominio custom
- Límite de conexiones simultáneas

## Para producción real

Esta opción es **solo temporal**. Para producción:
- Usa VPS (EC2, DigitalOcean, Oracle Cloud)
- Configura dominio propio
- SSL permanente (Let's Encrypt)
- Firewall UFW/security groups

## Detener exposición

Presiona `Ctrl+C` en el terminal de ngrok. El script hará cleanup automáticamente.

## Troubleshooting

**nginx no inicia:**
```bash
# Verificar puerto 8080 libre
sudo lsof -i :8080

# Matar proceso si está ocupado
sudo kill -9 $(sudo lsof -t -i:8080)
```

**ngrok no inicia:**
- Verifica que ngrok esté instalado
- Si faltan dependencias, instala manualmente desde ngrok.com

**Servicios Docker no responden:**
```bash
docker-compose -f docker-compose.prod.yml logs
```
