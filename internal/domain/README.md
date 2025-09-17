# Domain Layer

Esta camada contém apenas **conceitos de domínio puros**, seguindo rigorosamente os princípios de Domain-Driven Design (DDD).

## 🏛️ Princípios Seguidos

- **Independência de Infraestrutura**: Nenhuma dependência de frameworks, bancos de dados ou APIs externas
- **Regras de Negócio Puras**: Apenas lógica que reflete o conhecimento do domínio
- **Linguagem Ubíqua**: Termos e conceitos que fazem sentido para especialistas do domínio
- **Entidades Ricas**: Objetos com comportamento, não apenas estruturas de dados

## 📁 Estrutura

```
internal/domain/
└── session/                    # Contexto de Sessão meow
    ├── entity.go              # Entidade Session + Status enum
    ├── repository.go          # Interface para persistência
    ├── service.go             # Interface de serviços de domínio
    ├── domain_service.go      # Implementação de regras de negócio
    ├── valueobjects.go        # Value Objects (SessionID, SessionName, etc.)
    ├── validation.go          # Validações de domínio
    └── errors.go              # Erros específicos do domínio
```

## 🎯 Session - Único Contexto de Domínio

### Por que apenas Session?

- **Session** é a única entidade que possui **regras de negócio complexas**
- Tem **ciclo de vida próprio** (criação, conexão, desconexão, exclusão)
- Possui **invariantes de domínio** (transições de estado, validações)
- É o **conceito central** do negócio meow API

### O que foi removido e por quê?

#### ❌ Message
- **Não é entidade**: É um evento/comando, não tem ciclo de vida próprio
- **Sem regras de negócio**: Apenas transferência de dados
- **Movido para**: `internal/shared/types` como DTO

#### ❌ Media  
- **Não é conceito de domínio**: É detalhe técnico de processamento
- **Sem regras de negócio**: Apenas validações técnicas
- **Movido para**: `internal/infra/media` como serviço de infraestrutura

#### ❌ Webhooks
- **Configuração de infraestrutura**: Não é regra de negócio
- **Detalhe técnico**: Como notificar, não o que notificar
- **Movido para**: `internal/infra/webhooks` como serviço de infraestrutura

## 🔧 Componentes do Session

### `entity.go`
- Entidade `Session` com comportamentos ricos
- Enum `Status` com validações
- Métodos de negócio (`CanConnect`, `IsAuthenticated`, etc.)

### `domain_service.go`
- Implementação de regras de negócio complexas
- Validações que envolvem múltiplas entidades
- Lógica que não pertence a uma entidade específica

### `repository.go`
- Interface para persistência (implementada na infraestrutura)
- Apenas operações necessárias para o domínio

### `valueobjects.go`
- Objetos de valor imutáveis
- `SessionID`, `SessionName`, `ProxyURL`
- Validações intrínsecas aos value objects

### `validation.go`
- Funções puras de validação
- Regras de formato, tamanho, caracteres permitidos
- Validações de transição de estado

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
