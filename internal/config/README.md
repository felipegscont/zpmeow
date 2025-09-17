# 🔧 Configuration Module

Este módulo centraliza toda a configuração da aplicação meow, seguindo os princípios de Clean Architecture e fornecendo uma interface consistente para todas as configurações do sistema.

## 📁 Estrutura

```
internal/config/
├── config.go      # Estruturas principais e carregamento
├── interfaces.go  # Interfaces para abstração
├── defaults.go    # Configurações padrão e ambientes
└── README.md      # Esta documentação
```

## 🎯 Características

- ✅ **Centralizado**: Todas as configurações em um local
- ✅ **Tipado**: Configurações fortemente tipadas
- ✅ **Validado**: Validação automática de configurações
- ✅ **Flexível**: Suporte a múltiplos ambientes
- ✅ **Testável**: Interfaces para mocking em testes
- ✅ **Documentado**: Todas as variáveis documentadas

## 🏗️ Arquitetura

### Separação por Domínio

As configurações são organizadas por domínio funcional:

- **Database**: Configurações de banco de dados
- **Server**: Configurações do servidor HTTP
- **Auth**: Configurações de autenticação
- **Logging**: Configurações de logging
- **CORS**: Configurações de CORS
- **Webhook**: Configurações de webhooks
- **meow**: Configurações do cliente meow
- **Security**: Configurações de segurança

### Interfaces

Cada domínio possui uma interface específica que abstrai o acesso às configurações:

```go
type DatabaseConfigProvider interface {
    GetHost() string
    GetPort() string
    GetUser() string
    // ... outros métodos
}
```

## 📖 Uso

### Carregamento Básico

```go
import "meow/internal/config"

// Carregar configuração do ambiente
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Usar configuração
dbConfig := cfg.GetDatabase()
serverConfig := cfg.GetServer()
```

### Configurações por Ambiente

```go
// Desenvolvimento
cfg := config.DevelopmentConfig()

// Produção
cfg := config.ProductionConfig()

// Teste
cfg := config.TestConfig()
```

### Injeção de Dependência

```go
// Em vez de acessar configuração global
func NewDatabaseService(cfg config.DatabaseConfigProvider) *DatabaseService {
    return &DatabaseService{
        host: cfg.GetHost(),
        port: cfg.GetPort(),
    }
}
```

## 🌍 Variáveis de Ambiente

### ✅ Essenciais (obrigatórias no .env)

```bash
# Database
DB_HOST=localhost                    # Host do banco
DB_PORT=5432                        # Porta do banco
DB_USER=postgres                    # Usuário do banco
DB_PASSWORD=password                # Senha do banco
DB_NAME=meow                      # Nome do banco
DB_SSLMODE=disable                  # Modo SSL (require para produção)

# Server
SERVER_PORT=8080                    # Porta do servidor
GIN_MODE=debug                      # Modo (debug/release)

# Authentication
GLOBAL_API_KEY=your-secret-key      # Chave API global (obrigatória)

# Logging
LOG_LEVEL=info                      # Nível de log (debug/info/warn/error)
LOG_FORMAT=console                  # Formato (console/json)
```

### ⚙️ Avançadas (opcionais, com padrões sensatos)

```bash
# Database Pool (padrões: 25, 5, 5m)
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Server Timeouts (padrões: 30s, 30s, 120s)
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s

# Logging Avançado (padrão: arquivo desabilitado)
LOG_FILE_ENABLED=true
LOG_FILE_PATH=log/app.log
LOG_CONSOLE_COLOR=true
```

### 🔧 Configurações Automáticas (não precisam estar no .env)

O sistema possui padrões inteligentes para:

- **CORS**: Permite todas as origens em desenvolvimento
- **Webhook**: Timeout 30s, 3 tentativas, backoff exponencial
- **meow**: 3 tentativas, timeouts sensatos
- **Security**: Rate limiting desabilitado por padrão

Essas configurações podem ser customizadas via variáveis de ambiente se necessário, mas funcionam perfeitamente com os padrões.

## 🔄 Migração

### Antes (Configuração Espalhada)

```go
// Configuração hardcoded
config := cors.DefaultConfig()
config.AllowAllOrigins = true

// Configuração em múltiplos locais
timeout := 30 * time.Second
maxRetries := 3
```

### Depois (Configuração Centralizada)

```go
// Configuração centralizada e tipada
cfg, _ := config.LoadConfig()
corsConfig := cfg.GetCORS()
webhookConfig := cfg.GetWebhook()

// Uso através de interfaces
func NewService(cfg config.WebhookConfigProvider) *Service {
    return &Service{
        timeout: cfg.GetTimeout(),
        retries: cfg.GetMaxRetries(),
    }
}
```

## ✅ Benefícios

1. **Centralização**: Todas as configurações em um local
2. **Tipagem**: Configurações fortemente tipadas
3. **Validação**: Validação automática na inicialização
4. **Testabilidade**: Interfaces facilitam testes
5. **Flexibilidade**: Suporte a múltiplos ambientes
6. **Documentação**: Todas as variáveis documentadas
7. **Manutenibilidade**: Fácil de modificar e estender

## 🧪 Testes

```go
// Mock de configuração para testes
type MockConfig struct{}

func (m *MockConfig) GetTimeout() time.Duration {
    return 5 * time.Second
}

// Usar em testes
service := NewService(&MockConfig{})
```

Este módulo garante que todas as configurações sejam gerenciadas de forma consistente e centralizada, seguindo os princípios de Clean Architecture.
