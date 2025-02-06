# Cortex: AI-Powered Development Assistant

## Overview
Cortex is a terminal-native AI assistant designed to unify fragmented development workflows by integrating natural language interaction, context-aware memory, and direct command execution into a single interface.

### Development Philosophy
- **User Experience First**: Every feature begins with how it feels to use, not how it works internally.
- **Immediate Impact**: Users should be delighted from their first interaction with the welcome screen.
- **Progressive Enhancement**: Start minimal, layer in complexity only where it enhances the experience.
- **Rapid Iteration**: Development focuses on quick feedback loops and continuous refinement.

### Problem Statement
Developers juggle multiple tools (IDEs, terminals, AI chatbots), leading to context switching and productivity loss. Cortex solves this by:
- Providing a **unified terminal interface** for natural language and commands.
- Maintaining a **contextual memory** of decisions, code, and conversations.
- Acting as an **AI pair programmer** that understands project history.

### Key Features
- 🖥️ **Modern TUI**: Polished terminal interface with fluid animations and organic interactions.
- 🧠 **Associative Memory**: Recall past decisions using vector search and metadata.
- ⚡ **LLM Orchestration**: Model-agnostic integration with OpenAI/OpenRouter/etc.
- 🔒 **Safety First**: Block dangerous commands unless explicitly allowed.
- 📊 **Token Awareness**: Real-time tracking of LLM usage and limits.

### Design Principles
- **Minimal Yet Powerful**: Each feature serves a clear purpose with no unnecessary complexity.
- **Visual Hierarchy**: Important information stands out through subtle visual cues.
- **Responsive Flow**: Interface adapts smoothly to user actions and terminal constraints.
- **Contextual Awareness**: UI elements reflect system state and user intent.
- **Intuitive Interaction**: Commands and patterns follow natural developer workflows.

### Non-Features
- ❌ No terminal emulation (uses native shell environment)
- ❌ No file versioning (relies on Git/IDE)
- ❌ No sandboxed execution (runs in user's environment)

---

# Cortex Product Requirements Document (PRD)

## Objectives
1. Reduce context switching between development tools
2. Provide AI-assisted decision-making with project-specific memory
3. Maintain developer workflow integrity with minimal friction
4. Create a delightful and empowering user experience

## Development Priorities
1. **Core Experience**: Perfect the welcome screen and basic interaction patterns
2. **Essential Flow**: Streamline conversation and command handling
3. **Visual Feedback**: Implement minimal but effective animations and state indicators
4. **Progressive Features**: Layer in advanced capabilities without compromising simplicity

## Functional Requirements
| Component         | Requirements                                                                 |
|-------------------|-----------------------------------------------------------------------------|
| **TUI**           | Dynamic, fluid interface with ambient intelligence and spatial memory        |
|                   | Organic command suggestions and context-aware animations                     |
|                   | Subtle status indicators for system metrics (tokens/memory)                  |
|                   | Natural conversation flow with spatial navigation                            |
| **Memory Layer**  | Hybrid search (vector + tags), auto-archiving of old memories               |
| **LLM Integration**| Model-agnostic client with response streaming and token tracking           |
| **Safety**        | User-configurable allow/blocklist with confirmation for destructive commands |
| **Config**        | YAML/ENV configuration for API keys, themes, and default models             |

## Non-Functional Requirements
| Aspect            | Requirements                                                                 |
|-------------------|-----------------------------------------------------------------------------|
| Performance       | <16ms animation frame time; <200ms latency for memory recalls               |
| Aesthetics        | Fluid animations; organic transitions; ambient intelligence                  |
| Security          | API keys encrypted in transit; no storage of sensitive command outputs       |
| Compatibility     | macOS-first; PostgreSQL 15+ with pgvector                                   |
| Observability     | Structured JSON logging with Zerolog; Prometheus metrics for DB performance |

## Key Metrics
- ⏱️ **Query Latency**: 95th percentile <2s
- 🧠 **Memory Recall Accuracy**: >80% relevance in hybrid search
- 🔄 **LLM Usage**: <10% rate limit breaches per user/month
- 🎨 **Animation Performance**: 60fps smooth transitions
- 🌊 **Fluidity Score**: <2 frame drops per minute

## Implementation Strategy
1. **MVP Focus**: Start with core interaction patterns
2. **Iterative Enhancement**: Add features progressively based on user feedback
3. **Visual Refinement**: Polish animations and transitions continuously
4. **Performance Optimization**: Monitor and optimize as complexity grows

## Dependencies
- PostgreSQL 15+ with pgvector
- OpenAI/OpenRouter API access
- Bubble Tea (TUI framework)