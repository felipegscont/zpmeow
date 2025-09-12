package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "zpmeow API Support",
            "url": "https://github.com/your-username/zpmeow",
            "email": "support@zpmeow.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ping": {
            "get": {
                "description": "Returns a simple pong message to verify the API is running",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.PingResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/delete": {
            "post": {
                "description": "Delete a message from a chat",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Delete a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Delete request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatDeleteRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/download/audio": {
            "post": {
                "description": "Download an audio attachment from a WhatsApp message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Download an audio from a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Download request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatDownloadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/download/document": {
            "post": {
                "description": "Download a document attachment from a WhatsApp message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Download a document from a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Download request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatDownloadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/download/image": {
            "post": {
                "description": "Download an image attachment from a WhatsApp message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Download an image from a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Download request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatDownloadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/download/video": {
            "post": {
                "description": "Download a video attachment from a WhatsApp message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Download a video from a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Download request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatDownloadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/edit": {
            "post": {
                "description": "Edit a text message in a chat",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Edit a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Edit request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatEditRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/markread": {
            "post": {
                "description": "Mark one or more messages as read in a chat",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Mark messages as read",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Mark read request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatMarkReadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/presence": {
            "post": {
                "description": "Set presence status in a chat (typing, recording, paused)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Set chat presence",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Chat presence request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatPresenceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/chat/react": {
            "post": {
                "description": "Add an emoji reaction to a message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "React to a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "React request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.ChatReactRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/announce/set": {
            "post": {
                "description": "Set whether only admins can send messages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Set group announce mode",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Set announce request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupSetAnnounceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/create": {
            "post": {
                "description": "Create a new WhatsApp group with specified participants",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Create a new WhatsApp group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Group creation request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/ephemeral/set": {
            "post": {
                "description": "Set the duration for disappearing messages in a group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Set group ephemeral messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Set ephemeral request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupSetEphemeralRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/info": {
            "get": {
                "description": "Get detailed information about a specific group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Get group information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Group JID",
                        "name": "groupJid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/inviteinfo": {
            "post": {
                "description": "Get information about a group invite without joining",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Get invite information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Invite info request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupInviteInfoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/invitelink": {
            "get": {
                "description": "Get the invite link for a group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Get group invite link",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Group JID",
                        "name": "groupJid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/join": {
            "post": {
                "description": "Join a WhatsApp group using an invite code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Join a group via invite link",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Group join request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupJoinRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/leave": {
            "post": {
                "description": "Leave a WhatsApp group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Leave a group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Group leave request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupLeaveRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/list": {
            "get": {
                "description": "Get a list of all groups the session is part of",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "List all groups",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/locked/set": {
            "post": {
                "description": "Set whether only admins can edit group info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Set group locked mode",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Set locked request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupSetLockedRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/name/set": {
            "post": {
                "description": "Update the name of a group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Set group name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Set name request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupSetNameRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/participants/update": {
            "post": {
                "description": "Add, remove, promote, or demote group participants",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Update group participants",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update participants request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupUpdateParticipantsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/photo/remove": {
            "post": {
                "description": "Remove the photo of a group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Remove group photo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Remove photo request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupRemovePhotoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/photo/set": {
            "post": {
                "description": "Update the photo of a group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Set group photo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Set photo request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupSetPhotoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/group/topic/set": {
            "post": {
                "description": "Update the topic/description of a group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Set group topic/description",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Set topic request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.GroupSetTopicRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/audio": {
            "post": {
                "description": "Send an audio message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send an audio message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Audio message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendAudioRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/buttons": {
            "post": {
                "description": "Send an interactive buttons message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send an interactive buttons message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Buttons message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendButtonsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/contact": {
            "post": {
                "description": "Send a contact message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a contact message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Contact message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendContactRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/document": {
            "post": {
                "description": "Send a document message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a document message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Document message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendDocumentRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/image": {
            "post": {
                "description": "Send an image message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send an image message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Image message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendImageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/list": {
            "post": {
                "description": "Send an interactive list message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send an interactive list message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "List message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendListRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/location": {
            "post": {
                "description": "Send a location message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a location message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Location message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendLocationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/media": {
            "post": {
                "description": "Send any type of media message (image, audio, document, video) using a unified endpoint",
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a media message (unified endpoint)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Media message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendMediaRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/poll": {
            "post": {
                "description": "Send a poll message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a poll message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Poll message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendPollRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/sticker": {
            "post": {
                "description": "Send a sticker message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a sticker message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Sticker message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendStickerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/text": {
            "post": {
                "description": "Send a text message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a text message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Text message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendTextRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/{sessionId}/send/video": {
            "post": {
                "description": "Send a video message to a WhatsApp contact",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "send"
                ],
                "summary": "Send a video message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Video message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.SendVideoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.SendResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/create": {
            "post": {
                "description": "Creates a new WhatsApp session with the provided name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Create a new WhatsApp session",
                "parameters": [
                    {
                        "description": "Session creation request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/session.CreateSessionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/session.CreateSessionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/list": {
            "get": {
                "description": "Retrieves a list of all WhatsApp sessions in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "List all WhatsApp sessions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.SessionListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/connect": {
            "post": {
                "description": "Starts the connection process for a WhatsApp session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Connect a session to WhatsApp",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/session.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/delete": {
            "delete": {
                "description": "Deletes a WhatsApp session and logs out the client if active",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Delete a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Session deleted successfully"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/info": {
            "get": {
                "description": "Retrieves detailed information about a specific WhatsApp session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Get session information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.SessionInfoResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/logout": {
            "post": {
                "description": "Logs out a WhatsApp session and disconnects the client",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Logout a session from WhatsApp",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/pair": {
            "post": {
                "description": "Pairs a WhatsApp session with a phone number using pairing code method",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Pair session with phone number",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Phone number to pair",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/session.PairSessionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.PairSessionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/proxy/find": {
            "get": {
                "description": "Retrieves the current proxy configuration for a WhatsApp session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Get proxy configuration for session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.ProxyResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/proxy/set": {
            "post": {
                "description": "Sets or updates the proxy configuration for a WhatsApp session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Set proxy for session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Proxy configuration",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/session.ProxyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.ProxyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/qr": {
            "get": {
                "description": "Retrieves the QR code for a WhatsApp session to scan with mobile device",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sessions"
                ],
                "summary": "Get QR code for session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/session.QRCodeResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.PingResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "pong"
                }
            }
        },
        "session.CreateSessionRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "My WhatsApp Session"
                }
            }
        },
        "session.CreateSessionResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2023-01-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "name": {
                    "type": "string",
                    "example": "My WhatsApp Session"
                },
                "status": {
                    "type": "string",
                    "example": "disconnected"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-01-01T00:00:00Z"
                }
            }
        },
        "session.MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Operation completed successfully"
                }
            }
        },
        "session.PairSessionRequest": {
            "type": "object",
            "required": [
                "phone_number"
            ],
            "properties": {
                "phone_number": {
                    "type": "string",
                    "example": "5511999999999"
                }
            }
        },
        "session.PairSessionResponse": {
            "type": "object",
            "properties": {
                "pairing_code": {
                    "type": "string",
                    "example": "ABCD-EFGH"
                }
            }
        },
        "session.ProxyRequest": {
            "type": "object",
            "required": [
                "proxy_url"
            ],
            "properties": {
                "proxy_url": {
                    "type": "string",
                    "example": "http://proxy.example.com:8080"
                }
            }
        },
        "session.ProxyResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Proxy set successfully"
                },
                "proxy_url": {
                    "type": "string",
                    "example": "http://proxy.example.com:8080"
                }
            }
        },
        "session.QRCodeResponse": {
            "type": "object",
            "properties": {
                "qr_code": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
                },
                "status": {
                    "type": "string",
                    "example": "waiting_for_scan"
                }
            }
        },
        "session.SessionInfoResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2023-01-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "name": {
                    "type": "string",
                    "example": "My WhatsApp Session"
                },
                "proxy_url": {
                    "type": "string",
                    "example": "http://proxy.example.com:8080"
                },
                "qr_code": {
                    "type": "string",
                    "example": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
                },
                "status": {
                    "type": "string",
                    "example": "disconnected"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-01-01T00:00:00Z"
                },
                "whatsapp_jid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                }
            }
        },
        "session.SessionListResponse": {
            "type": "object",
            "properties": {
                "sessions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/session.SessionInfoResponse"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 5
                }
            }
        },
        "types.Button": {
            "type": "object",
            "required": [
                "buttonId",
                "buttonText"
            ],
            "properties": {
                "buttonId": {
                    "type": "string",
                    "example": "btn_1"
                },
                "buttonText": {
                    "$ref": "#/definitions/types.ButtonText"
                },
                "type": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "types.ButtonText": {
            "type": "object",
            "required": [
                "displayText"
            ],
            "properties": {
                "displayText": {
                    "type": "string",
                    "example": "Option 1"
                }
            }
        },
        "types.ChatDeleteRequest": {
            "type": "object",
            "required": [
                "messageId",
                "phone"
            ],
            "properties": {
                "forEveryone": {
                    "type": "boolean",
                    "example": false
                },
                "messageId": {
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.ChatDownloadRequest": {
            "type": "object",
            "required": [
                "messageId"
            ],
            "properties": {
                "messageId": {
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.ChatEditRequest": {
            "type": "object",
            "required": [
                "messageId",
                "newText",
                "phone"
            ],
            "properties": {
                "messageId": {
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                },
                "newText": {
                    "type": "string",
                    "example": "Edited message"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.ChatMarkReadRequest": {
            "type": "object",
            "required": [
                "messageIds",
                "phone"
            ],
            "properties": {
                "messageIds": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "3EB0C431C26A1916E07A",
                        "3EB0C431C26A1916E07B"
                    ]
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.ChatPresenceRequest": {
            "type": "object",
            "required": [
                "phone",
                "state"
            ],
            "properties": {
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "state": {
                    "description": "typing, recording, paused",
                    "type": "string",
                    "example": "typing"
                }
            }
        },
        "types.ChatReactRequest": {
            "type": "object",
            "required": [
                "emoji",
                "messageId",
                "phone"
            ],
            "properties": {
                "emoji": {
                    "type": "string",
                    "example": "👍"
                },
                "messageId": {
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.Contact": {
            "type": "object",
            "required": [
                "displayName",
                "vcard"
            ],
            "properties": {
                "displayName": {
                    "type": "string",
                    "example": "John Doe"
                },
                "vcard": {
                    "type": "string",
                    "example": "BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"
                }
            }
        },
        "types.ContextInfo": {
            "type": "object",
            "properties": {
                "participant": {
                    "type": "string",
                    "example": "+5511888888888@s.whatsapp.net"
                },
                "quotedMessage": {},
                "stanzaId": {
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                }
            }
        },
        "types.GroupCreateRequest": {
            "type": "object",
            "required": [
                "name",
                "participants"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "My Group"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "+5511999999999",
                        "+5511888888888"
                    ]
                }
            }
        },
        "types.GroupInviteInfoRequest": {
            "type": "object",
            "required": [
                "inviteCode"
            ],
            "properties": {
                "inviteCode": {
                    "type": "string",
                    "example": "CjQKOAokMjU5NzE4NzAtNzBiYy00"
                }
            }
        },
        "types.GroupJoinRequest": {
            "type": "object",
            "required": [
                "inviteCode"
            ],
            "properties": {
                "inviteCode": {
                    "type": "string",
                    "example": "CjQKOAokMjU5NzE4NzAtNzBiYy00"
                }
            }
        },
        "types.GroupLeaveRequest": {
            "type": "object",
            "required": [
                "groupJid"
            ],
            "properties": {
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                }
            }
        },
        "types.GroupRemovePhotoRequest": {
            "type": "object",
            "required": [
                "groupJid"
            ],
            "properties": {
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                }
            }
        },
        "types.GroupSetAnnounceRequest": {
            "type": "object",
            "required": [
                "announce",
                "groupJid"
            ],
            "properties": {
                "announce": {
                    "type": "boolean",
                    "example": true
                },
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                }
            }
        },
        "types.GroupSetEphemeralRequest": {
            "type": "object",
            "required": [
                "duration",
                "groupJid"
            ],
            "properties": {
                "duration": {
                    "description": "seconds",
                    "type": "integer",
                    "example": 86400
                },
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                }
            }
        },
        "types.GroupSetLockedRequest": {
            "type": "object",
            "required": [
                "groupJid",
                "locked"
            ],
            "properties": {
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                },
                "locked": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "types.GroupSetNameRequest": {
            "type": "object",
            "required": [
                "groupJid",
                "name"
            ],
            "properties": {
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                },
                "name": {
                    "type": "string",
                    "example": "New Group Name"
                }
            }
        },
        "types.GroupSetPhotoRequest": {
            "type": "object",
            "required": [
                "groupJid",
                "image"
            ],
            "properties": {
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                },
                "image": {
                    "type": "string",
                    "example": "data:image/jpeg;base64,/9j/4AAQ..."
                }
            }
        },
        "types.GroupSetTopicRequest": {
            "type": "object",
            "required": [
                "groupJid",
                "topic"
            ],
            "properties": {
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                },
                "topic": {
                    "type": "string",
                    "example": "New group description"
                }
            }
        },
        "types.GroupUpdateParticipantsRequest": {
            "type": "object",
            "required": [
                "action",
                "groupJid",
                "participants"
            ],
            "properties": {
                "action": {
                    "description": "add, remove, promote, demote",
                    "type": "string",
                    "example": "add"
                },
                "groupJid": {
                    "type": "string",
                    "example": "120363025246125486@g.us"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "+5511999999999",
                        "+5511888888888"
                    ]
                }
            }
        },
        "types.Row": {
            "type": "object",
            "required": [
                "rowId",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Row description"
                },
                "rowId": {
                    "type": "string",
                    "example": "row_1"
                },
                "title": {
                    "type": "string",
                    "example": "Row 1"
                }
            }
        },
        "types.Section": {
            "type": "object",
            "required": [
                "rows",
                "title"
            ],
            "properties": {
                "rows": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Row"
                    }
                },
                "title": {
                    "type": "string",
                    "example": "Section 1"
                }
            }
        },
        "types.SendAudioRequest": {
            "type": "object",
            "required": [
                "audio",
                "phone"
            ],
            "properties": {
                "audio": {
                    "type": "string",
                    "example": "data:audio/ogg;base64,T2dnU..."
                },
                "caption": {
                    "type": "string",
                    "example": "Audio caption"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendButtonsRequest": {
            "type": "object",
            "required": [
                "buttons",
                "phone",
                "text"
            ],
            "properties": {
                "buttons": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Button"
                    }
                },
                "footer": {
                    "type": "string",
                    "example": "Footer text"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "text": {
                    "type": "string",
                    "example": "Choose an option:"
                }
            }
        },
        "types.SendContactRequest": {
            "type": "object",
            "required": [
                "contact",
                "phone"
            ],
            "properties": {
                "contact": {
                    "$ref": "#/definitions/types.Contact"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendDocumentRequest": {
            "type": "object",
            "required": [
                "document",
                "phone"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Document caption"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "document": {
                    "type": "string",
                    "example": "data:application/pdf;base64,JVBERi0x..."
                },
                "filename": {
                    "type": "string",
                    "example": "document.pdf"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendImageRequest": {
            "type": "object",
            "required": [
                "image",
                "phone"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Image caption"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "image": {
                    "type": "string",
                    "example": "data:image/jpeg;base64,/9j/4AAQ..."
                },
                "mimeType": {
                    "type": "string",
                    "example": "image/jpeg"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendListRequest": {
            "type": "object",
            "required": [
                "buttonText",
                "phone",
                "sections",
                "text"
            ],
            "properties": {
                "buttonText": {
                    "type": "string",
                    "example": "Select Option"
                },
                "footer": {
                    "type": "string",
                    "example": "Footer text"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "sections": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Section"
                    }
                },
                "text": {
                    "type": "string",
                    "example": "Choose from the list:"
                }
            }
        },
        "types.SendLocationRequest": {
            "type": "object",
            "required": [
                "latitude",
                "longitude",
                "phone"
            ],
            "properties": {
                "address": {
                    "type": "string",
                    "example": "São Paulo, SP, Brazil"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "latitude": {
                    "type": "number",
                    "example": -23.5505
                },
                "longitude": {
                    "type": "number",
                    "example": -46.6333
                },
                "name": {
                    "type": "string",
                    "example": "São Paulo"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendMediaRequest": {
            "type": "object",
            "required": [
                "mediaType",
                "phone"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Media caption"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "filename": {
                    "type": "string",
                    "example": "document.pdf"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "media": {
                    "type": "string",
                    "example": "data:image/jpeg;base64,/9j/4AAQ..."
                },
                "mediaType": {
                    "type": "string",
                    "enum": [
                        "image",
                        "audio",
                        "document",
                        "video"
                    ],
                    "example": "image"
                },
                "mimeType": {
                    "type": "string",
                    "example": "image/jpeg"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendPollRequest": {
            "type": "object",
            "required": [
                "name",
                "options",
                "phone"
            ],
            "properties": {
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "name": {
                    "type": "string",
                    "example": "What's your favorite color?"
                },
                "options": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Red",
                        "Blue",
                        "Green"
                    ]
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "selectableCount": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "types.SendResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "The ID of the sent message",
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                },
                "messageId": {
                    "type": "string",
                    "example": "3EB0C431C26A1916E07A"
                },
                "sender": {
                    "description": "The identity the message was sent with (LID or PN)",
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "serverId": {
                    "description": "The server-specified ID of the sent message. Only present for newsletter messages.",
                    "type": "string",
                    "example": "wamid.HBgNNTU5OTgxNzY5NTM2FQIAERgSMzNFNzE4QzY5QzE5MjE2RTdB"
                },
                "success": {
                    "description": "Legacy fields for backward compatibility",
                    "type": "boolean",
                    "example": true
                },
                "timestamp": {
                    "description": "The message timestamp returned by the server (Unix timestamp)",
                    "type": "integer",
                    "example": 1640995200
                }
            }
        },
        "types.SendStickerRequest": {
            "type": "object",
            "required": [
                "phone",
                "sticker"
            ],
            "properties": {
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "sticker": {
                    "type": "string",
                    "example": "data:image/webp;base64,UklGRv4..."
                }
            }
        },
        "types.SendTextRequest": {
            "type": "object",
            "required": [
                "body",
                "phone"
            ],
            "properties": {
                "body": {
                    "type": "string",
                    "example": "Hello, World!"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "types.SendVideoRequest": {
            "type": "object",
            "required": [
                "phone",
                "video"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "Video caption"
                },
                "contextInfo": {
                    "$ref": "#/definitions/types.ContextInfo"
                },
                "id": {
                    "type": "string",
                    "example": "custom-message-id"
                },
                "phone": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "video": {
                    "type": "string",
                    "example": "data:video/mp4;base64,AAAAIGZ0eXA..."
                }
            }
        },
        "utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "details": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "utils.SuccessResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "zpmeow WhatsApp API",
	Description:      "A WhatsApp API server built with Go, inspired by wuzapi",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
