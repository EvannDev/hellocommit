# HelloCommit

A platform to discover and contribute to open source projects by finding "good first issues" from your GitHub starred repositories.

## Architecture

```
hellocommit/
├── backend-go/      # Go API with Fiber v3
├── frontend/        # Next.js 16 with shadcn/ui
├── docker-compose.yml
└── AGENTS.md
```

## Tech Stack

| Layer | Tech |
|-------|------|
| Frontend | Next.js 16, React 19, shadcn/ui, Tailwind v4 |
| Auth | NextAuth.js (GitHub OAuth) |
| Backend | Go 1.22+, Fiber v3 |
| Database | SQLite (prototype) |

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- GitHub OAuth App (create at https://github.com/settings/applications/new)

### Setup

1. **Clone and configure environment**

```bash
cp backend-go/.env.example backend-go/.env
cp frontend/.env.example frontend/.env.local
```

2. **Fill in environment variables**

In `backend-go/.env`:
- `GITHUB_TOKEN`: Your GitHub personal access token (for API rate limits)

In `frontend/.env.local`:
- `AUTH_SECRET`: Generate with `openssl rand -base64 32`
- `AUTH_GITHUB_ID`: Your GitHub OAuth App Client ID
- `AUTH_GITHUB_SECRET`: Your GitHub OAuth App Client Secret

3. **Run with Docker Compose**

```bash
docker-compose up --build
```

Or run individually:

**Backend:**
```bash
cd backend-go
go mod tidy
go run ./cmd/api
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

4. **Access the app**

- Frontend: http://localhost:3000
- API: http://localhost:8080/api/health

## GitHub OAuth App Setup

1. Go to https://github.com/settings/applications/new
2. Fill in:
   - Application name: HelloCommit
   - Homepage URL: http://localhost:3000
   - Authorization callback URL: http://localhost:3000/api/auth/callback/github
3. Copy Client ID and Client Secret to `.env.local`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/users` | Create/update user |
| GET | `/api/users/:id` | Get user by ID |
| GET | `/api/users/:id/starred` | Get starred repos |
| POST | `/api/users/:id/sync` | Sync starred repos |
| GET | `/api/repos/:owner/:name/issues` | Get issues for a repo |
| GET | `/api/issues/good-first` | Get good first issues |
| POST | `/api/sync/all/:userId` | Sync all data |

## Development

```bash
# Backend
cd backend-go
go build -o main ./cmd/api

# Frontend
cd frontend
npm run lint
```
