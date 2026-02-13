# Cortex

**Meta-cognitive architecture for AI-assisted development.**

## Vision

Agent orchestration isn't about skills—it's about **execution primitives** that compose reliably.

This project explores:
- **State passing** between agent invocations
- **Automatic chaining** via hooks and workflows
- **De-skillification**: promoting patterns from documentation to enforcement
- **Execution models** for composable agent mechanisms
- **Meta-cognitive tooling** that improves itself

## The Problem

- 46 skills, only 2-4 used (attention bankruptcy)
- Skills are "pretty please try this" (suggestions, not enforcement)
- No state passing (agents can't compose)
- No auto-chaining (manual orchestration)
- Wrong abstraction (skills should be hooks/rules/tools)

## The Insight

**Skills are documentation pretending to be enforcement.**

What we actually need:
- **Rules** → constraints (MUST)
- **Hooks** → automation (WHEN X, DO Y)
- **Tools** → capabilities (CAN)
- **Workflows** → composition (THEN)
- **State** → context preservation (PASS)

Skills should be the exception, not the default.

## Project Structure

```
cortex/
├── research/           # Session analyses, findings
├── specs/              # Execution model designs
├── prototypes/         # Experimental implementations
├── migrations/         # Skill → mechanism conversions
└── meta/               # Tools for improving tools
```

## Current Focus

**Phase 1:** Agent composition primitives
- Design state passing protocol
- Prototype executable skill spec
- Create migration path from current skills
- Build meta-cognitive feedback loop

## Meta-Observation

This project uses the broken tooling to design the fixed tooling.

**Dogfooding the dogfooding analysis.**

---

*Started: 2026-02-12*
*Context: Meta-analysis of session-mining revealed systematic gaps*
