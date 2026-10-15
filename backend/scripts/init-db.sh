#!/bin/bash
# Script para inicializar la base de datos manualmente después de levantar PostgreSQL

echo "Esperando a que PostgreSQL esté listo..."
until docker exec multicatalogo_db pg_isready -U multicatalogo_user -d multicatalogo > /dev/null 2>&1; do
  sleep 2
done

echo "PostgreSQL listo. Ejecutando init.sql..."
docker exec multicatalogo_db psql -U multicatalogo_user -d multicatalogo -f - < database/init.sql

echo "Inicialización completada."
