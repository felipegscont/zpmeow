#!/bin/bash

echo "🔧 Corrigindo imports nos handlers..."

cd internal/infra/http/handlers

# Corrigir imports comuns em todos os handlers
for file in *.go; do
    if [[ "$file" != "utils.go" ]]; then
        echo "   📝 Corrigindo $file"
        
        # Substituir imports antigos
        sed -i 's|zpmeow/internal/domain/session|zpmeow/internal/domain|g' "$file"
        sed -i 's|zpmeow/internal/types|zpmeow/internal/application|g' "$file"
        sed -i 's|zpmeow/internal/shared/utils|zpmeow/internal/shared|g' "$file"
        sed -i 's|zpmeow/internal/infra/http/utils|zpmeow/internal/infra/http/handlers|g' "$file"
        
        # Substituir referências de tipos
        sed -i 's|session\.SessionService|domain.SessionService|g' "$file"
        sed -i 's|session\.Session|domain.Session|g' "$file"
        sed -i 's|types\.|application.|g' "$file"
        sed -i 's|utils\.RespondWithError|RespondWithError|g' "$file"
        sed -i 's|utils\.RespondWithJSON|RespondWithJSON|g' "$file"
        sed -i 's|utils\.ValidateSessionIDParam|ValidateSessionIDParam|g' "$file"
        sed -i 's|utils\.ValidateAndBindJSON|ValidateAndBindJSON|g' "$file"
        
        # Corrigir referências de service
        sed -i 's|service\.MeowServiceImpl|infra.MeowServiceImpl|g' "$file"
        
        # Adicionar import do infra se necessário
        if grep -q "MeowServiceImpl" "$file"; then
            if ! grep -q "zpmeow/internal/infra" "$file"; then
                sed -i '/zpmeow\/internal\/domain/a\\t"zpmeow/internal/infra"' "$file"
            fi
        fi
    fi
done

echo "✅ Handlers corrigidos!"
