#!/bin/bash

echo "=== Iniciando stack local para exponer con ngrok ==="
echo ""

# Verificar si nginx está instalado
if ! command -v nginx &> /dev/null; then
    echo "❌ nginx no está instalado. Instalándolo..."
    sudo dnf install -y nginx
fi

# Verificar si ngrok está instalado
if ! command -v ngrok &> /dev/null; then
    echo "❌ ngrok no está instalado."
    echo "Instálalo con: curl o desde https://ngrok.com/download"
    exit 1
fi

# Verificar si docker/podman está instalado
if ! command -v docker &> /dev/null && ! command -v podman &> /dev/null; then
    echo "❌ docker/podman no está instalado."
    exit 1
fi

echo "✅ Dependencias verificadas"
echo ""

# Crear .env para docker-compose si no existe
if [ ! -f .env ]; then
    echo "Creando .env desde .env.prod.example..."
    cp .env.prod.example .env
fi

echo "🚀 Iniciando servicios con docker-compose..."
docker-compose -f docker-compose.prod.yml up -d

echo "⏳ Esperando que servicios estén listos..."
sleep 5

echo "🏗️  Iniciando nginx local en puerto 8080..."
NGINX_CONF=$(pwd)/nginx-local.conf
NGINX_TMP=/tmp/nginx-local
mkdir -p $NGINX_TMP/logs $NGINX_TMP/temp

# Si ya hay un nginx corriendo con esa configuración, recargar o continuar
if [ -f $NGINX_TMP/nginx.pid ] && kill -0 $(cat $NGINX_TMP/nginx.pid) 2>/dev/null; then
    echo "nginx ya está en ejecución (PID $(cat $NGINX_TMP/nginx.pid)). Recargando..."
    nginx -s reload -c "$NGINX_CONF"
else
    nginx -c "$NGINX_CONF"
fi

echo "🌐 Iniciando ngrok para exponer puerto 8080..."
echo "⚠️  ngrok mostrará la URL pública. Cópiala para compartirla."
echo ""
ngrok http 8080 --domain=ornamentally-experientialistic-desire.ngrok-free.dev

# Cleanup al terminar
echo ""
echo "🛑 Deteniendo servicios..."
if [ -f $NGINX_TMP/nginx.pid ]; then
    nginx -s stop -c "$NGINX_CONF" 2>/dev/null || true
fi
docker-compose -f docker-compose.prod.yml down
