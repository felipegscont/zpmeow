# Logger Package

Este pacote fornece uma camada de abstração de logging unificada para toda a aplicação, baseada na biblioteca **zerolog**.

## Características

- **Logger estruturado** com suporte a JSON e console
- **Cores configuráveis** para saída no console
- **Rotação automática** de arquivos de log
- **Múltiplos outputs** (console + arquivo)
- **Adapter para waLog** (whatsmeow)
- **Configuração centralizada** no módulo config

## Configuração

A configuração do logger agora está centralizada no módulo `internal/config`. Veja as variáveis de ambiente disponíveis:

```bash
# Nível de log (debug, info, warn, error, fatal)
LOG_LEVEL=info

# Formato do console (console, json)
LOG_FORMAT=console

# Habilitar cores no console
LOG_CONSOLE_COLOR=true

# Habilitar log em arquivo
LOG_FILE_ENABLED=true

# Caminho do arquivo de log
LOG_FILE_PATH=log/app.log

# Tamanho máximo antes da rotação (MB)
LOG_FILE_MAX_SIZE=100

# Número de backups a manter
LOG_FILE_MAX_BACKUPS=3

# Idade máxima dos logs (dias)
LOG_FILE_MAX_AGE=28

# Comprimir arquivos rotacionados
LOG_FILE_COMPRESS=true

# Formato do arquivo (json, console)
LOG_FILE_FORMAT=json
```

## Uso Básico

### Inicialização

```go
import (
    "zpmeow/internal/config"
    "zpmeow/internal/infra/logger"
)

// Usar configuração padrão
log := logger.Initialize(config.DefaultLoggerConfig())

// Ou carregar configuração do ambiente
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal("Failed to load config")
}
loggerConfig := cfg.GetLoggerConfig()
log := logger.Initialize(loggerConfig)

// Ou criar configuração customizada
customConfig := &config.LoggerConfig{
    Level:           "info",
    Format:          "console",
    ConsoleColor:    true,
    FileEnabled:     true,
    FilePath:        "log/app.log",
    FileMaxSize:     100,
    FileMaxBackups:  3,
    FileMaxAge:      28,
    FileCompress:    true,
    FileFormat:      "json",
}
log := logger.Initialize(customConfig)
```

### Logging Simples

```go
log := logger.GetLogger()

log.Info("Aplicação iniciada")
log.Infof("Usuário %s logado", userID)
log.Error("Erro ao conectar ao banco")
log.Errorf("Falha na operação: %v", err)
```

### Sub-loggers

```go
log := logger.GetLogger()
httpLogger := log.Sub("http")
dbLogger := log.Sub("database")

httpLogger.Info("Requisição recebida")
dbLogger.Error("Conexão perdida")
```

### Logging Estruturado

```go
log := logger.GetLogger()

log.With().
    Str("user_id", "123").
    Int("status_code", 200).
    Dur("duration", time.Second).
    Logger().
    Info("Requisição processada")
```

### Campos Múltiplos

```go
log := logger.GetLogger()

fields := map[string]interface{}{
    "user_id": "123",
    "action": "login",
    "ip": "192.168.1.1",
}

log.WithFields(fields).Info("Ação executada")
```

## Adapter para waLog

Para usar com bibliotecas que requerem `waLog.Logger` (como whatsmeow):

```go
import "zpmeow/internal/infra/logger"

// Criar adapter
waLogger := logger.GetWALogger("whatsapp")

// Usar com whatsmeow
client := whatsmeow.NewClient(deviceStore, waLogger)
```

## Estrutura de Arquivos

```
internal/
├── config/
│   └── config.go           # Configuração centralizada
└── infra/
    └── logger/
        ├── logger.go       # Logger principal consolidado
        ├── logger_test.go  # Testes
        └── log/
            └── app.log     # Arquivos de log
```

## Integração com Gin

O middleware de logging está automaticamente configurado para usar nosso logger:

```go
router.Use(middleware.Logger())
```

## Melhores Práticas

1. **Use a configuração centralizada** do módulo config
2. **Use sub-loggers** para diferentes módulos
3. **Prefira logging estruturado** em produção
4. **Configure níveis apropriados** por ambiente
5. **Use campos contextuais** para facilitar debugging
6. **Evite logging excessivo** em hot paths

## Migração

Se você estava usando `logger.NewConfigAdapter()`, agora use:

```go
// Antes
config := logger.NewConfigAdapter(...)

// Depois
cfg, _ := config.LoadConfig()
loggerConfig := cfg.GetLoggerConfig()
```
