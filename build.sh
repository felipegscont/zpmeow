#!/bin/bash

# Script para build do zpmeow com configurações específicas

echo "🔨 Building zpmeow..."

# Limpar cache
echo "🧹 Cleaning cache..."
go clean -cache

# Definir variáveis de ambiente
export GOTOOLCHAIN=local
export CGO_ENABLED=0
export GOROOT=/usr/local/go

# Criar diretório bin se não existir
mkdir -p bin

# Build
echo "📦 Building binary..."
go build -v -o bin/zpmeow ./cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful! Binary created at bin/zpmeow"
    echo "📊 Binary info:"
    ls -lh bin/zpmeow
else
    echo "❌ Build failed!"
    exit 1
fi
