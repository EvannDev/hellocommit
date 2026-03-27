# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HelloCommit is a full-stack web app that helps users discover "good first issues" from their GitHub starred repositories. It uses a Next.js 16 frontend and a Go/Fiber backend with SQLite.

## Commands

### From repo root (preferred)

```bash
make dev          # Run frontend + backend in parallel
make build        # Build both services
make lint         # ESLint + go vet
make type-check   # TypeScript type check
make test         # Go tests
make docker-up    # docker compose up --build
```

### Frontend only (`ui/`)

```bash
pnpm dev          # Dev server at http://localhost:3000
pnpm build        # Production build
pnpm lint         # Run ESLint — run before committing
pnpm type-check   # tsc --noEmit
```

### Backend only (repo root)

```bash
go run ./cmd/api  # Start API server at http://localhost:8080
go mod tidy       # Sync dependencies
```

### Docker (full stack)

```bash
docker compose up --build  # Build and run both services
```

## Architecture

**Frontend** (Next.js 16, React 19, TypeScript, Tailwind CSS v4, NextAuth.js v5) — lives in `ui/`:
- App Router with Server Components by default (`"use client"` only when needed)
- `ui/auth.ts` — GitHub OAuth via NextAuth
- `ui/lib/api.ts` — typed API client for the Go backend
- `ui/types/index.ts` — shared TypeScript interfaces
- `ui/app/(dashboard)/` — protected routes (dashboard, repos, settings)
- UI components follow shadcn/ui conventions using CVA for variants

**Backend** (Go 1.25, Fiber v3, SQLite via `modernc.org/sqlite`) — lives at repo root:
- `cmd/api/main.go` — entry point, routing
- `internal/` — layered: `handlers` → `services` → `repositories` → `database`
- `pkg/github/` — GitHub API client wrapper using `google/go-github`
- SQLite DB with `users`, `repos`, `issues` tables

**Data flow:** User authenticates via GitHub OAuth → frontend sends access token to backend → backend syncs starred repos and issues from GitHub API → data stored in SQLite → frontend displays filtered good first issues.

## Frontend Code Conventions

**This is Next.js 16** — APIs and conventions may differ from training data. Read `node_modules/next/dist/docs/` before writing Next.js-specific code.

- Use Server Components by default; add `"use client"` only for hooks, event handlers, or browser APIs
- Props in server components: `Readonly<{...}>`
- Tailwind CSS v4: use `@import "tailwindcss"` syntax; CSS variables in `@theme inline {}`; oklch color format
- No CSS modules — Tailwind only
- Imports: external packages → internal (`@/`) → relative → types (no blank lines within groups)
- Use `type` for aliases, `interface` for object shapes; avoid `any`
- Component files: kebab-case (`user-profile.tsx`); CVA for variants; `data-slot` attribute on root elements

## Environment Variables

Frontend (`.env.local`):
```
AUTH_SECRET=
AUTH_GITHUB_ID=
AUTH_GITHUB_SECRET=
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```
