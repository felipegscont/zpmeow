# WMeow - WhatsApp Integration Layer

Este pacote implementa a camada de integração com WhatsApp usando a biblioteca `whatsmeow`.

## Estrutura dos Arquivos

### 📁 Arquivos Principais

- **`service.go`** - Interface principal e implementação do serviço WhatsApp
- **`client.go`** - Cliente WhatsApp individual por sessão
- **`messages.go`** - Funções para envio de mensagens
- **`events.go`** - Processamento de eventos do WhatsApp
- **`session.go`** - Helpers para gerenciamento de sessões

## 🏗️ Arquitetura

### Interface Principal
```go
type Service interface {
    // Gerenciamento de Sessões
    StartClient(sessionID string) error
    StopClient(sessionID string) error
    LogoutClient(sessionID string) error
    
    // Status e Conexão
    IsClientConnected(sessionID string) bool
    GetClientStatus(sessionID string) types.Status
    GetQRCode(sessionID string) (string, error)
    PairPhone(sessionID, phoneNumber string) (string, error)
    
    // Mensagens
    SendTextMessage(ctx context.Context, sessionID, phone, text string) (*whatsmeow.SendResponse, error)
    SendImageMessage(ctx context.Context, sessionID, phone string, data []byte, caption, mimeType string) (*whatsmeow.SendResponse, error)
    // ... outras funções de mensagem
    
    // Grupos
    CreateGroup(ctx context.Context, sessionID, name string, participants []string) (*GroupInfo, error)
    ListGroups(ctx context.Context, sessionID string) ([]GroupInfo, error)
    // ... outras funções de grupo
    
    // Newsletters
    CreateNewsletter(ctx context.Context, sessionID string, params *whatsmeow.CreateNewsletterParams) (*waTypes.NewsletterMetadata, error)
    // ... outras funções de newsletter
}
```

### Implementação
```go
type serviceImpl struct {
    clients   map[string]*Client  // Clientes por sessão
    sessions  session.Repository  // Repositório de sessões
    logger    logging.Logger      // Logger
    container *sqlstore.Container // Container SQL
    waLogger  waLog.Logger        // Logger WhatsApp
    config    config.MeowConfigProvider // Configurações
    mu        sync.RWMutex        // Mutex para thread safety
}
```

## 📋 Responsabilidades por Arquivo

### `service.go`
- **Interface `Service`**: Define todas as operações disponíveis
- **Struct `serviceImpl`**: Implementação principal do serviço
- **Gerenciamento de clientes**: Criação, remoção e acesso aos clientes
- **Validações**: Validação de entrada e estado dos clientes
- **Operações de alto nível**: Grupos, newsletters, privacidade, etc.

### `client.go`
- **Struct `Client`**: Representa um cliente WhatsApp individual
- **Interface `EventHandler`**: Para processamento de eventos
- **Conexão**: Estabelecimento e manutenção da conexão
- **Autenticação**: QR Code e pareamento por telefone
- **Ciclo de vida**: Inicialização, reconexão, desconexão

### `messages.go`
- **Funções de envio**: Implementações específicas para cada tipo de mensagem
- **Validações**: Validação de entrada para mensagens
- **Upload de mídia**: Gerenciamento de upload de arquivos
- **Conversão de tipos**: Conversão entre tipos internos e whatsmeow

### `events.go`
- **Struct `EventProcessor`**: Processador central de eventos
- **Handlers específicos**: Para cada tipo de evento (mensagem, conexão, etc.)
- **Webhooks**: Envio de notificações via webhook
- **Atualização de estado**: Sincronização com o banco de dados

### `session.go`
- **Helpers de sessão**: Funções auxiliares para gerenciamento
- **QR Code**: Geração e exibição de códigos QR
- **Conexão**: Helpers para conexão segura
- **Status**: Gerenciamento de status das sessões
- **Device Store**: Gerenciamento do armazenamento de dispositivos

## 🔧 Padrões de Nomenclatura

### Structs
- `Service` - Interface principal
- `serviceImpl` - Implementação do serviço
- `Client` - Cliente WhatsApp
- `EventProcessor` - Processador de eventos
- `SessionHelper` - Helper de sessões

### Funções
- `NewService()` - Construtor do serviço
- `NewClient()` - Construtor do cliente
- `NewEventProcessor()` - Construtor do processador de eventos

### Métodos
- Seguem o padrão CamelCase
- Prefixos claros: `Get`, `Set`, `Send`, `Create`, `Update`, `Delete`
- Contexto como primeiro parâmetro quando aplicável

## 🚀 Uso

```go
// Criar o serviço
service := wmeow.NewService(container, waLogger, sessionRepo, config)

// Iniciar uma sessão
err := service.StartClient("session-id")

// Enviar mensagem
resp, err := service.SendTextMessage(ctx, "session-id", "5511999999999", "Olá!")

// Obter QR Code
qrCode, err := service.GetQRCode("session-id")
```

## 🔒 Thread Safety

- Todas as operações são thread-safe
- Uso de `sync.RWMutex` para proteção de mapas compartilhados
- Clientes individuais são isolados por sessão

## 📊 Logging

- Logger estruturado com níveis apropriados
- Contexto de sessão em todas as operações
- Logs de debug para troubleshooting
- Logs de erro com stack traces quando necessário
