# Skill Audit: 46 Skills → Classification

**Date:** 2026-02-12
**Goal:** Classify each skill as KEEP/PROMOTE/MERGE/DELETE
**Action:** Reduce cognitive load, strengthen enforcement

---

## Classification Criteria

- ✓ **KEEP** - Domain knowledge I can't infer (SolidJS, security patterns, etc.)
- → **PROMOTE** - Should be hook/rule/tool (automated, not invoked)
- ⚠️ **MERGE** - Overlaps with another skill
- ❌ **DELETE** - Never used, could infer, or redundant

---

## Skills Audit

### ✓ KEEP (Domain Knowledge - Actually Useful)

1. **building-with-solidjs** - Framework-specific patterns I wouldn't know
2. **building-with-solidstart** - SolidStart meta-framework specifics
3. **backend-first-security** - Opinionated security architecture (Supabase RLS)
4. **debugging-with-logs** - Wide events pattern (specific, non-obvious)
5. **designing-apis** - REST API best practices (useful checklist)
6. **self-hosting-vps-ubuntu** - Infrastructure setup (domain knowledge)
7. **security-checklist** - VPS hardening steps (concrete checklist)

**Total: 7 skills**

---

### → PROMOTE (Should Be Hook/Rule/Tool)

8. **timestamp-protocol** - Should be RULE in CLAUDE.md: "Begin each message with timestamp"
9. **fresh-eyes-review** - Could be HOOK: PostToolUse after Write → run review
10. **refactoring-framework** - Could be HOOK: PreToolUse before Edit → remind of dimensions
11. **ending-session** - Could be HOOK: Stop event → check if NEXT_STEPS.md needed
12. **starting-session** - Could be HOOK: Init event → read NEXT_STEPS.md if exists

**Total: 5 skills → promote to automation**

---

### ⚠️ MERGE (Overlapping Functionality)

13. **refactor** + **refactoring-framework** - Same purpose, merge into one
14. **session-mining** + **metaprompt-process** - Both extract patterns, consolidate
15. **evaluating-product-ideas** + **evaluating-business-strategy** - Similar evaluation, merge
16. **prd** + **project-documentation** - Both documentation generation, consolidate
17. **meta-agent-template** + **cognitive-approaches** - Both meta-cognitive frameworks, merge

**Total: 10 skills → 5 merged skills**

---

### ❌ DELETE (Never Used / Could Infer)

18. **anima** - Persistent memory (not actually used in practice)
19. **chat-directives** - Behavior modifiers (too meta, not actionable)
20. **collaborating-partner-mode** - Default mode (shouldn't need skill for this)
21. **creative-systems-thinking** - Ted Nelson thinking (interesting but unused)
22. **housekeeping** - Self-reflection for metaprompts (one-time use)
23. **idea-wizard** - We just used it, but I could do diverge/converge without it
24. **knowledge-graph-system** - PKM architecture (aspirational, not implemented)
25. **metaprompts** - Access library (redundant, just read files)
26. **metasighting-perspectives** - Metacognitive engine (too abstract)
27. **non-technical-async** - Explain to non-tech (could infer)
28. **opencode-ecosystem-analysis** - Multi-repo analysis (specific one-time use)
29. **openrouter** - Display models (utility, not skill)
30. **robot-mode-maker** - CLI design (could infer)
31. **step-workflow** - Phase-gated workflow (generic, could infer)
32. **big-brained-optimizer** - Deep optimization (could infer with constraints)
33. **reframing-problems** - Challenge problem formulation (could infer)
34. **repo-deep-dive-analysis** - Ingest repos (could do without skill)
35. **upgrade-analysis** - Evaluate against checklist (meta, one-time)
36. **workflow-runner** - Execute chains (infrastructure, not skill)

**Total: 19 skills to delete**

---

### 🤔 UNDECIDED (Need Your Input)

37. **metaprompt-process** - Process inbox captures (do you use this?)
38. **rams** - Accessibility review (do you use this?)

**Total: 2 skills - need user decision**

---

## Summary

| Category | Count | Action |
|----------|-------|--------|
| KEEP | 7 | No change - these are useful |
| PROMOTE | 5 | Convert to hooks/rules |
| MERGE | 10→5 | Consolidate overlapping |
| DELETE | 19 | Remove unused |
| UNDECIDED | 2 | User decides |

**Result: 46 skills → ~17 skills (63% reduction)**

---

## Immediate Actions

1. **Delete the 19 unused skills**
2. **Promote 5 to hooks/rules** (automate them)
3. **Merge 5 pairs** into consolidated versions
4. **Keep 7 domain knowledge** skills

**Net effect: Cognitive load drops dramatically. Only keep what's genuinely useful.**

