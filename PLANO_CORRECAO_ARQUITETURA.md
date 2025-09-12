# Plano de Correção da Arquitetura - ZPMeow

## 🎯 Objetivo

Corrigir os problemas arquiteturais mais graves identificados na estrutura de camadas do projeto, garantindo a correta implementação dos princípios da Clean Architecture.

## 🚨 Problemas Críticos Identificados

### 1. **Handlers Dependendo Diretamente do Domain**
- **Problema**: Handlers acessam diretamente `domain.SessionService`
- **Impacto**: Viola a separação de camadas e dificulta testes
- **Prioridade**: 🔴 CRÍTICA

### 2. **Camada Application Incompleta**
- **Problema**: Use cases são apenas stubs sem implementação
- **Impacto**: Lógica de aplicação não está centralizada
- **Prioridade**: 🔴 CRÍTICA

### 3. **Falta de Validação e Conversão**
- **Problema**: Não há validação de entrada nem conversão entre DTOs
- **Impacto**: Dados inválidos podem chegar ao domain
- **Prioridade**: 🟡 ALTA

### 4. **Shared Types Muito Simples**
- **Problema**: Types básicos sem validação ou comportamento
- **Impacto**: Falta de type safety e validações
- **Prioridade**: 🟢 MÉDIA

## 📋 Plano de Execução

### **Fase 1: Corrigir Dependências dos Handlers** 
*Estimativa: 2-3 horas*

#### 1.1 Criar Interface de Application Service
```bash
# Arquivo: internal/application/services/interface.go
```

**Ações**:
- [ ] Criar `SessionApplicationService` interface
- [ ] Definir métodos que recebem/retornam DTOs
- [ ] Implementar conversões Domain ↔ DTO

#### 1.2 Implementar Application Service
```bash
# Arquivo: internal/application/services/session.go
```

**Ações**:
- [ ] Implementar `SessionApplicationServiceImpl`
- [ ] Injetar `domain.SessionService` como dependência
- [ ] Implementar conversões e validações

#### 1.3 Atualizar Handlers
```bash
# Arquivo: internal/infra/http/handlers/session_handler.go
```

**Ações**:
- [ ] Alterar dependência para `application.SessionApplicationService`
- [ ] Remover imports diretos do domain
- [ ] Usar DTOs em vez de entidades domain

### **Fase 2: Implementar Use Cases Completos**
*Estimativa: 4-5 horas*

#### 2.1 Implementar Session Use Cases
```bash
# Arquivo: internal/application/usecase/session/session.go
```

**Ações**:
- [ ] Implementar `CreateSession` com validações
- [ ] Implementar `GetSession` com tratamento de erros
- [ ] Implementar `GetAllSessions` com paginação
- [ ] Implementar operações de conexão/desconexão
- [ ] Implementar operações de QR Code e pairing

#### 2.2 Adicionar Validações de Negócio
```bash
# Arquivo: internal/application/validation/validator.go
```

**Ações**:
- [ ] Criar validador de DTOs de entrada
- [ ] Implementar regras de negócio específicas
- [ ] Integrar com validações do domain

### **Fase 3: Melhorar Sistema de Validação**
*Estimativa: 2-3 horas*

#### 3.1 Fortalecer Shared Types
```bash
# Arquivo: internal/shared/types/session.go
```

**Ações**:
- [ ] Criar `SessionID` com validação
- [ ] Criar `SessionName` com regras de negócio
- [ ] Implementar `Status` com transições válidas
- [ ] Adicionar métodos de validação

#### 3.2 Implementar Conversores
```bash
# Arquivo: internal/application/converters/converter.go
```

**Ações**:
- [ ] Criar `SessionConverter` interface
- [ ] Implementar conversões Domain → DTO
- [ ] Implementar conversões DTO → Domain
- [ ] Adicionar validações nas conversões

### **Fase 4: Atualizar Injeção de Dependências**
*Estimativa: 1-2 horas*

#### 4.1 Atualizar main.go
```bash
# Arquivo: cmd/server/main.go
```

**Ações**:
- [ ] Instanciar `SessionApplicationService`
- [ ] Injetar no lugar de `domain.SessionService` nos handlers
- [ ] Atualizar ordem de inicialização

#### 4.2 Atualizar Re-exports
```bash
# Arquivo: internal/application/application.go
```

