#!/bin/bash

echo "=== Iniciando stack local para exponer con ngrok ==="
echo ""

# Verificar si nginx está instalado
if ! command -v nginx &> /dev/null; then
    echo "❌ nginx no está instalado. Instalando..."
    sudo dnf install -y nginx
fi

# Verificar si ngrok está instalado
if ! command -v ngrok &> /dev/null; then
    echo "❌ ngrok no está instalado."
    echo "Instálalo con: sudo dnf install ngrok"
    echo "O descarga desde: https://ngrok.com/download"
    exit 1
fi

# Verificar si docker está instalado
if ! command -v docker &> /dev/null; then
    echo "❌ docker no está instalado."
    echo "Instálalo con: sudo dnf install -y docker"
    exit 1
fi

# Iniciar docker si no está corriendo
sudo systemctl start docker

echo "✅ Dependencias verificadas"
echo ""

# Crear .env para docker-compose si no existe
if [ ! -f .env ]; then
    echo "Creando .env desde .env.prod.example..."
    cp .env.prod.example .env
    echo "⚠️  Edita .env con tus valores seguros antes de continuar"
    read -p "Presiona Enter para continuar..."
fi

echo "🚀 Iniciando servicios con docker-compose..."
docker-compose -f docker-compose.prod.yml up -d

echo "⏳ Esperando que servicios estén listos (30 segundos)..."
sleep 30

echo "🏗️  Iniciando nginx local en puerto 8080..."
NGINX_CONF=$(pwd)/nginx-local.conf
NGINX_TMP=/tmp/nginx-local
mkdir -p $NGINX_TMP
sudo nginx -c $NGINX_CONF -p $NGINX_TMP

echo "🌐 Iniciando ngrok para exponer puerto 8080..."
echo "⚠️  ngrok mostrará la URL pública. Cópiala para compartirla."
echo ""
ngrok http 8080

# Cleanup al terminar
echo ""
echo "🛑 Deteniendo servicios..."
sudo nginx -s stop
docker-compose -f docker-compose.prod.yml down
