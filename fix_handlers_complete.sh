#!/bin/bash

echo "🔧 Corrigindo todos os handlers sistematicamente..."

cd internal/infra/http/handlers

# Primeiro, corrigir todos os package names
echo "   📝 Corrigindo package names..."
for file in *.go; do
    sed -i 's/^package handler$/package handlers/' "$file"
done

# Segundo, padronizar imports básicos
echo "   📝 Padronizando imports..."
for file in *.go; do
    if [[ "$file" != "utils.go" ]]; then
        # Criar backup temporário
        cp "$file" "$file.tmp"
        
        # Extrair apenas o conteúdo após os imports
        awk '
        BEGIN { in_import = 0; import_done = 0 }
        /^import \(/ { in_import = 1; next }
        /^\)/ && in_import { in_import = 0; import_done = 1; next }
        in_import { next }
        /^import / && !in_import { next }
        import_done || (!in_import && !/^import/) { print }
        ' "$file.tmp" > "$file.body"
        
        # Criar novo arquivo com imports padronizados
        cat > "$file" << 'EOF'
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"zpmeow/internal/application"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/shared"
)

EOF
        
        # Adicionar o corpo do arquivo
        cat "$file.body" >> "$file"
        
        # Limpar arquivos temporários
        rm "$file.tmp" "$file.body"
    fi
done

# Terceiro, corrigir referências de tipos
echo "   📝 Corrigindo referências de tipos..."
for file in *.go; do
    if [[ "$file" != "utils.go" ]]; then
        # Corrigir tipos de domínio
        sed -i 's/session\.SessionService/domain.SessionService/g' "$file"
        sed -i 's/session\.Session/domain.Session/g' "$file"
        
        # Corrigir tipos de aplicação
        sed -i 's/types\./application./g' "$file"
        sed -i 's/whatsapp\./application./g' "$file"
        sed -i 's/webhook\./application./g' "$file"
        
        # Corrigir tipos de infraestrutura
        sed -i 's/service\.MeowService/\*infra.MeowServiceImpl/g' "$file"
        sed -i 's/service\.MeowServiceImpl/infra.MeowServiceImpl/g' "$file"
        
        # Corrigir funções utilitárias
        sed -i 's/utils\.RespondWithError/RespondWithError/g' "$file"
        sed -i 's/utils\.RespondWithJSON/RespondWithJSON/g' "$file"
        sed -i 's/utils\.ValidateSessionIDParam/ValidateSessionIDParam/g' "$file"
        sed -i 's/utils\.ValidateAndBindJSON/ValidateAndBindJSON/g' "$file"
        
        # Corrigir shared utilities
        sed -i 's/sharedUtils\./shared./g' "$file"
    fi
done

echo "✅ Handlers corrigidos sistematicamente!"
