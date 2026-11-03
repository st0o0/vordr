## ADDED Requirements

### Requirement: Community files
The repository SHALL include CODE_OF_CONDUCT.md (Contributor Covenant v2.1), CONTRIBUTING.md (dev setup, PR rules), SECURITY.md (private reporting via GitHub advisories), and NOTICE (copyright + dependency summary).

#### Scenario: Community files present
- **WHEN** the repository is cloned
- **THEN** CODE_OF_CONDUCT.md, CONTRIBUTING.md, SECURITY.md, and NOTICE SHALL exist at the repository root

### Requirement: License
The repository SHALL be licensed under MIT with a LICENSE file at the repository root.

#### Scenario: License file
- **WHEN** inspecting the repository root
- **THEN** a LICENSE file with MIT license text SHALL exist

### Requirement: GitHub Issue templates
The repository SHALL include issue templates for bug reports and feature requests, with blank issues disabled and security issues routed to GitHub advisories.

#### Scenario: Bug report template
- **WHEN** a user creates a new issue
- **THEN** they SHALL see a structured bug report form with fields for what happened, reproduction steps, logs, image tag, and host info

#### Scenario: Feature request template
- **WHEN** a user creates a new issue
- **THEN** they SHALL see a structured feature request form with fields for problem, proposal, and alternatives

### Requirement: AI configuration files
The repository SHALL include CLAUDE.md with project context, architecture overview, build commands, and conventions. It SHALL also include AGENTS.md with code change instructions and OpenSpec usage guidance.

#### Scenario: CLAUDE.md content
- **WHEN** an AI agent reads CLAUDE.md
- **THEN** it SHALL understand the project purpose, how to build/test/lint, and what conventions to follow

#### Scenario: AGENTS.md content
- **WHEN** an AI agent starts a code change
- **THEN** AGENTS.md SHALL instruct it to read CLAUDE.md first, run build+test+lint after changes, and use OpenSpec for planning

### Requirement: Git and Docker ignore files
The repository SHALL include .gitignore and .dockerignore files matching the standard template used across all Go repos.

#### Scenario: Binary not tracked
- **WHEN** running `go build -o vordr ./cmd/vordr`
- **THEN** the `vordr` binary SHALL be excluded by .gitignore

#### Scenario: Docker context minimal
- **WHEN** running `docker build`
- **THEN** .git, .github, tests, docs, *.md (except LICENSE), and tooling configs SHALL be excluded from the build context
