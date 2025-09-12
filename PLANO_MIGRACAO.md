# Plano de Migração - Reorganização do Projeto

## Visão Geral
Este plano detalha a reorganização da estrutura do projeto seguindo uma arquitetura em camadas mais clara e modular, baseado na análise detalhada dos arquivos existentes.

## Análise de Responsabilidades dos Arquivos

### Domain/Session - Análise Atual
- **dto.go**: Contém DTOs de API (Request/Response) - **VIOLAÇÃO**: DTOs não pertencem ao domain
- **entity.go**: Entidade Session com regras de negócio - **CORRETO**: Permanece no domain
- **service.go**: Use cases e orquestração - **VIOLAÇÃO**: Use cases devem estar na application layer
- **repository.go**: Interface do repositório - **VIOLAÇÃO**: Implementação deve estar na infra
- **validation.go**: Validações de domínio e infraestrutura misturadas - **VIOLAÇÃO**: Separar responsabilidades

### Types/Utils - Análise Atual
- **types/common.go**: Tipos básicos (Status, ID) - **CORRETO**: Shared types
- **types/wuzapi.go**: DTOs de API WhatsApp - **VIOLAÇÃO**: DTOs devem estar na application
- **utils/media.go**: Utilitários de mídia - **CORRETO**: Shared utilities
- **utils/response.go**: Utilitários HTTP - **VIOLAÇÃO**: Deve estar na infra/http
- **utils/validation.go**: Validações básicas - **CORRETO**: Shared utilities

### Infra/Database - Análise Atual
- **postgres.go**: Repository implementation + Models + Connection - **VIOLAÇÃO**: Múltiplas responsabilidades

### Infra/Meow - Análise Atual
- **client.go**: Cliente WhatsApp core - **CORRETO**: Core functionality
- **manager.go**: Gerenciamento de clientes - **CORRETO**: Core functionality
- **service.go**: Implementação do WhatsApp service - **CORRETO**: Service layer
- **event.go**: Manipulação de eventos - **CORRETO**: Event handling
- **messages.go**: Builders de mensagens - **CORRETO**: Service utilities
- **constants.go**: Constantes - **CORRETO**: Core constants
- **utils.go**: Utilitários específicos do meow - **CORRETO**: Adapter utilities
- **validation.go**: Validações específicas do meow - **CORRETO**: Adapter validation

### Infra/Webhook - Análise Atual
- **service.go**: Serviço de webhook com HTTP client - **CORRETO**: Infrastructure service
- **WebhookPayload**: Estrutura de dados do webhook - **VIOLAÇÃO**: Deveria ser DTO
- **SendWebhook methods**: Lógica de envio - **CORRETO**: Infrastructure service

### Infra/HTTP - Análise Completa

#### **Handlers - Violações Identificadas**
- **utils.go**:
  - **VIOLAÇÃO**: `SessionToDTOConverter` - lógica de conversão deveria estar na application layer
  - **VIOLAÇÃO**: `domainErrorMappings` - mapeamento de erros deveria estar em shared/errors
  - **CORRETO**: `ValidateSessionIDParam` - validação HTTP específica
- **session_handler.go**:
  - **VIOLAÇÃO**: Validação de phone numbers duplicada em múltiplos handlers
  - **VIOLAÇÃO**: Lógica de mapeamento de status inline
  - **CORRETO**: HTTP handling e routing
- **send_handler.go**:
  - **VIOLAÇÃO**: `handleSendResponse` - lógica de conversão duplicada
  - **VIOLAÇÃO**: Validação de mídia duplicada (`ValidateMediaSize`, `ValidateMediaType`)
  - **VIOLAÇÃO**: Switch case para tipos de mídia - deveria ser strategy pattern
- **chat_handler.go**:
  - **VIOLAÇÃO**: `mapPresenceState` - lógica de negócio no handler
  - **VIOLAÇÃO**: Validação de phone number duplicada
  - **VIOLAÇÃO**: Construção de response inline
- **group_handler.go**:
  - **VIOLAÇÃO**: Construção de participants array inline
  - **VIOLAÇÃO**: Validação de phone numbers duplicada
  - **VIOLAÇÃO**: Lógica de conversão de dados inline
