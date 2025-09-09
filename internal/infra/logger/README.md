# Logger Package

Este pacote fornece uma camada de abstração de logging unificada para toda a aplicação, baseada na biblioteca **zerolog**.

## Características

- **Logger estruturado** com suporte a JSON e console
- **Cores configuráveis** para saída no console
- **Rotação automática** de arquivos de log
- **Múltiplos outputs** (console + arquivo)
- **Adapter para waLog** (whatsmeow)
- **Configuração via .env**

## Configuração

### Variáveis de Ambiente

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
import "zpmeow/internal/infra/logger"

// Usar configuração padrão
log := logger.Initialize(logger.DefaultConfig())

// Ou criar configuração customizada
config := logger.NewConfigAdapter(
    "info",           // level
    "console",        // format
    "log/app.log",    // filePath
    "json",           // fileFormat
    true,             // consoleColor
    true,             // fileEnabled
    true,             // fileCompress
    100,              // fileMaxSize
    3,                // fileMaxBackups
    28,               // fileMaxAge
)
log := logger.Initialize(config)
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
log/
├── app.log          # Log atual
├── app.log.1        # Backup 1
├── app.log.2        # Backup 2
├── app.log.3.gz     # Backup comprimido
└── README.md        # Documentação
```

## Níveis de Log

- **DEBUG**: Informações detalhadas para debugging
- **INFO**: Informações gerais sobre funcionamento
- **WARN**: Avisos sobre situações que podem precisar atenção
- **ERROR**: Erros que não impedem o funcionamento
- **FATAL**: Erros críticos que causam encerramento

## Formatos de Saída

### Console (Desenvolvimento)
```
2024-01-15 10:30:45 INF Starting server module=app
2024-01-15 10:30:45 DBG Database connected module=database
```

### JSON (Produção)
```json
{"level":"info","module":"app","time":"2024-01-15T10:30:45Z","message":"Starting server"}
{"level":"debug","module":"database","time":"2024-01-15T10:30:45Z","message":"Database connected"}
```

## Integração com Gin

O middleware de logging está automaticamente configurado para usar nosso logger:

```go
router.Use(middleware.Logger())
```

## Melhores Práticas

1. **Use sub-loggers** para diferentes módulos
2. **Prefira logging estruturado** em produção
3. **Configure níveis apropriados** por ambiente
4. **Use campos contextuais** para facilitar debugging
5. **Evite logging excessivo** em hot paths
