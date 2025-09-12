#!/bin/bash

echo "🔧 Removendo package duplicado..."

cd internal/infra/http/handlers

for file in *.go; do
    if [[ "$file" != "utils.go" ]]; then
        echo "   📝 Corrigindo $file"
        # Remover linhas duplicadas do package
        sed -i '/^)$/,/^package handlers$/{
            /^package handlers$/d
        }' "$file"
        
        # Remover linhas vazias extras
        sed -i '/^)$/,/^type/{
            /^$/N
            /^\n$/d
        }' "$file"
    fi
done

echo "✅ Package duplicado removido!"