- **webhook_handler.go**:
  - **VIOLAÇÃO**: `supportedEventTypes` hardcoded
  - **VIOLAÇÃO**: Validação de eventos duplicada
- **newsletter_handler.go**:
  - **VIOLAÇÃO**: Conversão de newsletter data inline
  - **VIOLAÇÃO**: Lógica de construção de response complexa

#### **Middleware - Análise**
- **logging.go**:
  - **CORRETO**: HTTP logging específico
  - **VIOLAÇÃO**: `HTTPLogEntry` - poderia ser shared type
  - **CORRETO**: Skip paths configuration
- **cors.go**: **CORRETO**: HTTP middleware específico

#### **Router - Análise**
- **routes.go**: **CORRETO**: HTTP routing configuration

### Infra/Meow - Análise Completa

#### **Core Components - Violações Identificadas**
- **constants.go**:
  - **VIOLAÇÃO**: Status constants duplicados com `types/common.go`
  - **VIOLAÇÃO**: Validation patterns deveriam estar em shared/validation
  - **CORRETO**: Timeouts e configurações específicas do meow
- **client.go**:
  - **VIOLAÇÃO**: Validação inline (`Validation.ValidateSessionID`)
  - **VIOLAÇÃO**: Error wrapping inline (`Error.WrapError`)
  - **CORRETO**: WhatsApp client management
- **manager.go**:
  - **VIOLAÇÃO**: Database queries inline (deveria usar repository)
  - **VIOLAÇÃO**: Validação duplicada
  - **CORRETO**: Client lifecycle management
- **service.go**:
  - **VIOLAÇÃO**: Database queries diretas (linha 1180+)
  - **VIOLAÇÃO**: `hasDeviceCredentials` - lógica de negócio complexa
  - **VIOLAÇÃO**: `waitForConnectionsToEstablish` - polling logic
  - **CORRETO**: WhatsApp service implementation

#### **Utilities e Validation - Violações**
- **utils.go**:
  - **VIOLAÇÃO**: `JIDUtils`, `ErrorUtils` - deveriam estar em shared/utils
  - **VIOLAÇÃO**: `IsRetryableError` - lógica de negócio
  - **CORRETO**: Utilities específicas do meow
- **validation.go**:
  - **VIOLAÇÃO**: Validações genéricas misturadas com específicas
  - **VIOLAÇÃO**: `ValidateMediaSize` duplicada com utils/media.go
  - **CORRETO**: Validações específicas do WhatsApp

#### **Event Handling - Análise**
- **event.go**:
  - **CORRETO**: Event handling específico do WhatsApp
  - **VIOLAÇÃO**: `sendWebhook` - poderia ser mais genérico
- **messages.go**:
  - **CORRETO**: Message builders específicos do WhatsApp

## Nova Estrutura Proposta (Corrigida)

