# Cortex: AI-Powered Development Assistant

## Overview
Cortex is a terminal-native AI assistant designed to unify fragmented development workflows by integrating natural language interaction, context-aware memory, and direct command execution into a single interface.

### Problem Statement
Developers juggle multiple tools (IDEs, terminals, AI chatbots), leading to context switching and productivity loss. Cortex solves this by:
- Providing a **unified terminal interface** for natural language and commands.
- Maintaining a **contextual memory** of decisions, code, and conversations.
- Acting as an **AI pair programmer** that understands project history.

### Key Features
- 🖥️ **Modern TUI**: Polished terminal interface with syntax highlighting and animations.
- 🧠 **Associative Memory**: Recall past decisions using vector search and metadata.
- ⚡ **LLM Orchestration**: Model-agnostic integration with OpenAI/OpenRouter/etc.
- 🔒 **Safety First**: Block dangerous commands unless explicitly allowed.
- 📊 **Token Awareness**: Real-time tracking of LLM usage and rate limits.

### Non-Features
- ❌ No terminal emulation (uses native shell environment).
- ❌ No file versioning (relies on Git/IDE).
- ❌ No sandboxed execution (runs in user’s environment).

---

# Cortex Product Requirements Document (PRD)

## Objectives
1. Reduce context switching between development tools.
2. Provide AI-assisted decision-making with project-specific memory.
3. Maintain developer workflow integrity with minimal friction.

## Functional Requirements
| Component         | Requirements                                                                 |
|-------------------|-----------------------------------------------------------------------------|
| **TUI**           | Dynamic, fluid interface with ambient intelligence and spatial memory.       |
|                   | Organic command suggestions and context-aware animations.                    |
|                   | Subtle status indicators for system metrics (tokens/memory).                 |
|                   | Natural conversation flow with spatial navigation.                           |
| **Memory Layer**  | Hybrid search (vector + tags), auto-archiving of old memories.              |
| **LLM Integration**| Model-agnostic client with response streaming and token tracking.          |
| **Safety**        | User-configurable allow/blocklist with confirmation for destructive commands. |
| **Config**        | YAML/ENV configuration for API keys, themes, and default models.           |

## Non-Functional Requirements
| Aspect            | Requirements                                                                 |
|-------------------|-----------------------------------------------------------------------------|
| Performance       | <16ms animation frame time; <200ms latency for memory recalls.              |
| Aesthetics        | Fluid animations; organic transitions; ambient intelligence.                 |
| Security          | API keys encrypted in transit; no storage of sensitive command outputs.    |
| Compatibility     | macOS-first; PostgreSQL 15+ with pgvector.                                 |
| Observability     | Structured JSON logging with Zerolog; Prometheus metrics for DB performance.|

## Key Metrics
- ⏱️ **Query Latency**: 95th percentile <2s
- 🧠 **Memory Recall Accuracy**: >80% relevance in hybrid search
- 🔄 **LLM Usage**: <10% rate limit breaches per user/month
- 🎨 **Animation Performance**: 60fps smooth transitions
- 🌊 **Fluidity Score**: <2 frame drops per minute

## Dependencies
- PostgreSQL 15+ with pgvector
- OpenAI/OpenRouter API access
- Bubble Tea (TUI framework)
