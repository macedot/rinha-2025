#!/bin/bash

# https://github.com/andersongomes001/rinha-2025

set -e
APP_NAME="rinha-2025"
DOCKER_USER="macedot"
VERSION=$(git rev-parse --short HEAD)
IMAGE_NAME="$DOCKER_USER/$APP_NAME"

echo "🐳 Build da imagem Docker..."
docker build -t $IMAGE_NAME:$VERSION -t $IMAGE_NAME:latest .

echo "✅ Build concluído:"
echo "  - $IMAGE_NAME:$VERSION"
echo "  - $IMAGE_NAME:latest"

read -p "Deseja fazer push da imagem para Docker Hub? (s/n): " resposta
if [[ "$resposta" =~ ^[sS]$ ]]; then
    echo "🔐 Enviando imagens..."
    docker push $IMAGE_NAME:$VERSION
    docker push $IMAGE_NAME:latest
    echo "🎉 Imagens enviadas!"
fi