```
internal/
├── application/
│   ├── dto/
│   │   ├── session/          # DTOs de sessão (de domain/session/dto.go)
│   │   │   ├── request.go    # CreateSessionRequest, PairSessionRequest, etc.
│   │   │   └── response.go   # SessionInfoResponse, QRCodeResponse, etc.
│   │   ├── whatsapp/         # DTOs do WhatsApp (de types/wuzapi.go)
│   │   │   ├── message.go    # SendTextRequest, SendMediaRequest, etc.
│   │   │   └── response.go   # SendResponse, etc.
│   │   └── webhook/          # DTOs de webhook (de types/wuzapi.go + webhook/service.go)
│   │       ├── request.go    # SetWebhookRequest, UpdateWebhookRequest
│   │       ├── response.go   # WebhookResponse
│   │       └── payload.go    # WebhookPayload (de webhook/service.go)
│   ├── services/
│   │   ├── conversion.go     # SessionToDTOConverter (de http/handler/utils.go)
│   │   ├── media_strategy.go # Strategy pattern para tipos de mídia
│   │   └── validation.go     # Validações de aplicação centralizadas
│   └── usecase/
│       └── session/
│           ├── session.go    # SessionService implementation (de domain/session/service.go)
│           └── session_test.go
├── domain/
│   └── session/
│       ├── entity.go         # Session entity (mantido)
│       ├── repository.go     # Repository interface (mantido como contrato)
│       ├── service.go        # WhatsApp service interface (mantido como contrato)
│       ├── errors.go         # Domain errors
│       └── validation.go     # Validações de domínio puras
├── shared/
│   ├── types/
│   │   ├── common.go         # Status, ID, Timestamp (de types/common.go)
│   │   └── http.go           # HTTPLogEntry (de middleware/logging.go)
│   ├── utils/
│   │   ├── media.go          # Utilitários de mídia (de utils/media.go)
│   │   ├── validation.go     # Validações básicas (de utils/validation.go)
│   │   ├── jid.go            # JIDUtils (de meow/utils.go)
│   │   └── conversion.go     # Utilities de conversão reutilizáveis
│   ├── errors/
│   │   ├── errors.go         # ErrorUtils (de meow/utils.go)
│   │   ├── mapping.go        # Domain error mappings (de http/handler/utils.go)
│   │   └── retry.go          # IsRetryableError logic
│   └── patterns/
│       ├── strategy.go       # Strategy pattern para media types
│       └── converter.go      # Base converter interfaces
├── infra/
│   ├── database/
│   │   ├── migrations/       # Mantido
│   │   ├── models/
│   │   │   └── session.go    # sessionModel + conversões (de postgres.go)
│   │   ├── repositories/
│   │   │   └── session.go    # PostgresSessionRepository (de postgres.go, sem queries inline)
│   │   └── connection.go     # Database connection (de postgres.go)
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── session.go    # SessionHandler (limpo, sem conversões inline)
│   │   │   ├── send.go       # SendHandler (usando strategy pattern)
│   │   │   ├── chat.go       # ChatHandler (sem lógica de negócio)
│   │   │   ├── group.go      # GroupHandler (sem conversões inline)
│   │   │   ├── webhook.go    # WebhookHandler (usando constantes centralizadas)
│   │   │   ├── user.go       # UserHandler (limpo)
│   │   │   ├── newsletter.go # NewsletterHandler (sem conversões inline)
│   │   │   └── health.go     # HealthHandler (mantido)
│   │   ├── middleware/
│   │   │   ├── cors.go       # Mantido
│   │   │   └── logging.go    # Mantido (HTTPLogEntry movido para shared)
│   │   ├── router/
│   │   │   └── routes.go     # Mantido
│   │   └── utils/
│   │       ├── response.go   # HTTP response utilities (de utils/response.go)
│   │       └── validation.go # Validações HTTP específicas
│   ├── webhook/
│   │   ├── service.go        # WebhookService (mantido, sem WebhookPayload)
│   │   ├── client.go         # HTTP client específico para webhooks
│   │   └── retry.go          # Lógica de retry separada
│   └── meow/
│       ├── core/
│       │   ├── client.go     # MeowClient (sem validações inline)
│       │   ├── manager.go    # ClientManager (sem database queries inline)
│       │   └── constants.go  # Constantes específicas do meow + Event constants
│       ├── service/
│       │   ├── service.go    # MeowServiceImpl (sem database queries inline)
│       │   └── messages.go   # MessageBuilder, MediaUploader
│       ├── event/
│       │   └── handler.go    # EventHandler (renomeado de event.go)
│       └── adapter/
│           ├── utils.go      # Utilitários específicos do meow apenas
│           └── validation.go # Validações específicas do WhatsApp apenas
```

## Fases da Migração (Revisadas)

### Fase 1: Criar Nova Estrutura de Diretórios
**Objetivo**: Criar todos os diretórios da nova estrutura

**Ações**:
- Criar `internal/application/dto/session/`
- Criar `internal/application/dto/whatsapp/`
- Criar `internal/application/usecase/session/`
- Criar `internal/shared/types/`
- Criar `internal/shared/utils/`
- Criar `internal/shared/errors/`
- Criar `internal/infra/database/models/`
- Criar `internal/infra/database/repositories/`
- Criar `internal/infra/http/utils/`
- Criar `internal/infra/meow/core/`
- Criar `internal/infra/meow/service/`
- Criar `internal/infra/meow/event/`
- Criar `internal/infra/meow/adapter/`

### Fase 2: Migrar Shared Components
**Objetivo**: Mover componentes compartilhados primeiro para evitar dependências circulares

