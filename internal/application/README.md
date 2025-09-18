# 🎯 Application Layer - Clean Architecture

Esta camada contém os **casos de uso da aplicação** (Use Cases) seguindo rigorosamente os princípios de Clean Architecture e idiomaticidade do Go.

## 📋 Estrutura Atual

```
internal/application/
├── README.md          # Esta documentação
├── interfaces.go      # Interfaces para infraestrutura
├── session.go         # Casos de uso de sessões
├── chat.go           # Casos de uso de chat
├── messaging.go      # Casos de uso de mensagens
├── webhook.go        # Casos de uso de webhooks
├── dispatcher.go     # Despachador de eventos
├── conversion.go     # Conversores de dados
├── group.go          # Casos de uso de grupos
├── newsletter.go     # Casos de uso de newsletter
└── whatsapp.go       # Casos de uso específicos do WhatsApp
```

## 🎯 Responsabilidades da Application Layer

### ✅ O que esta camada PODE e DEVE fazer:

#### **1. Orquestração de Use Cases**
- Coordenar chamadas entre Domain e Infrastructure
- Implementar fluxos de casos de uso complexos
- Gerenciar transações que envolvem múltiplos agregados

#### **2. Dependências Permitidas**
- ✅ **Domain Layer**: `internal/domain/*` - Usar entidades e services
- ✅ **DTOs**: `internal/interfaces/dto` - Para comunicação com interfaces
- ✅ **Shared Utilities**: `internal/shared/*` - Validação, tipos, etc.
- ✅ **Bibliotecas Externas**: `github.com/google/uuid`, etc.
- ✅ **Standard Library**: `context`, `fmt`, `time`, etc.

#### **3. Definição de Interfaces**
- Definir interfaces para Infrastructure (Dependency Inversion)
- Abstrair dependências externas via interfaces
- Permitir injeção de dependências

#### **4. Conversão de Dados**
- Converter entre DTOs e entidades de Domain
- Transformar dados de entrada/saída
- Mapear estruturas entre camadas

#### **5. Validação de Entrada**
- Validar DTOs de request
- Verificar parâmetros de entrada
- Sanitizar dados antes de passar para Domain

### ❌ O que esta camada NÃO PODE fazer:

#### **1. Dependências Proibidas**
- ❌ **Infrastructure**: `internal/infrastructure/*` - Violação de dependência
- ❌ **Interface Handlers**: `internal/interfaces/http/*` - Inversão incorreta
- ❌ **Detalhes de Implementação**: Banco, HTTP, filesystem diretamente

#### **2. Responsabilidades Proibidas**
- ❌ **Regras de Negócio**: Lógica complexa de domínio (vai para Domain)
- ❌ **Implementações Concretas**: Detalhes de infraestrutura
- ❌ **Validações de Domínio**: Regras de negócio (delegado para Domain)
- ❌ **Persistência Direta**: Acesso direto a banco/storage

## 🏗️ Padrões de Implementação

### **1. Dependency Inversion Pattern**

A Application Layer define interfaces para Infrastructure:

```go
// ✅ CORRETO: Application define interface
type WebhookSender interface {
    SendWebhook(ctx context.Context, url string, payload interface{}) error
}

type SessionApp struct {
    sessionRepo    session.Repository    // Domain interface
    webhookSender  WebhookSender        // Application interface
    validator      *validation.Validator // Shared utility
}

// ❌ INCORRETO: Importar infrastructure diretamente
import "meow/internal/infrastructure/webhooks" // VIOLAÇÃO!
```

### **2. Use Case Pattern**

Cada caso de uso segue o padrão de orquestração:

```go
func (s *SessionApp) CreateSession(ctx context.Context, req *dto.CreateSessionRequest) (*dto.SessionResponse, error) {
    // 1. Validação de entrada (Application responsibility)
    if err := s.validator.Validate(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Criar entidade de domínio (Domain responsibility)
    sessionID, err := session.NewSessionID(req.ID)
    if err != nil {
        return nil, fmt.Errorf("invalid session ID: %w", err)
    }

    sessionName, err := session.NewSessionName(req.Name)
    if err != nil {
        return nil, fmt.Errorf("invalid session name: %w", err)
    }

    // 3. Aplicar regras de negócio (Domain responsibility)
    sessionEntity, err := session.NewSession(sessionID, sessionName, proxyURL)
    if err != nil {
        return nil, fmt.Errorf("failed to create session: %w", err)
    }

    // 4. Persistir (Infrastructure via interface)
    if err := s.sessionRepo.Create(ctx, sessionEntity); err != nil {
        return nil, fmt.Errorf("failed to save session: %w", err)
    }

    // 5. Converter para DTO de resposta (Application responsibility)
    return &dto.SessionResponse{
        ID:     sessionEntity.ID.Value(),
        Name:   sessionEntity.Name.Value(),
        Status: string(sessionEntity.Status),
    }, nil
}
```

### **3. Interface Segregation**

Interfaces pequenas e específicas:

```go
// ✅ CORRETO: Interfaces específicas
type MessageSender interface {
    SendTextMessage(ctx context.Context, sessionID, chatJID, content string) error
}

type MediaUploader interface {
    UploadImage(ctx context.Context, data []byte) (string, error)
}

// ❌ INCORRECTTO: Interface muito grande
type MegaService interface {
    SendMessage(...)
    UploadMedia(...)
    CreateSession(...)
    // ... 50 métodos
}
```

