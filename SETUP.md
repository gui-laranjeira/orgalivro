# Setup & Deployment

## Prerequisites

| Tool | Version | Notes |
|---|---|---|
| Go | 1.22+ | `go version` |
| Node.js | 18+ | `node --version` |
| pnpm | 8+ | `npm install -g pnpm` |

No C compiler needed — the SQLite driver is pure Go.

---

## Development

### 1. Backend

```bash
cd backend
go run ./cmd/server
```

The API starts on `http://localhost:8080`. The SQLite database is created automatically at `./orgalivro.db` on first run.

**Hot reload** (optional — requires [air](https://github.com/air-verse/air)):

```bash
go install github.com/air-verse/air@latest
cd backend
air
```

**Environment variables** (all optional):

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `DB_PATH` | `./orgalivro.db` | SQLite database file path |
| `ALLOWED_ORIGINS` | `http://localhost:5173` | CORS allowed origins (comma-separated) |

### 2. Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

The dev server starts on `http://localhost:5173` and proxies `/api` to `http://localhost:8080` automatically.

---

## Production (single-machine / self-hosted)

### Build

```bash
# Build frontend
cd frontend
pnpm install
pnpm build
# Output: frontend/dist/

# Build backend binary
cd backend
go build -o orgalivro ./cmd/server
```

### Run

The simplest deployment: serve the frontend static files from the Go binary itself.

**Option A — Nginx reverse proxy (recommended)**

1. Copy `frontend/dist/` to your web root (e.g. `/var/www/orgalivro`).
2. Configure Nginx:

```nginx
server {
    listen 80;
    server_name your-host;

    root /var/www/orgalivro;
    index index.html;

    # All non-API routes serve the SPA
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Proxy API calls to the Go backend
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

3. Start the backend:

```bash
DB_PATH=/var/lib/orgalivro/orgalivro.db \
ALLOWED_ORIGINS=http://your-host \
PORT=8080 \
./orgalivro
```

**Option B — systemd service**

Create `/etc/systemd/system/orgalivro.service`:

```ini
[Unit]
Description=Orgalivro backend
After=network.target

[Service]
Type=simple
User=orgalivro
WorkingDirectory=/opt/orgalivro
ExecStart=/opt/orgalivro/orgalivro
Environment=PORT=8080
Environment=DB_PATH=/var/lib/orgalivro/orgalivro.db
Environment=ALLOWED_ORIGINS=http://your-host
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now orgalivro
```

---

## Docker (optional)

Create `Dockerfile` in the repo root:

```dockerfile
# ---- Build frontend ----
FROM node:18-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm build

# ---- Build backend ----
FROM golang:1.22-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN go build -o /orgalivro ./cmd/server

# ---- Final image ----
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /orgalivro ./orgalivro
COPY --from=frontend /app/frontend/dist ./static

EXPOSE 8080
ENV DB_PATH=/data/orgalivro.db

CMD ["./orgalivro"]
```

> **Note:** The Go binary does not currently serve the frontend static files directly. You would need to embed `frontend/dist/` using Go's `embed` package or add Nginx as a sidecar to serve static files and proxy `/api`. See Option A above.

Build and run:

```bash
docker build -t orgalivro .
docker run -p 8080:8080 -v orgalivro-data:/data orgalivro
```

---

## Future: Tauri Desktop App

The project is architected for eventual Tauri packaging:

- Hash-based routing (`#/`) works in WebView without a server
- CGO-free SQLite driver allows the Go binary to run as a Tauri [sidecar](https://tauri.app/develop/sidecar/)
- The frontend build output is identical for both web and desktop targets

When ready:

1. Add Tauri to the frontend: `pnpm add -D @tauri-apps/cli`
2. Configure the Go binary as a sidecar in `tauri.conf.json`
3. The frontend's `VITE_API_BASE=/api/v1` points to the sidecar's local port
4. `pnpm tauri build` produces a native installer

---

## Database

The SQLite database is a single file. To back it up, copy it while the server is stopped (or use SQLite's online backup API):

```bash
sqlite3 orgalivro.db ".backup orgalivro-backup.db"
```

To start fresh, delete the file — the backend recreates schema on next startup.