**Ações**:
- Mover `internal/types/common.go` → `internal/shared/types/common.go`
- Mover `internal/utils/media.go` → `internal/shared/utils/media.go`
- Mover `internal/utils/validation.go` → `internal/shared/utils/validation.go`
- Mover `internal/utils/response.go` → `internal/infra/http/utils/response.go`
- Extrair `HTTPLogEntry` de `middleware/logging.go` → `internal/shared/types/http.go`
- Extrair `JIDUtils` de `meow/utils.go` → `internal/shared/utils/jid.go`
- Extrair `ErrorUtils` de `meow/utils.go` → `internal/shared/errors/errors.go`
- Extrair `domainErrorMappings` de `http/handler/utils.go` → `internal/shared/errors/mapping.go`
- Criar `internal/shared/patterns/strategy.go` para strategy pattern
- Criar `internal/shared/patterns/converter.go` para interfaces de conversão

### Fase 3: Migrar Application Layer - DTOs
**Objetivo**: Separar DTOs por responsabilidade

**Ações**:
- Separar `internal/domain/session/dto.go` em:
  - `internal/application/dto/session/request.go` (CreateSessionRequest, PairSessionRequest, etc.)
  - `internal/application/dto/session/response.go` (SessionInfoResponse, QRCodeResponse, etc.)
- Separar `internal/types/wuzapi.go` em:
  - `internal/application/dto/whatsapp/message.go` (SendTextRequest, SendMediaRequest, etc.)
  - `internal/application/dto/whatsapp/response.go` (SendResponse, etc.)
- Separar DTOs de webhook em:
  - `internal/application/dto/webhook/request.go` (SetWebhookRequest, UpdateWebhookRequest de types/wuzapi.go)
  - `internal/application/dto/webhook/response.go` (WebhookResponse de types/wuzapi.go)
  - `internal/application/dto/webhook/payload.go` (WebhookPayload de webhook/service.go)

### Fase 4: Migrar Domain Layer
**Objetivo**: Limpar domain mantendo apenas entidades e contratos

**Ações**:
- Manter `internal/domain/session/entity.go` (apenas limpar imports)
- Manter `internal/domain/session/repository.go` como interface
- Extrair interface WhatsAppService de `service.go` → `internal/domain/session/service.go`
- Separar validações de domínio puras em `internal/domain/session/validation.go`
- Criar `internal/domain/session/errors.go` com erros de domínio

### Fase 5: Migrar Application Layer - Services e Use Cases
**Objetivo**: Mover lógica de aplicação para camada correta

**Ações**:
- Mover `SessionToDTOConverter` de `http/handler/utils.go` → `internal/application/services/conversion.go`
- Criar `internal/application/services/media_strategy.go` com strategy pattern para tipos de mídia
- Criar `internal/application/services/validation.go` centralizando validações de aplicação
- Mover implementação SessionService de `internal/domain/session/service.go` → `internal/application/usecase/session/session.go`
- Mover `internal/domain/session/service_test.go` → `internal/application/usecase/session/session_test.go`
- Atualizar imports e dependências

### Fase 6: Reorganizar Database Infrastructure
**Objetivo**: Separar responsabilidades da infraestrutura de banco

**Ações**:
- Extrair `sessionModel` e funções de conversão de `postgres.go` → `internal/infra/database/models/session.go`
- Extrair `PostgresSessionRepository` de `postgres.go` → `internal/infra/database/repositories/session.go`
- Manter configuração de conexão em `internal/infra/database/connection.go`
- Remover `postgres.go` original

### Fase 7: Reorganizar Webhook Infrastructure
**Objetivo**: Organizar infraestrutura de webhook de forma modular

**Ações**:
- Refatorar `internal/infra/webhook/service.go`:
  - Remover `WebhookPayload` (já movido para application/dto)
  - Separar lógica de retry em `internal/infra/webhook/retry.go`
  - Criar `internal/infra/webhook/client.go` para HTTP client específico
- Mover `internal/infra/http/handler/webhook_handler.go` → `internal/infra/http/handlers/webhook.go`
- Atualizar `supportedEventTypes` para usar constantes de `internal/infra/meow/core/constants.go`

