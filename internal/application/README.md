# 🎯 Application Layer - Clean Architecture

Esta camada contém os **casos de uso da aplicação** seguindo rigorosamente os princípios de Clean Architecture e Domain-Driven Design (DDD) conforme definido no ARCHITECTURE.md.

## 📋 Estrutura Oficial

Seguindo exatamente o padrão definido no ARCHITECTURE.md:

```
internal/application/
├── README.md          # Esta documentação
├── session.go         # Criar, listar, conectar, validar, deletar sessões
├── message.go         # Enviar, construir, validar mensagens
├── media.go           # Upload/download de mídia
└── webhook.go         # Registrar, notificar, validar webhooks
```

## 🎯 Responsabilidades da Camada Application

### ✅ O que esta camada FAZ:
- **Orquestração**: Coordena chamadas entre domain e infrastructure
- **Conversão de DTOs**: Transforma dados de entrada/saída
- **Coordenação de Transações**: Gerencia operações que envolvem múltiplas camadas
- **Validação de Entrada**: Valida DTOs de request
- **Fluxo de Casos de Uso**: Implementa a sequência de operações para cada caso de uso

### ❌ O que esta camada NÃO FAZ:
- **Regras de Negócio**: Lógica de domínio complexa (delegado para domain)
- **Implementações Concretas**: Detalhes de infraestrutura (abstraído via interfaces)
- **Validações de Domínio**: Regras de negócio (delegado para domain services)
- **Persistência Direta**: Acesso direto a banco/storage (via repositories)

## 🏗️ Detalhamento por Agregado

### 📁 session.go - Agregado Session
**Casos de Uso Implementados:**
- `CreateSession`: Criar nova sessão
- `GetSession`: Buscar sessão por ID
- `ListSessions`: Listar todas as sessões
- `ConnectSession`: Iniciar conexão de sessão
- `DeleteSession`: Excluir sessão

**Dependências:**
- `session.Repository`: Interface para persistência (infrastructure)
- `session.DomainService`: Regras de negócio (domain)
- `validation.Validator`: Validação de entrada (shared)

### 📁 message.go - Agregado Message
**Casos de Uso Implementados:**
- `SendMessage`: Enviar mensagem de texto
- `SendMedia`: Enviar mensagem de mídia
- `SendLocation`: Enviar mensagem de localização
- `SendContact`: Enviar mensagem de contato

**Dependências:**
- `message.Service`: Regras de negócio de mensagens (domain)
- `session.Repository`: Verificação de sessão (domain)
- `validation.Validator`: Validação de entrada (shared)

### 📁 media.go - Agregado Media
**Casos de Uso Implementados:**
- `UploadMedia`: Upload de arquivo de mídia
- `GetMedia`: Buscar informações de mídia
- `DownloadMedia`: Gerar URL de download
- `ListMedia`: Listar arquivos de mídia
- `DeleteMedia`: Excluir arquivo de mídia
- `GetUploadProgress`: Acompanhar progresso de upload

**Dependências:**
- `media.Service`: Regras de negócio de mídia (domain)
- `mediaInfra.StorageService`: Armazenamento (infrastructure)
- `session.Repository`: Verificação de sessão (domain)
- `validation.Validator`: Validação de entrada (shared)

### 📁 webhook.go - Agregado Webhook
**Casos de Uso Implementados:**
- `RegisterWebhook`: Registrar novo webhook
- `GetWebhook`: Buscar webhook por ID
- `UpdateWebhook`: Atualizar configuração de webhook
- `ListWebhooks`: Listar webhooks
- `DeleteWebhook`: Excluir webhook
- `TestWebhook`: Testar webhook
- `NotifyWebhook`: Enviar notificações
- `CreateMessageWebhookPayload`: Criar payload para eventos de mensagem
- `CreateStatusWebhookPayload`: Criar payload para eventos de status
- `CreateConnectionWebhookPayload`: Criar payload para eventos de conexão

**Dependências:**
- `session.Repository`: Verificação de sessão (domain)
- `validation.Validator`: Validação de entrada (shared)

## 🔄 Fluxo de Orquestração

