# Orgalivro — System Specification

## Overview

Orgalivro is a personal book library management system designed for shared use on a single device. Multiple named profiles (e.g., family members, roommates) each maintain their own reading list and status, while sharing a single physical book catalog. No passwords or authentication — profiles are selected by name and persisted in the browser's localStorage.

---

## Goals

- Track books owned/read/wanted across multiple named profiles on one device
- Zero recurring cost (SQLite on disk, no cloud services required)
- ISBN-13 barcode scanning support (future — mobile browser or Tauri)
- Start as a local web app, evolve to a cross-platform desktop app (Tauri v2)
- Fast to use: scan an ISBN and a book is prefilled in seconds

## Non-Goals

- Multi-device sync or cloud backup (out of scope for v1)
- User authentication, passwords, or access control
- Social features (sharing lists, recommendations)
- eBook management (physical books only)
- Public deployment / multi-tenant hosting

---

## Tech Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Frontend | Vite + Svelte 5 + TypeScript | Fast builds, excellent reactivity model, small bundle |
| Styling | TailwindCSS v4 | Utility-first, no CSS files to maintain |
| Routing | svelte-spa-router (hash routing) | Works in Tauri WebView and static file servers without server rewrites |
| Backend | Go + Gin | Lightweight HTTP, easy cross-compilation |
| ORM | GORM | AutoMigrate, type-safe queries |
| Database | SQLite via glebarez/sqlite (pure-Go) | CGO-free → straightforward Tauri sidecar compilation |
| ISBN Lookup | Open Library API (free, no key needed) | Sufficient metadata, no rate limits for personal use |

---

## Data Model

### Rationale

Books are stored once in a shared catalog (`books` table), keyed by ISBN-13. Each profile maintains its own reading state via `library_entries` (status, rating, notes). This means:
- No duplicate rows for the same physical book
- Any profile can view any book's metadata
- Each profile's reading history is independent

### Entity Relationship

```
profiles (1) ──< library_entries >── (1) books
books >──< book_authors >──< authors
books >──< book_genres  >──< genres
```

### Schema

```sql
CREATE TABLE profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    avatar_url TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    isbn13 TEXT UNIQUE,
    cover_url TEXT,
    description TEXT,
    year INTEGER,
    language TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE authors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE book_authors (
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE genres (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE book_genres (
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, genre_id)
);

CREATE TABLE library_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK(status IN ('want_to_read','reading','read')) DEFAULT 'want_to_read',
    rating INTEGER CHECK(rating IS NULL OR (rating >= 1 AND rating <= 5)),
    notes TEXT,
    added_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(profile_id, book_id)
);
```

---

## REST API

Base path: `/api/v1`

All responses use JSON. Error responses: `{"error": "message"}`.

### Profiles

| Method | Path | Body | Response | Notes |
|--------|------|------|----------|-------|
| GET | /profiles | — | `Profile[]` | All profiles |
| POST | /profiles | `{name, avatar_url?}` | `Profile` 201 | Name must be unique |
| DELETE | /profiles/:id | — | 204 | Cascades library entries |

### Books (shared catalog)

| Method | Path | Body / Query | Response | Notes |
|--------|------|------|----------|-------|
| GET | /books | `?q=&author=&genre=&year=&language=&page=&limit=` | `{data: Book[], total, page, limit}` | Full-text search on title/description |
| POST | /books | `{title, isbn13?, cover_url?, description?, year?, language?, authors[], genres[]}` | `Book` 201 | Authors/genres auto-created (FirstOrCreate) |
| GET | /books/:id | — | `Book` | Includes authors, genres |
| PUT | /books/:id | same as POST body (partial) | `Book` | Updates metadata |
| DELETE | /books/:id | — | 204 | Only deletes if no library entries reference it; else 409 |

### Authors & Genres

| Method | Path | Response |
|--------|------|----------|
| GET | /authors | `Author[]` (all distinct) |
| GET | /genres | `Genre[]` (all distinct) |

### Library Entries (per-profile)

| Method | Path | Body / Query | Response | Notes |
|--------|------|------|----------|-------|
| GET | /profiles/:pid/library | `?status=&q=&page=&limit=` | `{data: Entry[], total, page, limit}` | Entries with book details |
| POST | /profiles/:pid/library | `{book_id, status?, rating?, notes?}` | `Entry` 201 | Creates entry; 409 if duplicate |
| PUT | /profiles/:pid/library/:book_id | `{status?, rating?, notes?}` | `Entry` | Partial update |
| DELETE | /profiles/:pid/library/:book_id | — | 204 | Removes from profile only |

### ISBN Lookup

| Method | Path | Response |
|--------|------|----------|
| GET | /isbn/:isbn | `IsbnResult` (normalized from Open Library) |

`IsbnResult`:
```json
{
  "title": "The Lord of the Rings",
  "isbn13": "9780261102354",
  "cover_url": "https://covers.openlibrary.org/b/isbn/9780261102354-L.jpg",
  "description": "...",
  "year": 1968,
  "language": "en",
  "authors": ["J.R.R. Tolkien"],
  "genres": []
}
```

---

## ISBN Lookup Flow

1. User types or scans an ISBN-13 into the `IsbnLookup` component
2. Component calls `GET /api/v1/isbn/:isbn` on the backend
3. Backend fetches `https://openlibrary.org/api/books?bibkeys=ISBN:{isbn}&format=json&jscmd=data`
4. Backend normalizes the response into `IsbnResult` and returns it
5. Frontend pre-fills the `BookForm` with the result
6. User reviews/edits and submits to `POST /api/v1/books` then `POST /api/v1/profiles/:pid/library`