### Fase 8: Reorganizar Meow Infrastructure
**Objetivo**: Organizar infraestrutura WhatsApp de forma modular e limpar violações

**Ações**:
- Mover para `core/`:
  - `client.go` → `internal/infra/meow/core/client.go` (remover validações inline)
  - `manager.go` → `internal/infra/meow/core/manager.go` (remover database queries inline)
  - `constants.go` → `internal/infra/meow/core/constants.go` (incluir event constants, remover duplicações)
- Mover para `service/`:
  - `service.go` → `internal/infra/meow/service/service.go` (remover database queries inline)
  - `messages.go` → `internal/infra/meow/service/messages.go`
- Mover para `event/`:
  - `event.go` → `internal/infra/meow/event/handler.go`
- Mover para `adapter/`:
  - `utils.go` → `internal/infra/meow/adapter/utils.go` (apenas utilities específicas do meow)
  - `validation.go` → `internal/infra/meow/adapter/validation.go` (apenas validações específicas do WhatsApp)

### Fase 9: Refatorar HTTP Handlers
**Objetivo**: Limpar handlers HTTP removendo violações de responsabilidade

**Ações**:
- Refatorar todos os handlers para usar `application/services/conversion.go`
- Atualizar handlers para usar `application/services/validation.go`
- Implementar strategy pattern em `send_handler.go` para tipos de mídia
- Remover lógica de negócio inline dos handlers (mapPresenceState, etc.)
- Atualizar `webhook_handler.go` para usar constantes centralizadas
- Centralizar validações de phone number
- Remover construção de responses inline

### Fase 10: Atualizar Imports e Dependências
**Objetivo**: Corrigir todos os imports e dependências

**Ações**:
- Atualizar imports em todos os arquivos Go
- Corrigir dependências circulares
- Atualizar referências em testes
- Verificar interfaces e implementações
- Atualizar handlers HTTP para usar novos services
- Atualizar event handler para usar constantes centralizadas
- Atualizar meow components para usar shared utilities

### Fase 11: Testes e Validação Final
**Objetivo**: Garantir que tudo funciona após a migração

**Ações**:
- Executar `go mod tidy`
- Executar todos os testes: `go test ./...`
- Verificar build: `go build ./...`
- Testar funcionalidades principais (sessões, mensagens, webhooks)
- Validar que não há código duplicado
- Validar que validações estão centralizadas
- Validar que conversões usam services apropriados
- Testar strategy pattern para tipos de mídia
- Validar arquitetura final
- Testar integração webhook end-to-end

## Considerações Importantes

### Violações de Responsabilidade Identificadas
1. **DTOs no Domain**: DTOs de API estão misturados com entidades de domínio
2. **Use Cases no Domain**: Lógica de aplicação está na camada de domínio
3. **Repository Implementation no Domain**: Implementação de repositório deveria estar na infra
4. **Validações Misturadas**: Validações de domínio, aplicação e infraestrutura estão juntas
5. **HTTP Utils em Utils Gerais**: Utilitários HTTP específicos estão em utils compartilhados
6. **Múltiplas Responsabilidades**: `postgres.go` tem models, repository e connection

### Princípios Aplicados na Reorganização
1. **Separation of Concerns**: Cada camada tem responsabilidade específica
2. **Dependency Inversion**: Domain não depende de infraestrutura
3. **Single Responsibility**: Cada arquivo tem uma responsabilidade clara
4. **Clean Architecture**: Camadas bem definidas (Domain → Application → Infrastructure)

### Dependências e Ordem de Migração
- **Shared primeiro**: Para evitar dependências circulares
- **Domain depois**: Para estabelecer contratos
- **Application em seguida**: Para implementar use cases
- **Infrastructure por último**: Para implementar detalhes técnicos

### Testes
- Todos os testes devem ser movidos junto com seus arquivos correspondentes
- Mocks podem precisar ser atualizados para novas interfaces
- Testes de integração podem precisar de ajustes nos imports

### Compatibilidade
- Esta migração quebrará builds temporariamente
- **OBRIGATÓRIO**: Fazer em uma branch separada
- Commits incrementais por fase para facilitar rollback
- Testar cada fase antes de prosseguir

