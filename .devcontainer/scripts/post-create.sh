#!/bin/bash
set -e

echo "Configurando entorno de desarrollo Condominio..."

# Fix git safe directory for containers
git config --global --add safe.directory /workspaces/condominio

# Create .env from example if not exists
if [ ! -f .env ]; then
    cp .env.example .env
    echo ".env creado desde .env.example"
fi

# Install dependencies
echo "Instalando dependencias..."
bun install

# Generate Prisma client
echo "Generando cliente Prisma..."
bun prisma generate

# Setup git hooks
echo "Configurando git hooks..."
bunx lefthook install

# Setup Go dependencies
echo "Configurando Go..."
cd apps/api
go mod download
cd ../..

# Install moon tools
echo "Instalando toolchains de moon..."
moon sync

echo "Entorno listo!"
echo ""
echo "Para iniciar desarrollo:"
echo "  bun dev              # Todos los servicios"
echo "  moon run panel:dev   # Solo frontend"
echo "  moon run api:serve   # Solo backend"
