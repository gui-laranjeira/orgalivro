# Orgalivro

A personal book library management system designed for shared household use. Multiple named profiles can coexist on a single device — no accounts, no passwords. Each profile owns the books they add and can independently track reading status, ratings, and notes for any book in the shared catalog.

---

## Features

- **Named profiles** — switch between readers instantly; active profile persisted locally
- **Book ownership** — the profile that adds a book is its owner; ownership can be transferred
- **Shared catalog** — all books are visible to everyone in the Library view
- **Per-profile reading state** — each profile independently tracks status (Want to Read / Reading / Read), star rating (1–5), and personal notes
- **ISBN-13 lookup** — auto-fills book metadata via BrasilAPI (Brazilian books) with Open Library as fallback
- **Category filtering** — filter the Library by genre
- **Reader tags** — book cards show which profiles are reading each book and their status

---

## Architecture

```
orgalivro/
├── backend/       Go HTTP API
└── frontend/      Svelte SPA
```

### Backend

**Stack:** Go 1.22 · Gin · GORM · SQLite (CGO-free via `glebarez/sqlite`)

```
backend/
├── cmd/server/main.go          Entry point — manual DI wiring
└── internal/
    ├── config/                 Env var configuration
    ├── db/                     SQLite open + WAL pragma + AutoMigrate
    ├── model/                  GORM models (Profile, Book, Author, Genre, LibraryEntry)
    ├── dto/                    Request/response structs
    ├── repository/             DB queries (BookRepo, EntryRepo, ProfileRepo)
    ├── service/                Business logic (BookService, EntryService, ISBNService…)
    ├── handler/                HTTP handlers (Gin)
    └── router/                 Route registration + CORS
```

**Layering:** `Handler → Service → Repository → GORM → SQLite`

**Database design:**

| Table | Purpose |
|---|---|
| `profiles` | Named user profiles |
| `books` | Shared catalog — one row per ISBN, with `owner_profile_id` |
| `authors` | Deduplicated author names |
| `genres` | Deduplicated genre/category names |
| `book_authors` | Many-to-many join |
| `book_genres` | Many-to-many join |
| `library_entries` | Per-profile reading state for a book (status, rating, notes) |

**Key decisions:**
- `PRAGMA journal_mode=WAL` + `SetMaxOpenConns(1)` — safe single-writer SQLite
- `PRAGMA foreign_keys=ON` re-enabled after each migration run
- CGO-free SQLite driver — required for future Tauri sidecar packaging
- ISBN lookup proxied through backend — avoids CORS issues in browser/WebView

### Frontend

**Stack:** Svelte 5 (runes) · Vite 5 · TailwindCSS v4 · TypeScript · svelte-spa-router

```
frontend/src/
├── App.svelte                  Router setup (hash-based for Tauri compatibility)
├── pages/                      LibraryPage, CatalogPage, BookPage, AddBookPage, ProfilesPage
├── components/
│   ├── book/                   BookCard, BookDetail, BookForm, BookStatusBadge
│   ├── isbn/                   IsbnLookup
│   ├── library/                EntryCard, LibraryFilters
│   ├── layout/                 Navbar, Toast
│   └── profile/                ProfileForm, ProfileList
└── lib/
    ├── api.ts                  Typed fetch wrappers grouped by resource
    ├── types.ts                TypeScript interfaces
    ├── utils.ts                STATUS_LABELS, STATUS_COLORS, debounce, starRating
    └── stores/                 profile.ts (localStorage), toast.ts (auto-dismiss)
```

**Key decisions:**
- Hash routing (`#/`) — works in Tauri WebView without server rewrites
- Vite proxy `/api → http://localhost:8080` — identical paths in dev and prod
- Svelte 5 runes — `$state`, `$derived`, `$effect`, `$props` throughout

### REST API

| Method | Path | Description |
|---|---|---|
| GET | `/api/v1/profiles` | List profiles |
| POST | `/api/v1/profiles` | Create profile |
| DELETE | `/api/v1/profiles/:id` | Delete profile |
| GET | `/api/v1/books` | List books (supports `q`, `genre`, `owner_profile_id` filters) |
| POST | `/api/v1/books` | Create book |
| GET | `/api/v1/books/:id` | Get book |
| PUT | `/api/v1/books/:id` | Update book |
| DELETE | `/api/v1/books/:id` | Delete book |
| PUT | `/api/v1/books/:id/owner` | Transfer book ownership |
| GET | `/api/v1/authors` | List all authors |
| GET | `/api/v1/genres` | List all genres |
| GET | `/api/v1/profiles/:id/library` | List a profile's reading entries |
| POST | `/api/v1/profiles/:id/library` | Add book to profile's reading list |
| PUT | `/api/v1/profiles/:id/library/:book_id` | Update reading entry |
| DELETE | `/api/v1/profiles/:id/library/:book_id` | Remove reading entry |
| GET | `/api/v1/isbn/:isbn` | ISBN-13 lookup |

---

## Current Status

The core application is fully functional and usable:

- [x] Profile management (create, switch, delete)
- [x] Book catalog with ISBN lookup, manual entry, and genre filtering
- [x] Book ownership model with transfer support
- [x] Per-profile reading entries (status, rating, notes)
- [x] BrasilAPI ISBN lookup with Open Library fallback
- [x] Reader tags on book cards
- [x] Responsive UI with TailwindCSS v4

### What's left

- [ ] **Tauri packaging** — wrap as a desktop app with Go binary as sidecar
- [ ] **Cover images** — BrasilAPI often returns `null` for `cover_url`; integrate a dedicated cover CDN (e.g. Google Books covers or Open Library covers by ISBN)
- [ ] **Book editing UI** — the backend supports `PUT /books/:id` but there is no edit form on the frontend yet
- [ ] **Delete book** — backend enforces no entries exist before deletion; frontend has no delete button
- [ ] **Pagination** — API supports `page`/`limit` but the UI always loads the first 20 results
- [ ] **Profile avatars** — the model has `avatar_url` but the UI only shows the first-letter avatar
- [ ] **Reading statistics** — books read per profile, average rating, reading streak, etc.
- [ ] **Export / backup** — export a profile's library to CSV or JSON
