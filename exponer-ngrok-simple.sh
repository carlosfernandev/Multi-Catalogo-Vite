#!/bin/bash

echo "=== Exponer con ngrok (servicios ya corriendo localmente) ==="
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

echo "📋 Asegúrate de tener corriendo:"
echo "   - Frontend preview: cd frontend && pnpm run preview (puerto 4173)"
echo "   - Backend Go: cd backend && go run . (puerto 3000)"
echo "   - FastAPI SaaS: cd Unidad3_SaaS && python main.py (puerto 8000)"
echo ""
read -p "¿Servicios corriendo? (y/n): " confirm

if [ "$confirm" != "y" ]; then
    echo "❌ Cancelado. Inicia los servicios primero."
    exit 1
fi

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
echo "🛑 Deteniendo nginx..."
sudo nginx -s stop
