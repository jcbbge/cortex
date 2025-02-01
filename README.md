# Cortex

A terminal-native AI assistant designed to unify development workflows through natural language interaction, context-aware memory, and direct command execution.

## Status: Early Development 🚧
This is my first Go tool, currently in active development. Core features are being implemented.

## Overview
Cortex solves context switching between development tools by:
- Providing a **unified terminal interface** for natural language and commands
- Maintaining **contextual memory** of decisions, code, and conversations
- Acting as an **AI pair programmer** that understands project history

## Key Features (In Development)
- 🖥️ **Modern TUI**: Fluid, dynamic interface with ambient intelligence
- 🧠 **Associative Memory**: Vector-based search with PostgreSQL/pgvector
- ⚡ **LLM Orchestration**: Model-agnostic integration
- 🔒 **Safety First**: Command allow/blocklist system
- 📊 **Token Awareness**: Real-time usage tracking

## Prerequisites
- Go 1.21+
- Docker & Docker Compose
- OpenAI/OpenRouter API access

## Development Setup

### 1. Clone & Install Tools
# Clone repository
git clone https://github.com/jcbbge/cortex.git
cd cortex

# Install development tools
make install-tools
make deps

### 2. Database Setup
# Start PostgreSQL with pgvector
docker compose up -d

# Environment variables (create .env file)
DB_PASSWORD=cortexai

The database will be available at:
- Host: localhost
- Port: 5433
- User: cortexai
- Database: cortexai
- Password: from DB_PASSWORD or default 'cortexai'

### 3. Run Development Server
# Run with hot reload
make dev

## Project Structure
.
├── cmd/cortex/     # Main application entry
├── internal/       # Private application code
├── pkg/           # Public libraries
├── ai/            # Project documentation
└── scripts/       # Utility scripts

## Docker Configuration

### Database Container
The project uses PostgreSQL 15 with pgvector for similarity search:

### Database Optimizations
PostgreSQL is configured with optimized settings for vector operations:
- HNSW indexing enabled
- Shared buffers: 1GB
- Effective cache size: 3GB
- Maintenance work mem: 1GB
- Vector-specific settings:
  - `pgvector.hnsw_ef_search = 100`
  - `pgvector.hnsw_m = 16`

## Contributing
This is a personal learning project as I explore Go development. Feel free to open issues for bugs or feature suggestions.

## License
MIT
