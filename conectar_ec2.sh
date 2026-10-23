#!/bin/bash

echo "Conectando a tu instancia EC2 Ubuntu por el puerto 443..."
KEY_PATH="$HOME/Documentos/aws/angel-rsa-2.pem"
EC2_HOST="ubuntu@3.144.255.71"

echo "Host: $EC2_HOST"
echo "Llave: $KEY_PATH"
echo ""

ssh -o StrictHostKeyChecking=accept-new -p 443 -i "$KEY_PATH" "$EC2_HOST"

echo ""