**Ações**:
- [ ] Exportar novos services da application layer
- [ ] Remover exports desnecessários do domain
- [ ] Organizar imports

## 🔧 Implementação Detalhada

### **Estrutura de Arquivos Após Correções**

**Princípios da Nova Organização:**
- ✅ **Sem underscores** nos nomes de arquivos
- ✅ **Separação por pastas** para manter responsabilidades claras
- ✅ **Nomes descritivos** e simples
- ✅ **Agrupamento lógico** por funcionalidade

```
internal/
├── application/
│   ├── services/
│   │   ├── interface.go                # Interfaces dos services
│   │   └── session.go                  # Implementação session service
│   ├── validation/
│   │   └── validator.go                # Validações de entrada
│   ├── converters/
│   │   └── converter.go                # Conversões Domain ↔ DTO
│   ├── usecase/
│   │   └── session/                    # Pasta existente
│   │       └── session.go              # Use cases implementados
│   └── dto/ (existente)
├── domain/ (sem alterações)
├── infra/
│   └── http/handlers/
│       └── session.go                  # Atualizado para usar application
└── shared/
    └── types/
        ├── session.go                  # Types específicos
        └── common.go                   # Types gerais
```

### **Exemplo de Interface Application Service**

```go
// internal/application/services/interface.go
type SessionApplicationService interface {
    CreateSession(ctx context.Context, req *dto.CreateSessionRequest) (*dto.CreateSessionResponse, error)
    GetSession(ctx context.Context, req *dto.GetSessionRequest) (*dto.SessionInfoResponse, error)
    GetAllSessions(ctx context.Context) (*dto.SessionListResponse, error)
    DeleteSession(ctx context.Context, req *dto.DeleteSessionRequest) error
    ConnectSession(ctx context.Context, req *dto.ConnectSessionRequest) error
    DisconnectSession(ctx context.Context, req *dto.DisconnectSessionRequest) error
    GetQRCode(ctx context.Context, req *dto.GetQRCodeRequest) (*dto.QRCodeResponse, error)
    PairWithPhone(ctx context.Context, req *dto.PairSessionRequest) (*dto.PairSessionResponse, error)
}
```

## ✅ Critérios de Aceitação

### **Fase 1 - Completa quando:**
- [ ] Handlers não importam mais pacotes `domain`
- [ ] Handlers usam apenas DTOs da application layer
- [ ] Application service está funcionando como intermediário

### **Fase 2 - Completa quando:**
- [ ] Todos os use cases têm implementação real
- [ ] Validações de negócio estão implementadas
- [ ] Testes unitários passam

### **Fase 3 - Completa quando:**
- [ ] Types têm validação adequada
- [ ] Conversões são automáticas e seguras
- [ ] Erros são tratados adequadamente

### **Fase 4 - Completa quando:**
- [ ] Aplicação inicia sem erros
- [ ] Endpoints funcionam corretamente
- [ ] Arquitetura está limpa e testável

## 🧪 Estratégia de Testes

### **Durante a Implementação:**
1. **Testes Unitários**: Para cada service/converter implementado
2. **Testes de Integração**: Para handlers atualizados
3. **Testes de Contrato**: Entre camadas

### **Validação Final:**
1. **Smoke Tests**: Endpoints principais funcionando
2. **Architecture Tests**: Verificar dependências corretas
3. **Performance Tests**: Garantir que mudanças não degradaram performance

## 📊 Cronograma Sugerido

| Fase | Duração | Dependências | Responsável |
|------|---------|--------------|-------------|
| Fase 1 | 2-3h | - | Dev |
| Fase 2 | 4-5h | Fase 1 | Dev |
| Fase 3 | 2-3h | Fase 2 | Dev |
| Fase 4 | 1-2h | Fase 3 | Dev |
| **Total** | **9-13h** | - | - |

## 🚀 Próximos Passos

1. **Revisar e aprovar** este plano
2. **Criar branch** `fix/architecture-layers`
3. **Executar Fase 1** e fazer commit
4. **Testar** cada fase antes de prosseguir
5. **Fazer merge** após todas as fases completas

## 📝 Notas Importantes

- **Backward Compatibility**: Manter APIs existentes funcionando durante a transição
- **Testes**: Executar testes após cada fase
- **Documentação**: Atualizar documentação da arquitetura após conclusão
- **Code Review**: Revisar cada fase antes de prosseguir para a próxima