## Cronograma Estimado (Revisado com HTTP e Meow)
- **Fase 1**: 30 minutos (criar diretórios)
- **Fase 2**: 1.5 horas (shared components + extrações)
- **Fase 3**: 1.5 horas (DTOs incluindo webhooks)
- **Fase 4**: 45 minutos (domain cleanup)
- **Fase 5**: 1.5 horas (application services + use cases)
- **Fase 6**: 1.5 horas (database refactor)
- **Fase 7**: 1 hora (webhook refactor)
- **Fase 8**: 1.5 horas (meow refactor + cleanup)
- **Fase 9**: 2 horas (HTTP handlers refactor)
- **Fase 10**: 2-3 horas (imports e dependências)
- **Fase 11**: 2 horas (testes e validação completa)

**Total Estimado**: 12-14 horas

## Riscos e Mitigações
- **Risco**: Dependências circulares → **Mitigação**: Migrar shared components primeiro
- **Risco**: Quebra de interfaces → **Mitigação**: Manter interfaces no domain como contratos
- **Risco**: Imports quebrados → **Mitigação**: Usar IDE para refatoração automática quando possível
- **Risco**: Testes falhando → **Mitigação**: Atualizar testes incrementalmente por fase
- **Risco**: Funcionalidade perdida → **Mitigação**: Testes de integração após cada fase

## Definição de Constantes de Eventos WhatsApp

Durante a migração, as seguintes constantes de eventos devem ser definidas em `internal/infra/meow/core/constants.go`:

