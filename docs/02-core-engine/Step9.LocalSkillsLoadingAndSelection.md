# Step 9 — Local Skills Loading and Selection

**Knowledge depth: 6/10**

Load and select small local SKILL.md instruction packages safely.

## Step outcome

Enabled skills are validated, ordered deterministically, and injected without granting tool permissions.

## Task 1 — Create local skills

### Theory

- **Architectural role:** The loader converts a filesystem package into immutable `SkillContent`. A skill changes model guidance, while the tool registry and policy still control permission to act.
- **Why:** Validating names, paths, and size at the boundary prevents traversal and uncontrolled context growth. Frontmatter supports discovery without loading the full instructions.
- **GoClaw reference:** Inspect the parser and loader in [`goclaw/internal/skills/loader.go`](../../goclaw/internal/skills/loader.go), metadata search in [`goclaw/internal/skills/search.go`](../../goclaw/internal/skills/search.go), path guards in [`goclaw/internal/skills/guard.go`](../../goclaw/internal/skills/guard.go), and a real package such as [`goclaw/skills/workspace-organizing/SKILL.md`](../../goclaw/skills/workspace-organizing/SKILL.md).

A skill is a package of task instructions, usually a `SKILL.md` plus optional scripts or references. It changes model guidance; it does not bypass the tool registry.

Skill publishing, versioned managed storage, and runtime dependency installation are not implemented in this course.

### Goal

Learn to package instructions and metadata as selectable modules instead of placing every instruction in the system prompt.

### Guide to implement

Create two small skills:

```text
skills/
├── writing/
│   └── SKILL.md
└── workspace/
    └── SKILL.md
```

Use simple YAML frontmatter:

```markdown
---
name: writing
description: Improve and structure user-facing prose.
---

# Instructions

Use direct language. Preserve facts supplied by the user.
```

Implement a loader that validates the directory name, reads only `SKILL.md`, limits file size, parses frontmatter, and returns immutable `SkillContent` values.

## Task 2 — Select and inject skills

### Theory

- **Architectural role:** Agent settings define the allow-list, the selector resolves slugs, the loader reads and validates files, and the prompt builder only formats selected content. Each layer has one reason to change.
- **Why:** Loading every skill wastes tokens and can create conflicting instructions. Stable sorting gives the same prompt for the same input, which helps caching, snapshot tests, and debugging.
- **GoClaw reference:** Review visibility and grant filtering in [`goclaw/internal/skills/visibility.go`](../../goclaw/internal/skills/visibility.go), BM25 selection in [`goclaw/internal/skills/search.go`](../../goclaw/internal/skills/search.go), and how skill filters enter the prompt in [`goclaw/internal/agent/loop_history.go`](../../goclaw/internal/agent/loop_history.go).

Loading every skill wastes context and can create conflicting instructions. Start with explicit selection; semantic auto-selection can be added later.

### Goal

Add only relevant, enabled instructions to the current run context.

### Guide to implement

Read `enabled_skills` from the selected agent's settings. For each enabled slug:

1. Resolve it under the configured `skills/` root.
2. Reject traversal and symlink escapes.
3. Load the validated `SKILL.md`.
4. Sort selected skills by slug for deterministic prompts.
5. Add them under a clearly marked `Selected skills` section.

Missing or invalid skills should produce a clear configuration error. They should not disappear silently.