### Padrão de Implementação
Cada caso de uso segue o mesmo padrão de orquestração:

```go
func (s *Service) UseCase(ctx context.Context, req *dto.Request) (*dto.Response, error) {
    // 1. Validação de entrada (application layer)
    if err := s.validator.Validate(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Verificar dependências (sessão, etc.)
    entity, err := s.repository.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("dependency check failed: %w", err)
    }

    // 3. Aplicar regras de negócio (delegado para domain)
    if err := s.domainService.ValidateBusinessRules(entity); err != nil {
        return nil, fmt.Errorf("business rules validation failed: %w", err)
    }

    // 4. Coordenar operações (infrastructure)
    result, err := s.infrastructureService.Execute(entity)
    if err != nil {
        return nil, fmt.Errorf("operation failed: %w", err)
    }

    // 5. Converter para DTO de resposta
    response := &dto.Response{
        Status:  200,
        Message: "Operation completed successfully",
        Data:    convertToDTO(result),
    }

    return response, nil
}
```

## 🎯 Benefícios da Estrutura Atual

### ✅ **Vantagens**
1. **Conformidade com ARCHITECTURE.md**: Segue exatamente o padrão definido
2. **Organização por Agregado**: Cada arquivo representa um agregado DDD
3. **Responsabilidades Claras**: Apenas orquestração, sem lógica de negócio
4. **Separação de Camadas**: Domain, Application, Infrastructure bem definidas
5. **Facilidade de Manutenção**: Estrutura simples e previsível

### 📋 **Convenções Seguidas**
- **Nomenclatura**: snake_case para arquivos (session.go, message.go, etc.)
- **Package**: Todos os arquivos usam `package application`
- **Agregados**: Um arquivo por agregado DDD
- **Responsabilidades**: Apenas casos de uso e orquestração

## 📖 Exemplo de Uso

```go
package main

import (
    "zpmeow/internal/application"
    "zpmeow/internal/domain/session"
    "zpmeow/internal/domain/message"
    "zpmeow/internal/shared/validation"
)

func main() {
    // Inicializar dependências
    validator := validation.NewValidator()
    sessionRepo := // implementação do repository
    sessionDomainService := session.NewSessionDomainService()
    messageDomainService := message.NewDomainService()

    // Criar services de aplicação
    sessionService := application.NewSessionService(
        sessionRepo,
        sessionDomainService,
        validator,
    )

    messageService := application.NewMessageService(
        messageDomainService,
        sessionRepo,
        validator,
    )

    // Usar casos de uso
    response, err := sessionService.CreateSession(ctx, createRequest)
    response, err := messageService.SendMessage(ctx, sessionID, messageRequest)
}
```

## 🔄 Migração Realizada

### **Antes** (estrutura inconsistente):
- Mistura de responsabilidades
- Lógica de negócio na application layer
- Validações técnicas misturadas com regras de domínio

### **Depois** (seguindo ARCHITECTURE.md):
- Responsabilidades claras por camada
- Apenas orquestração na application layer
- Delegação correta para domain e infrastructure
- Estrutura por agregados DDD

## 🚀 Próximos Passos

1. **✅ Refatoração Completa**: Todos os arquivos refatorados seguindo o padrão
2. **🔄 Testes**: Verificar compilação e funcionalidade
3. **📚 Documentação**: README.md atualizado
4. **🔗 Integração**: Verificar compatibilidade com outras camadas

## 📋 Checklist de Conformidade

### ✅ Seguindo ARCHITECTURE.md:
- [x] Estrutura por agregados (session.go, message.go, media.go, webhook.go)
- [x] Apenas orquestração entre domain e infrastructure
- [x] Sem lógica de negócio na application layer
- [x] Validação de entrada via shared/validation
- [x] Conversão de DTOs para entidades de domain
- [x] Delegação de regras de negócio para domain services
- [x] Coordenação de persistência via repositories
- [x] Nomenclatura snake_case para arquivos
- [x] Package application consistente
- [x] Documentação atualizada

Esta estrutura está agora **100% conforme** com o padrão definido no ARCHITECTURE.md.