```go
// WhatsApp Event Types
const (
	// Messages and Communication
	EventMessage                    = "Message"
	EventUndecryptableMessage      = "UndecryptableMessage"
	EventReceipt                   = "Receipt"
	EventMediaRetry                = "MediaRetry"
	EventReadReceipt               = "ReadReceipt"

	// Groups and Contacts
	EventGroupInfo                 = "GroupInfo"
	EventJoinedGroup               = "JoinedGroup"
	EventPicture                   = "Picture"
	EventBlocklistChange           = "BlocklistChange"
	EventBlocklist                 = "Blocklist"

	// Connection and Session
	EventConnected                 = "Connected"
	EventDisconnected              = "Disconnected"
	EventConnectFailure            = "ConnectFailure"
	EventKeepAliveRestored         = "KeepAliveRestored"
	EventKeepAliveTimeout          = "KeepAliveTimeout"
	EventLoggedOut                 = "LoggedOut"
	EventClientOutdated            = "ClientOutdated"
	EventTemporaryBan              = "TemporaryBan"
	EventStreamError               = "StreamError"
	EventStreamReplaced            = "StreamReplaced"
	EventPairSuccess               = "PairSuccess"
	EventPairError                 = "PairError"
	EventQR                        = "QR"
	EventQRScannedWithoutMultidevice = "QRScannedWithoutMultidevice"

	// Privacy and Settings
	EventPrivacySettings           = "PrivacySettings"
	EventPushNameSetting           = "PushNameSetting"
	EventUserAbout                 = "UserAbout"

	// Synchronization and State
	EventAppState                  = "AppState"
	EventAppStateSyncComplete      = "AppStateSyncComplete"
	EventHistorySync               = "HistorySync"
	EventOfflineSyncCompleted      = "OfflineSyncCompleted"
	EventOfflineSyncPreview        = "OfflineSyncPreview"

	// Calls
	EventCallOffer                 = "CallOffer"
	EventCallAccept                = "CallAccept"
	EventCallTerminate             = "CallTerminate"
	EventCallOfferNotice           = "CallOfferNotice"
	EventCallRelayLatency          = "CallRelayLatency"

	// Presence and Activity
	EventPresence                  = "Presence"
	EventChatPresence              = "ChatPresence"

	// Identity
	EventIdentityChange            = "IdentityChange"

	// Errors
	EventCATRefreshError           = "CATRefreshError"

	// Newsletter (WhatsApp Channels)
	EventNewsletterJoin            = "NewsletterJoin"
	EventNewsletterLeave           = "NewsletterLeave"
	EventNewsletterMuteChange      = "NewsletterMuteChange"
	EventNewsletterLiveUpdate      = "NewsletterLiveUpdate"

	// Facebook/Meta Bridge
	EventFBMessage                 = "FBMessage"

	// Special - receives all events
	EventAll                       = "All"
)

// Event Categories for easier management
var (
	MessageEvents = []string{
		EventMessage,
		EventUndecryptableMessage,
		EventReceipt,
		EventMediaRetry,
		EventReadReceipt,
	}

	ConnectionEvents = []string{
		EventConnected,
		EventDisconnected,
		EventConnectFailure,
		EventKeepAliveRestored,
		EventKeepAliveTimeout,
		EventLoggedOut,
		EventClientOutdated,
		EventTemporaryBan,
		EventStreamError,
		EventStreamReplaced,
		EventPairSuccess,
		EventPairError,
		EventQR,
		EventQRScannedWithoutMultidevice,
	}

	GroupEvents = []string{
		EventGroupInfo,
		EventJoinedGroup,
		EventPicture,
		EventBlocklistChange,
		EventBlocklist,
	}

	CallEvents = []string{
		EventCallOffer,
		EventCallAccept,
		EventCallTerminate,
		EventCallOfferNotice,
		EventCallRelayLatency,
	}

	PresenceEvents = []string{
		EventPresence,
		EventChatPresence,
	}

	SyncEvents = []string{
		EventAppState,
		EventAppStateSyncComplete,
		EventHistorySync,
		EventOfflineSyncCompleted,
		EventOfflineSyncPreview,
	}

	NewsletterEvents = []string{
		EventNewsletterJoin,
		EventNewsletterLeave,
		EventNewsletterMuteChange,
		EventNewsletterLiveUpdate,
	}

	AllEvents = append(append(append(append(append(append(
		MessageEvents,
		ConnectionEvents...),
		GroupEvents...),
		CallEvents...),
		PresenceEvents...),
		SyncEvents...),
		NewsletterEvents...,
		EventPrivacySettings,
		EventPushNameSetting,
		EventUserAbout,
		EventIdentityChange,
		EventCATRefreshError,
		EventFBMessage,
	)
)

// Helper functions for event validation
func IsValidEvent(event string) bool {
	for _, validEvent := range AllEvents {
		if event == validEvent {
			return true
		}
	}
	return event == EventAll
}

func GetEventCategory(event string) string {
	eventCategories := map[string]string{
		// Message events
		EventMessage:              "message",
		EventUndecryptableMessage: "message",
		EventReceipt:              "message",
		EventMediaRetry:           "message",
		EventReadReceipt:          "message",

		// Connection events
		EventConnected:                      "connection",
		EventDisconnected:                   "connection",
		EventConnectFailure:                 "connection",
		EventKeepAliveRestored:              "connection",
		EventKeepAliveTimeout:               "connection",
		EventLoggedOut:                      "connection",
		EventClientOutdated:                 "connection",
		EventTemporaryBan:                   "connection",
		EventStreamError:                    "connection",
		EventStreamReplaced:                 "connection",
		EventPairSuccess:                    "connection",
		EventPairError:                      "connection",
		EventQR:                             "connection",
		EventQRScannedWithoutMultidevice:    "connection",

		// Group events
		EventGroupInfo:       "group",
		EventJoinedGroup:     "group",
		EventPicture:         "group",
		EventBlocklistChange: "group",
		EventBlocklist:       "group",

		// Call events
		EventCallOffer:        "call",
		EventCallAccept:       "call",
		EventCallTerminate:    "call",
		EventCallOfferNotice:  "call",
		EventCallRelayLatency: "call",

		// Presence events
		EventPresence:     "presence",
		EventChatPresence: "presence",

		// Sync events
		EventAppState:             "sync",
		EventAppStateSyncComplete: "sync",
		EventHistorySync:          "sync",
		EventOfflineSyncCompleted: "sync",
		EventOfflineSyncPreview:   "sync",

		// Newsletter events
		EventNewsletterJoin:       "newsletter",
		EventNewsletterLeave:      "newsletter",
		EventNewsletterMuteChange: "newsletter",
		EventNewsletterLiveUpdate: "newsletter",

		// Other events
		EventPrivacySettings:  "settings",
		EventPushNameSetting:  "settings",
		EventUserAbout:        "settings",
		EventIdentityChange:   "identity",
		EventCATRefreshError:  "error",
		EventFBMessage:        "bridge",
		EventAll:              "all",
	}

	if category, exists := eventCategories[event]; exists {
		return category
	}
	return "unknown"
}
```

