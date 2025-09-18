# Domain Layer - DDD Clean Architecture

Esta camada contém apenas **conceitos de domínio puros**, seguindo rigorosamente os princípios de Domain-Driven Design (DDD).

## 🏛️ Princípios DDD Aplicados

### Core Concepts
- **Entities**: Objetos com identidade única e ciclo de vida
- **Value Objects**: Objetos imutáveis definidos por seus valores
- **Aggregates**: Cluster de objetos tratados como uma unidade
- **Aggregate Root**: Ponto de entrada único para o aggregate
- **Domain Services**: Lógica que não pertence a uma entidade específica
- **Repository Interfaces**: Contratos para persistência (implementados na infra)
- **Domain Events**: Eventos que representam mudanças importantes no domínio

### Design Principles
- **Independência de Infraestrutura**: Zero dependências externas
- **Regras de Negócio Puras**: Apenas lógica de domínio
- **Linguagem Ubíqua**: Termos do negócio WhatsApp API
- **Invariantes de Domínio**: Regras que sempre devem ser verdadeiras
- **Encapsulamento**: Estado interno protegido

## 📁 Estrutura DDD Limpa

```
internal/domain/
├── session/                    # Bounded Context: Session Management
│   ├── aggregate.go           # Session Aggregate Root
│   ├── entity.go              # Session Entity (core business object)
│   ├── value_objects.go       # Value Objects (SessionName, ApiKey, etc.)
│   ├── repository.go          # Repository Interface
│   ├── service.go             # Domain Service Interface
│   ├── events.go              # Domain Events
│   └── errors.go              # Domain-specific Errors
└── common/                     # Shared Domain Concepts
    ├── value_objects.go       # Common Value Objects (ID, Timestamp)
    └── events.go              # Base Domain Event types
```

## 🎯 Bounded Context: Session Management

### Session Aggregate
**Session** é o **Aggregate Root** principal do sistema:

#### Core Responsibilities
- **Gerenciar ciclo de vida**: Criação → Conexão → Desconexão → Exclusão
- **Manter invariantes**: Estado consistente, transições válidas
- **Encapsular regras de negócio**: Validações, autenticação, configuração
- **Publicar eventos de domínio**: Mudanças de estado importantes

#### Aggregate Composition
- **Session Entity** (Aggregate Root)
  - Identity: SessionID (Value Object)
  - Core attributes: Name, Status, Timestamps
  - Behavior: Connect, Disconnect, Authenticate, Configure

- **Value Objects**
  - SessionName: Nome único da sessão
  - ApiKey: Chave de autenticação
  - WaJID: Identificador WhatsApp
  - QRCode: Código QR para pareamento
  - ProxyURL: Configuração de proxy
  - WebhookConfig: Configuração de webhooks

#### Domain Services
- **SessionDomainService**: Regras que envolvem múltiplas entidades
- **SessionIdentifierService**: Resolução e validação de identificadores

#### Domain Events
- SessionCreated, SessionConnected, SessionDisconnected
- SessionAuthenticated, SessionConfigurationChanged

## 🔧 Componentes do Session

### `entity.go`
- Entidade `Session` com comportamentos ricos
- Enum `Status` com validações
- Métodos de negócio (`CanConnect`, `IsAuthenticated`, etc.)

### `services.go`
- Implementação de regras de negócio complexas
- Validações que envolvem múltiplas entidades
- Lógica que não pertence a uma entidade específica

### `repository.go`
- Interface para persistência (implementada na infraestrutura)
- Apenas operações necessárias para o domínio

### `valueobjects.go`
- Objetos de valor imutáveis
- `SessionName`, `ProxyURL` (SessionID movido para shared/types)
- Validações intrínsecas aos value objects

### `identifier.go`
- Serviço de identificação e resolução de sessões
- Validação de formato UUID e nomes
- Normalização de identificadores

### `errors.go`
- Erros específicos do domínio
- Representam violações de regras de negócio
- Linguagem ubíqua nos nomes dos erros

## ✅ Benefícios desta Estrutura

1. **Testabilidade**: Domínio pode ser testado sem infraestrutura
2. **Clareza**: Fica óbvio o que é regra de negócio vs. detalhe técnico
3. **Manutenibilidade**: Mudanças de infraestrutura não afetam o domínio
4. **Reutilização**: Regras de negócio podem ser usadas em diferentes contextos
5. **Evolução**: Fácil adicionar novos contextos de domínio quando necessário

## 🚫 O que NÃO deve estar aqui

- Interfaces de APIs externas (meow, HTTP, etc.)
- Detalhes de persistência (SQL, NoSQL, etc.)
- Processamento de arquivos/mídia
- Configurações de infraestrutura
- DTOs de transferência de dados
- Validações puramente técnicas (formato de arquivo, etc.)

## 📖 Exemplo de Uso

```go
// Criando uma sessão
session := NewSession("uuid", "my-session")

// Usando domain service para validar regras de negócio
domainService := NewSessionDomainService()

if domainService.CanConnect(session) {
    session.SetStatus(StatusConnecting)
}

// Validando configuração completa
if err := domainService.ValidateSessionConfiguration(session); err != nil {
    // Tratar erro de domínio
}
```

Esta estrutura garante que o domínio permaneça puro e focado apenas nas regras de negócio essenciais.
