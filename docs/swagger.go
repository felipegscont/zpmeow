// Package docs contains the Swagger documentation for the zpmeow API
package docs

// @title ZPMeow API
// @version 1.0
// @description ZPMeow is a WhatsApp API service that allows you to manage WhatsApp sessions, send messages, and interact with WhatsApp Web programmatically.
// @description
// @description ## Features
// @description - Create and manage multiple WhatsApp sessions
// @description - QR code authentication and phone pairing
// @description - Proxy support for sessions
// @description - Real-time session status monitoring
// @description
// @description ## Authentication
// @description Currently, the API does not require authentication. This may change in future versions.
// @description
// @description ## Session Management
// @description Sessions represent individual WhatsApp connections. Each session has:
// @description - Unique ID for identification
// @description - Name for easy reference
// @description - Status (disconnected, connecting, connected, etc.)
// @description - Optional QR code for mobile scanning
// @description - Optional proxy configuration
// @description
// @description ## Getting Started
// @description 1. Create a new session using POST /sessions/create
// @description 2. Connect the session using POST /sessions/{id}/connect
// @description 3. Get the QR code using GET /sessions/{id}/qr and scan it with your phone
// @description 4. Monitor the session status using GET /sessions/{id}/info
// @description
// @description Alternatively, you can use phone pairing:
// @description 1. Create a new session
// @description 2. Connect the session
// @description 3. Use POST /sessions/{id}/pair with your phone number
// @description 4. Enter the pairing code on your phone
//
// @contact.name ZPMeow Support
// @contact.url https://github.com/your-repo/zpmeow
// @contact.email support@zpmeow.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /
//
// @schemes http https
//
// @tag.name health
// @tag.description Health check endpoints
//
// @tag.name sessions
// @tag.description WhatsApp session management endpoints