Proxying through the backend avoids CORS restrictions from Open Library in the browser.

---

## Configuration

The backend reads configuration from environment variables with defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PATH` | `./orgalivro.db` | SQLite database file path |
| `PORT` | `8080` | HTTP server port |
| `ALLOWED_ORIGINS` | `http://localhost:5173` | CORS allowed origins (comma-separated) |

---

## CORS Strategy

- In development: Vite proxy forwards `/api` requests to `:8080` — no CORS headers needed for the browser
- In production (web): Backend configured with `ALLOWED_ORIGINS` pointing to the static file server origin
- In Tauri: WebView makes requests to the Go sidecar on `localhost`; origin is `tauri://localhost` — add to `ALLOWED_ORIGINS`

---

## Frontend State Management

### Stores

- **`activeProfile`** (`src/lib/stores/profile.ts`): `Profile | null`, persisted to `localStorage` under key `orgalivro_active_profile`. Reading from `localStorage` on store initialization.
- **`toast`** (`src/lib/stores/toast.ts`): Array of `{id, message, type}`. Each toast auto-dismisses after 3 seconds.

### Routing

Hash-based routing via `svelte-spa-router`:

| Route | Page | Description |
|-------|------|-------------|
| `#/` | LibraryPage | Active profile's library (redirects to /profiles if none set) |
| `#/catalog` | CatalogPage | All books in the shared catalog |
| `#/books/:id` | BookPage | Book detail + add to library |
| `#/books/add` | AddBookPage | Add new book (with ISBN lookup) |
| `#/profiles` | ProfilesPage | List and create profiles |

---

## SQLite Pragmas & Connection Settings

Applied immediately after opening the database:

```sql
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;
```

`SetMaxOpenConns(1)` is set on the underlying `*sql.DB` because SQLite allows only one writer at a time. WAL mode allows concurrent readers during a write.

---

## Go Dependency Wiring

Manual wiring in `main.go` (no DI framework):

```
main → config.Load()
     → db.Open(cfg.DBPath)
     → db.Migrate(gormDB)
     → repositories (profileRepo, bookRepo, entryRepo)
     → services (profileSvc, bookSvc, entrySvc, isbnSvc)
     → handlers (profileHandler, bookHandler, entryHandler, isbnHandler)
     → router.New(cfg, handlers...)
     → r.Run(":" + cfg.Port)
```

---

## Future Tauri Integration

1. Compile Go backend as a Tauri sidecar binary (`orgalivro-server`)
2. Tauri spawns the sidecar on startup and waits for it to be ready
3. Frontend's `VITE_API_BASE` points to `http://localhost:8080/api/v1`
4. Add `tauri://localhost` to `ALLOWED_ORIGINS`
5. ISBN barcode scanning: use `@tauri-apps/plugin-barcode-scanner` or mobile browser `BarcodeDetector` API

No frontend changes required for Tauri — hash routing and the API abstraction layer mean the same build works in both contexts.

---

## Project Structure

```
orgalivro/
├── spec.md
├── task_01.md
├── task_02.md
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── db/db.go
│   │   ├── db/migrate.go
│   │   ├── model/profile.go
│   │   ├── model/book.go
│   │   ├── model/entry.go
│   │   ├── repository/profile_repo.go
│   │   ├── repository/book_repo.go
│   │   ├── repository/entry_repo.go
│   │   ├── service/profile_service.go
│   │   ├── service/book_service.go
│   │   ├── service/entry_service.go
│   │   ├── service/isbn_service.go
│   │   ├── handler/profile_handler.go
│   │   ├── handler/book_handler.go
│   │   ├── handler/entry_handler.go
│   │   ├── handler/isbn_handler.go
│   │   ├── router/router.go
│   │   └── dto/
│   │       ├── profile_dto.go
│   │       ├── book_dto.go
│   │       └── entry_dto.go
│   ├── go.mod
│   └── .air.toml
└── frontend/
    ├── src/
    │   ├── main.ts
    │   ├── App.svelte
    │   ├── lib/
    │   │   ├── api.ts
    │   │   ├── types.ts
    │   │   ├── utils.ts
    │   │   └── stores/
    │   │       ├── profile.ts
    │   │       └── toast.ts
    │   ├── components/
    │   │   ├── layout/Navbar.svelte
    │   │   ├── layout/Toast.svelte
    │   │   ├── book/BookCard.svelte
    │   │   ├── book/BookDetail.svelte
    │   │   ├── book/BookForm.svelte
    │   │   ├── book/BookStatusBadge.svelte
    │   │   ├── library/LibraryFilters.svelte
    │   │   ├── library/EntryCard.svelte
    │   │   ├── isbn/IsbnLookup.svelte
    │   │   ├── profile/ProfileList.svelte
    │   │   └── profile/ProfileForm.svelte
    │   └── pages/
    │       ├── LibraryPage.svelte
    │       ├── CatalogPage.svelte
    │       ├── BookPage.svelte
    │       ├── AddBookPage.svelte
    │       └── ProfilesPage.svelte
    ├── vite.config.ts
    ├── tailwind.config.ts
    ├── package.json
    ├── .env.development
    └── .env.production
```
