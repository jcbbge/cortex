# Next Steps

**Updated:** 2026-02-12
**Branch:** main
**Context:** Initial Cortex setup after meta-analysis of session-mining/ending-session

---

## Current State

Cortex initialized as proper playground for agent orchestration research.

**Completed:**
- Created repository structure
- Moved session mining meta-analysis to `research/001_*`
- Identified core problem: skills are documentation, need execution primitives
- Recognized "hidden architecture" trying to emerge

**Key Finding:** 46 skills exist, only 2-4 used. Skills can't chain, pass state, or auto-trigger. They're suggestions, not enforcement.

---

## Next Steps

### 1. Spawn "Agent Composition Primitives" Exploration

**New session with context:**

```
Read /Users/jcbbge/cortex/research/001_session_mining_meta_analysis.md

Context: Meta-analysis revealed:
- Skills are documentation, not enforcement (wrong abstraction)
- Missing primitives: state passing, auto-chaining, execution model
- 46 skills but only 2-4 used (attention bankruptcy)
- Hidden architecture trying to emerge (skills want to compose)

Task: Design agent composition primitives

Goals:
1. State passing protocol (skills → input/output)
2. Auto-chaining via conditions (workflows)
3. Execution model (not just markdown)
4. Migration strategy (46 skills → proper mechanisms)

Questions to answer:
- What's the execution model for composable agents?
- How does state pass between skill invocations?
- Can we "de-skillify" by promoting to hooks/rules?
- What's the spec for "executable skills"?
- How to migrate incrementally (no big rewrite)?

Output:
- Prototype: executable skill spec (in specs/)
- Design: state passing protocol
- Strategy: de-skillification migration plan
- Examples: 3-5 conversions (skill → hook/rule/workflow)

Constraints:
- Must work with OpenCode/Claude architecture
- Incremental adoption (backwards compatible)
- Reduce cognitive load (simplify, don't add)
```

### 2. De-skillification Audit (After Composition Design)

Run on `/Users/jcbbge/Documents/metaprompts/_skills/`:
- Evaluate all 46 skills against AGENT_UPGRADE_ANALYSIS_CHECKLIST.md
- Identify: which should be hooks? rules? tools? workflows?
- Prioritize: most impactful conversions first
- Document in `migrations/audit_results.md`

### 3. Implement Session-Mining v2 + Ending-Session v5

Apply findings from meta-analysis:
- session-mining: Add full taxonomy evaluation (8 mechanism types)
- ending-session: Add auto-mining trigger logic
- Test integrated workflow

---

## Context

### Key Files

- `research/001_session_mining_meta_analysis.md` - Full findings from meta-analysis
- `/Users/jcbbge/Documents/metaprompts/AGENT_UPGRADE_ANALYSIS_CHECKLIST.md` - 8 mechanism taxonomy
- `/Users/jcbbge/Documents/metaprompts/_skills/session-mining/SKILL.md` - Current (incomplete)
- `/Users/jcbbge/Documents/metaprompts/_skills/ending-session/SKILL.md` - Current (isolated)

### Decisions Made

- **Cortex is the playground** for agent orchestration R&D
- **Skills are wrong abstraction** - promote to hooks/rules where possible
- **State passing is critical** - can't compose without it
- **Incremental migration** - no big rewrite, adapt existing

### Key Insights

1. **Skills = documentation pretending to be enforcement**
2. **12.5% taxonomy coverage** in session-mining (only synthesizes skills)
3. **Hidden architecture emerging** - execution model trying to form
4. **46 skills is attention bankruptcy** - need stronger mechanisms
5. **Dogfooding the analysis** - using broken tools to design fixed tools

---

## Blockers/Questions

**None.** Path is clear:

1. Design composition primitives (new session)
2. Audit existing skills (migration plan)
3. Implement upgrades (v2/v5)
4. Test integrated workflow

---

## Meta-Note

This is the handoff from "meta-analysis session" to "composition primitives exploration."

**We're practicing what we're designing:** using session-mining + ending-session to create artifacts that spawn the next session.

The tools examining themselves. 🔄
