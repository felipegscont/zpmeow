# zpmeow - WhatsApp API

Uma API WhatsApp construída em Go, inspirada no wuzapi, com funcionalidades avançadas e documentação completa.

## 🚀 Funcionalidades

- ✅ **Gerenciamento de Sessões** - Criar, conectar e gerenciar múltiplas sessões WhatsApp
- ✅ **Reconexão Automática** - Sessões reconectam automaticamente após restart (como wuzapi)
- ✅ **Envio de Mensagens** - Texto, imagem, áudio, vídeo, documentos, stickers, etc.
- ✅ **Gerenciamento de Grupos** - Criar, gerenciar e interagir com grupos
- ✅ **Chat Management** - Marcar como lido, reagir, deletar, editar mensagens
- ✅ **Documentação Swagger** - API totalmente documentada
- ✅ **Banco de Dados** - PostgreSQL com interface de gerenciamento
- ✅ **Logs Estruturados** - Sistema de logging avançado

## 🛠️ Tecnologias

- **Go** - Linguagem principal
- **PostgreSQL** - Banco de dados
- **Redis** - Cache (opcional)
- **MinIO** - Armazenamento de arquivos
- **DBGate** - Interface de gerenciamento do banco
- **Docker** - Containerização
- **Swagger** - Documentação da API

## 📦 Instalação

### Usando Docker (Recomendado)

```bash
# Clonar o repositório
git clone <repository-url>
cd zpmeow

# Iniciar todos os serviços
docker compose up -d

# Criar o banco de dados
make db-create

# Iniciar a aplicação
make run
```

### Desenvolvimento Local

```bash
# Instalar dependências
make deps

# Iniciar apenas o banco (PostgreSQL)
docker compose up -d postgres

# Criar banco e iniciar app
make dev-local
```

## 🌐 Acesso aos Serviços

- **API Principal:** http://localhost:8080
- **Documentação Swagger:** http://localhost:8080/swagger/index.html
- **DBGate (Gerenciamento DB):** http://localhost:3000
- **MinIO Console:** http://localhost:9001

## 📚 Documentação

- [Documentação da API](http://localhost:8080/swagger/index.html) - Swagger UI completo
- [DBGate - Gerenciamento do Banco](docs/DBGATE.md) - Como usar o DBGate
- [Arquitetura](ARQUITETURA.md) - Documentação da arquitetura

## 🔧 Comandos Úteis

```bash
# Gerar documentação Swagger
make swagger

# Gerenciar banco de dados
make db-create    # Criar banco
make db-drop      # Deletar banco
make db-reset     # Resetar banco
make db-test      # Testar conexão

# Docker
make up           # Iniciar serviços
make down         # Parar serviços
```

## 📊 Gerenciamento do Banco

O zpmeow inclui o **DBGate** para facilitar o gerenciamento do PostgreSQL:

```bash
# Iniciar DBGate
docker compose up -d dbgate

# Acessar: http://localhost:3000
```

Veja mais detalhes em [docs/DBGATE.md](docs/DBGATE.md).

## 🔄 Reconexão Automática

Inspirado no wuzapi, o zpmeow reconecta automaticamente sessões que estavam conectadas:

- Sessões com status `connected` são reconectadas no startup
- Verifica credenciais antes de tentar reconectar
- Logs detalhados do processo de reconexão

## 📝 Exemplo de Uso

```bash
# 1. Criar uma sessão
curl -X POST http://localhost:8080/sessions/create \
  -H "Content-Type: application/json" \
  -d '{"name": "Minha Sessão"}'

# 2. Conectar a sessão
curl -X POST http://localhost:8080/sessions/{session-id}/connect

# 3. Obter QR Code
curl http://localhost:8080/sessions/{session-id}/qr

# 4. Enviar mensagem
curl -X POST http://localhost:8080/session/{session-id}/send/text \
  -H "Content-Type: application/json" \
  -d '{"to": "5511999999999@s.whatsapp.net", "text": "Olá!"}'
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.