## Análise Detalhada da Implementação de Webhooks

### Componentes Atuais de Webhook

#### 1. **WebhookService** (`internal/infra/webhook/service.go`)
**Responsabilidades Atuais**:
- HTTP client para envio de webhooks
- Estrutura `WebhookPayload` (violação - deveria ser DTO)
- Lógica de retry com backoff exponencial
- Envio síncrono e assíncrono

**Problemas Identificados**:
- `WebhookPayload` misturado com lógica de infraestrutura
- Lógica de retry acoplada ao service
- HTTP client não é reutilizável

#### 2. **WebhookHandler** (`internal/infra/http/handler/webhook_handler.go`)
**Responsabilidades Atuais**:
- Endpoints REST para configuração de webhooks
- Validação de eventos suportados
- CRUD de configurações de webhook

**Problemas Identificados**:
- Lista `supportedEventTypes` hardcoded
- Não usa constantes centralizadas de eventos
- Validação duplicada

#### 3. **DTOs de Webhook** (`internal/types/wuzapi.go`)
**Responsabilidades Atuais**:
- `SetWebhookRequest`, `UpdateWebhookRequest`, `WebhookResponse`

**Problemas Identificados**:
- DTOs misturados com outros tipos
- Não estão na camada de aplicação

#### 4. **Integração com Eventos** (`internal/infra/meow/event.go`)
**Responsabilidades Atuais**:
- `sendWebhook()` method no EventHandler
- Verificação de subscrição de eventos
- Envio assíncrono de webhooks

**Funciona Corretamente**: Boa integração entre eventos e webhooks

### Reorganização Proposta para Webhooks

#### **Application Layer** - DTOs
```
internal/application/dto/webhook/
├── request.go     # SetWebhookRequest, UpdateWebhookRequest
├── response.go    # WebhookResponse
└── payload.go     # WebhookPayload (estrutura enviada via HTTP)
```

#### **Infrastructure Layer** - Webhook Service
```
internal/infra/webhook/
├── service.go     # WebhookService (sem DTOs)
├── client.go      # HTTP client reutilizável
└── retry.go       # Lógica de retry separada
```

#### **Infrastructure Layer** - HTTP Handlers
```
internal/infra/http/handlers/
└── webhook.go     # WebhookHandler (usando constantes centralizadas)
```

### Melhorias na Reorganização

#### 1. **Separação de Responsabilidades**
- **DTOs**: Estruturas de dados na application layer
- **Service**: Lógica de negócio de webhook
- **Client**: HTTP client reutilizável
- **Retry**: Estratégias de retry configuráveis

#### 2. **Uso de Constantes Centralizadas**
- Handler usará constantes de `internal/infra/meow/core/constants.go`
- Validação consistente de eventos
- Facilita manutenção

#### 3. **Testabilidade Melhorada**
- HTTP client mockável
- Retry logic testável separadamente
- DTOs testáveis independentemente

#### 4. **Configurabilidade**
- Timeouts configuráveis
- Estratégias de retry configuráveis
- Headers customizáveis

### Fluxo de Webhook Após Reorganização

1. **Evento WhatsApp** → `EventHandler.HandleEvent()`
2. **Verificação** → `Session.IsEventSubscribed()`
3. **Criação Payload** → `application/dto/webhook/payload.go`
4. **Envio** → `infra/webhook/service.go`
5. **HTTP Request** → `infra/webhook/client.go`
6. **Retry se necessário** → `infra/webhook/retry.go`

## Benefícios Esperados
1. **Arquitetura Limpa**: Separação clara de responsabilidades
2. **Manutenibilidade**: Código mais fácil de manter e evoluir
3. **Testabilidade**: Melhor isolamento para testes unitários
4. **Reutilização**: Componentes shared podem ser reutilizados
5. **Escalabilidade**: Estrutura preparada para crescimento do projeto
6. **Eventos Organizados**: Constantes de eventos bem definidas e categorizadas
7. **Webhooks Modulares**: Infraestrutura de webhook bem organizada e testável
