# Architecture Documentation

## Overview

ZPMeow implements **Clean Architecture** principles for a WhatsApp session management API. This document defines the architectural layers, folder structure, and file responsibilities.

## Architectural Pattern: Clean Architecture

### Core Principles

1. **Dependency Inversion**: Outer layers depend on inner layers
2. **Separation of Concerns**: Each layer has a single responsibility  
3. **Testability**: Business logic isolated from external dependencies
4. **Independence**: Domain layer independent of frameworks and databases

### Layer Structure

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                       │
│                   (HTTP Handlers & Routes)                  │
├─────────────────────────────────────────────────────────────┤
│                   INFRASTRUCTURE LAYER                      │
│              (Database, WhatsApp, External APIs)            │
├─────────────────────────────────────────────────────────────┤
│                     DOMAIN LAYER                           │
│                (Business Logic & Entities)                  │
├─────────────────────────────────────────────────────────────┤
│                    SHARED COMPONENTS                        │
│                (Config, Types, Utils)                       │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
zpmeow/
├── cmd/
│   └── server/                  # ENTRY POINT
│       └── main.go              # Dependency injection and bootstrap
├── internal/
│   ├── domain/                  # DOMAIN LAYER
│   │   └── session/             # Bounded Context: WhatsApp Sessions
│   │       ├── entity.go        # Pure business entities
│   │       ├── dto.go           # Data Transfer Objects
│   │       ├── service.go       # Business logic interfaces
│   │       └── repository.go    # Persistence contracts
│   ├── infra/                   # INFRASTRUCTURE LAYER
│   │   ├── database/            # Persistence
│   │   │   └── postgres.go      # PostgreSQL implementation
│   │   ├── http/                # HTTP interface
│   │   │   ├── handler/         # HTTP controllers
│   │   │   ├── middleware/      # Middlewares
│   │   │   └── router/          # Route configuration
│   │   └── whatsapp/            # WhatsApp integration
│   │       └── service.go       # WhatsApp implementation
│   ├── config/                  # SHARED: Configuration
│   ├── types/                   # SHARED: Common types
│   └── utils/                   # SHARED: Utilities
├── docs/                        # Documentation
└── ref/                         # External references
```

## Layer Definitions

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Contains pure business logic and core application rules.

**Characteristics**:
- **Pure Entities**: No external dependencies (no database or JSON tags)
- **Business Rules**: Centralized domain logic
- **Interfaces**: Contracts defined in domain
- **Independence**: No dependencies on external layers

**File Structure**:
```
internal/domain/session/
├── entity.go        # Business entities with domain methods
├── dto.go           # Data Transfer Objects for API contracts
├── service.go       # Business logic interfaces and implementations
└── repository.go    # Persistence contracts (interfaces only)
```

**Key Components**:

- **entity.go**: Pure business entities
```go
type Session struct {
    ID        string
    Name      string
    DeviceJID string
    Status    types.Status  // No tags!
    QRCode    string
    ProxyURL  string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

- **service.go**: Business logic interfaces
```go
type SessionService interface {
    CreateSession(ctx context.Context, name string) (*Session, error)
    GetSession(ctx context.Context, id string) (*Session, error)
    ConnectSession(ctx context.Context, id string) error
}
```

- **repository.go**: Persistence contracts
```go
type SessionRepository interface {
    Save(ctx context.Context, session *Session) error
    FindByID(ctx context.Context, id string) (*Session, error)
    FindAll(ctx context.Context) ([]*Session, error)
}
```

### 2. Infrastructure Layer (`internal/infra/`)

**Purpose**: Technical implementations and adapters for external systems.

**Characteristics**:
- **Concrete Implementations**: PostgreSQL, HTTP, WhatsApp integrations
- **Adapters**: Convert between domain and infrastructure
- **Dependency Direction**: Depends on domain (dependency inversion)

**File Structure**:
```
internal/infra/
├── database/
│   └── postgres.go      # PostgreSQL repository implementations
├── http/
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middlewares
│   └── router/          # Route configuration
└── whatsapp/
    └── service.go       # WhatsApp service implementation
```

**Key Components**:

- **database/postgres.go**: Repository implementations
```go
type PostgresSessionRepository struct {
    db *sqlx.DB
}

// Database model with tags
type sessionModel struct {
    ID        string    `db:"id"`
    Name      string    `db:"name"`
    DeviceJID string    `db:"device_jid"`
    Status    string    `db:"status"`
    QRCode    string    `db:"qr_code"`
    ProxyURL  string    `db:"proxy_url"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

- **http/handler/**: HTTP controllers
```go
type SessionHandler struct {
    sessionService session.SessionService
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
    // HTTP logic, validation, serialization
}
```

- **whatsapp/service.go**: External service integration
```go
type WhatsAppServiceImpl struct {
    db        *sqlx.DB
    container *sqlstore.Container
    clients   map[string]*whatsmeow.Client
}
```

### 3. Shared Components (`internal/config/`, `internal/types/`, `internal/utils/`)

**Purpose**: Reusable code across all layers.

**File Structure**:
```
internal/
├── config/
│   └── config.go        # Application configuration
├── types/
│   └── common.go        # Shared types and enums
└── utils/
    ├── validation.go    # Validation utilities
    └── response.go      # HTTP response utilities
```

**Components**:
- **config/**: Centralized configuration management
- **types/**: Common types and enums (Status, ID, etc.)
- **utils/**: Helper functions (validation, response formatting)

### 4. Entry Point (`cmd/server/`)

**Purpose**: Application bootstrap and dependency injection.

**File Structure**:
```
cmd/server/
└── main.go              # Application entry point
```

**Responsibilities**:
- Initialize database connections
- Configure dependencies
- Wire up repositories, services, and handlers
- Start HTTP server

## Data Flow

```
HTTP Request → Handler → Service → Repository → Database
                ↓           ↓          ↓
            Validation   Business   Persistence
                ↓        Logic        ↓
            DTO/JSON   ←  Entity  ←  SQL Model
```

## Design Patterns

### Repository Pattern
- **Abstraction**: Persistence layer abstracted through interfaces
- **Contracts**: Defined in domain layer
- **Implementations**: Located in infrastructure layer

### Service Pattern
- **Business Logic**: Centralized in service layer
- **Interfaces**: Well-defined contracts
- **Reusability**: Shared across multiple handlers

### Dependency Injection
- **Inversion of Control**: Dependencies injected at startup
- **Configuration**: Centralized in main.go
- **Flexibility**: Easy to swap implementations

## Entity-Database Mapping

```go
// Domain Entity (no tags)
type Session struct {
    ID        string        // → sessions.id
    Name      string        // → sessions.name  
    DeviceJID string        // → sessions.device_jid
    Status    types.Status  // → sessions.status
    QRCode    string        // → sessions.qr_code
    ProxyURL  string        // → sessions.proxy_url
    CreatedAt time.Time     // → sessions.created_at
    UpdatedAt time.Time     // → sessions.updated_at
}

// Database Model (with tags)
type sessionModel struct {
    ID        string    `db:"id"`
    Name      string    `db:"name"`
    DeviceJID string    `db:"device_jid"`
    Status    string    `db:"status"`
    QRCode    string    `db:"qr_code"`
    ProxyURL  string    `db:"proxy_url"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

## Benefits

- **Testability**: Business logic isolated and easily testable
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations
- **Scalability**: Structure supports growth and new features
