# AGENTS.md

## Project Overview

Litho API is a backend administration API service providing foundational modules: users, organizations, roles, permissions, and routes.

- **Language**: Go
- **HTTP Framework**: [Hertz](https://github.com/cloudwego/hertz)
- **Database**: PostgreSQL
- **Cache**: Redis
- **License**: BSD-3-Clause

## Repository Layout

- `main.go` — application entry point
- `bootstrap/` — application startup and initialization
- `api/` — HTTP handlers
- `model/` — data models
- `common/` — shared utilities
- `config/` — configuration
- `builder/` — build-related helpers
- `Dockerfile` — container image build
- `.cnb.yml` — CNB pipeline configuration

## CNB Platform

This project is hosted on [CNB](https://cnb.cool) (repo: `kainonly/litho-api`). CI/CD is driven by `.cnb.yml`:

- **Branch push**: syncs code to the GitHub mirror repository
- **Tag push**: syncs the tag to GitHub, then builds and pushes a Docker image to the CNB registry

CNB-related skills (`cnb-*`) are available globally for interacting with issues, pull requests, pipelines, and the package registry via the `cnb` CLI.

## Conventions

- Keep commit messages and PR titles in semantic commit format; PR titles must not contain parentheses or branch names
- Reference related issues in PR descriptions with `Ref: #<ISSUE_ID>`
- Do not merge or close PRs as an AI; leave that to humans
- After creating or pushing a PR, stop immediately — do not poll CI or review status