## 📋 Regras de Dependência

### **✅ DEPENDÊNCIAS PERMITIDAS**

| Tipo | Exemplo | Justificativa |
|------|---------|---------------|
| **Standard Library** | `context`, `fmt`, `time` | Sempre permitido |
| **Domain Layer** | `internal/domain/session` | Application usa Domain |
| **DTOs** | `internal/interfaces/dto` | Comunicação com interfaces |
| **Shared Utilities** | `internal/shared/validation` | Utilitários compartilhados |
| **External Libraries** | `github.com/google/uuid` | Bibliotecas específicas |

### **❌ DEPENDÊNCIAS PROIBIDAS**

| Tipo | Exemplo | Por que é proibido |
|------|---------|-------------------|
| **Infrastructure** | `internal/infrastructure/database` | Violação de dependência |
| **Interface Handlers** | `internal/interfaces/http` | Inversão incorreta |
| **Frameworks** | `gin`, `echo` | Detalhes de implementação |

### **🔄 COMO CORRIGIR VIOLAÇÕES**

#### **Problema**: Application importando Infrastructure
```go
// ❌ INCORRETO
import "meow/internal/infrastructure/webhooks"

type WebhookApp struct {
    webhookService *webhooks.Service // Dependência direta!
}
```

#### **Solução**: Definir interface na Application
```go
// ✅ CORRETO
type WebhookSender interface {
    SendWebhook(ctx context.Context, url string, payload interface{}) error
}

type WebhookApp struct {
    webhookSender WebhookSender // Interface definida na Application
}

// Infrastructure implementa a interface
func NewWebhookApp(sender WebhookSender) *WebhookApp {
    return &WebhookApp{webhookSender: sender}
}
```

## 🎯 Benefícios da Arquitetura Correta

### ✅ **Vantagens da Application Layer**
1. **Testabilidade**: Fácil de testar com mocks das interfaces
2. **Flexibilidade**: Pode trocar implementações de Infrastructure
3. **Reutilização**: Use cases podem ser reutilizados em diferentes interfaces
4. **Manutenibilidade**: Responsabilidades bem definidas
5. **Evolução**: Fácil adicionar novos casos de uso

### 📋 **Convenções Go Idiomáticas**
- **Interfaces pequenas**: Preferir interfaces específicas
- **Dependency Injection**: Via construtores, não globals
- **Error Handling**: Sempre retornar erros explícitos
- **Context**: Sempre primeiro parâmetro
- **Naming**: Interfaces terminam com -er quando possível

## 📖 Exemplo Completo

```go
// interfaces.go - Definir contratos
type MessageSender interface {
    SendTextMessage(ctx context.Context, sessionID, chatJID, content string) error
}

// messaging.go - Implementar use case
type MessageApp struct {
    sessionRepo   session.Repository
    messageSender MessageSender
    validator     *validation.Validator
}

func NewMessageApp(repo session.Repository, sender MessageSender, validator *validation.Validator) *MessageApp {
    return &MessageApp{
        sessionRepo:   repo,
        messageSender: sender,
        validator:     validator,
    }
}

func (m *MessageApp) SendMessage(ctx context.Context, req *dto.SendMessageRequest) error {
    // 1. Validar entrada
    if err := m.validator.Validate(req); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // 2. Verificar sessão existe
    session, err := m.sessionRepo.GetByID(ctx, req.SessionID)
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }

    // 3. Verificar regras de negócio (Domain)
    if !session.IsConnected() {
        return fmt.Errorf("session not connected")
    }

    // 4. Executar operação (Infrastructure via interface)
    if err := m.messageSender.SendTextMessage(ctx, req.SessionID, req.ChatJID, req.Content); err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }

    return nil
}
```

## 🔍 Checklist de Conformidade

### ✅ **Dependências Corretas**
- [ ] Apenas stdlib Go
- [ ] Domain layer (`internal/domain/*`)
- [ ] DTOs (`internal/interfaces/dto`)
- [ ] Shared utilities (`internal/shared/*`)
- [ ] Bibliotecas externas específicas

### ❌ **Dependências Proibidas**
- [ ] Infrastructure (`internal/infrastructure/*`)
- [ ] Interface handlers (`internal/interfaces/http/*`)
- [ ] Frameworks web diretamente

### 🏗️ **Padrões Implementados**
- [ ] Dependency Inversion (interfaces definidas na Application)
- [ ] Use Case pattern (um método por caso de uso)
- [ ] Error handling idiomático
- [ ] Context como primeiro parâmetro
- [ ] Validação de entrada
- [ ] Conversão DTO ↔ Domain

### 📋 **Estrutura de Arquivos**
- [ ] Um arquivo por domínio/contexto
- [ ] Interfaces em arquivo separado
- [ ] Conversores em arquivo separado
- [ ] Nomenclatura clara e consistente

## 🚀 Status Atual

A Application Layer está **em conformidade** com Clean Architecture e idiomaticidade Go, seguindo corretamente:

- ✅ **Dependency Rule**: Depende apenas de camadas internas
- ✅ **Interface Segregation**: Interfaces pequenas e específicas
- ✅ **Dependency Inversion**: Application define interfaces para Infrastructure
- ✅ **Single Responsibility**: Cada arquivo tem responsabilidade clara
- ✅ **Go Idioms**: Seguindo convenções da linguagem

**Status**: 🎯 **ARQUITETURA CORRETA E IDIOMÁTICA**
