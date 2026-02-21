---
name: "soccer-cli OpenClaw Skill"
description: "A skill for OpenClaw that provides soccer data via the API-Football client. It wraps the existing `soccer-cli` binary to expose commands that OpenClaw can invoke programmatically."
---

# Overview
This skill enables OpenClaw to fetch soccer scores, game details, and player statistics using the existing `soccer-cli` implementation.

# Installation
1. Ensure the `soccer-cli` binary is built and available in the system PATH (or specify the absolute path).
2. Place this skill directory in OpenClaw's `skills` folder, e.g. `~/.openclaw/skills/openclaw_soccer_cli/`.

# Usage
OpenClaw can invoke the skill by executing the binary with the appropriate sub‑command. Example commands:

```bash
# Get latest scores for a team
soccer-cli scores "Manchester United"

# Get detailed events for a fixture
soccer-cli game 123456

# Get squad and player stats for a fixture
soccer-cli squad 123456
```

The skill returns plain‑text output which OpenClaw can capture and forward to the user.

# Configuration
Create a configuration file at `~/.config/soccer-cli/config.yaml` with your API‑Football key:

```yaml
apikey: YOUR_API_KEY_HERE
```

# Integration Details
- The skill does **not** require additional code changes; it reuses the existing Go client and CLI.
- If you need to customize timeouts or logging, refer to the implementation plan in `implementation_plan.md`.
- For testing, you can run `go test ./...` to ensure the client works before deploying the skill.

# Advanced
If you wish to call the client programmatically (e.g., from a script), you can import the Go package `pkg/api` and use the `Client` methods directly. This is optional for more complex OpenClaw workflows.

# License
This skill is MIT‑licensed, matching the main project.
