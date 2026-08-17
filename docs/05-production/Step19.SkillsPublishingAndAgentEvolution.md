# Step 19 — Skills, Publishing, and Agent Evolution

**Knowledge depth: 8/10**

Read [14 — Skills Runtime](../14-skills-runtime.md), [15 — Core Skills System](../15-core-skills-system.md), and [16 — Skill Publishing System](../16-skill-publishing.md) in that order. They explain runtime dependencies, bundled skills, access rules, and the publishing lifecycle. Read [21 — Agent Evolution and Skill Management](../21-agent-evolution-and-skill-management.md) before allowing the system to suggest or apply changes to itself.

## From local instructions to managed capabilities

Early in the course, a skill is simply a selected instruction package. In a platform, skills need metadata, versions, access policies, dependencies, and a predictable runtime environment.

```text
skill source → validate/package → publish version → grant access → select for a run
```

This lifecycle turns informal prompt files into something an operator can understand and govern.

## Publishing is not execution

Publishing a skill should not immediately run its scripts or grant it to every agent. Treat installation, dependency setup, authorization, and activation as separate concepts. This makes the system easier to reason about and lets the UI present a safe management flow.

## Evolution needs evidence

GoClaw's evolution design starts with metrics, produces suggestions, and places guardrails around any application or rollback. The teaching point is simple: an agent should not silently rewrite its own behavior because one response looked poor. First make runs observable; then make changes reviewable.

