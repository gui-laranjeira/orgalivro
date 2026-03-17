# Task 02 — Frontend Implementation

Vite + Svelte 5 + TailwindCSS v4 frontend for Orgalivro.

---

## Prerequisites

- Node.js 20+
- pnpm: `npm i -g pnpm`
- Backend running on `:8080`

---

## Step 1 — Scaffold Vite + Svelte 5

```bash
pnpm create vite frontend --template svelte-ts
cd frontend
pnpm install
```

---

## Step 2 — TailwindCSS v4

```bash
pnpm add -D tailwindcss @tailwindcss/vite
```

Create `tailwind.config.ts` (v4 no longer needs a config file, but create it for IDE support):

```ts
// tailwind.config.ts — v4 uses CSS-based config, this file is optional
export default {
  content: ["./src/**/*.{svelte,ts,html}"],
};
```

Add the Tailwind plugin to `vite.config.ts`:

```ts
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
```

In `src/app.css` (or `src/index.css`), replace all contents with:

```css
@import "tailwindcss";
```

Import it in `src/main.ts`:

```ts
import "./app.css";
import App from "./App.svelte";
import { mount } from "svelte";

const app = mount(App, { target: document.getElementById("app")! });
export default app;
```

---

## Step 3 — Router

```bash
pnpm add svelte-spa-router
```

---

## Step 4 — Environment Files

### `.env.development`

```
VITE_API_BASE=/api/v1
```

### `.env.production`

```
VITE_API_BASE=/api/v1
```

> In production the static files and API live behind the same origin. For Tauri, override with `VITE_API_BASE=http://localhost:8080/api/v1`.

---

## Step 5 — `src/lib/types.ts`

```ts
export interface Profile {
  id: number;
  name: string;
  avatar_url: string;
  created_at: string;
}

export interface Author {
  id: number;
  name: string;
}

export interface Genre {
  id: number;
  name: string;
}

export interface Book {
  id: number;
  title: string;
  isbn13: string;
  cover_url: string;
  description: string;
  year: number;
  language: string;
  authors: Author[];
  genres: Genre[];
  created_at: string;
  updated_at: string;
}

export type ReadingStatus = "want_to_read" | "reading" | "read";

export interface LibraryEntry {
  id: number;
  profile_id: number;
  book: Book;
  status: ReadingStatus;
  rating: number | null;
  notes: string;
  added_at: string;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface IsbnResult {
  title: string;
  isbn13: string;
  cover_url: string;
  description: string;
  year: number;
  language: string;
  authors: string[];
  genres: string[];
}
```

---

## Step 6 — `src/lib/api.ts`

```ts
import type {
  Profile,
  Book,
  LibraryEntry,
  PaginatedResponse,
  IsbnResult,
  ReadingStatus,
} from "./types";

const BASE = import.meta.env.VITE_API_BASE as string;

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error ?? res.statusText);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

// --- Profiles ---
export const profiles = {
  list: () => req<Profile[]>("/profiles"),
  create: (name: string, avatar_url?: string) =>
    req<Profile>("/profiles", {
      method: "POST",
      body: JSON.stringify({ name, avatar_url }),
    }),
  delete: (id: number) => req<void>(`/profiles/${id}`, { method: "DELETE" }),
};

// --- Books ---
export interface BookQuery {
  q?: string;
  author?: string;
  genre?: string;
  year?: number;
  language?: string;
  page?: number;
  limit?: number;
}

export const books = {
  list: (query?: BookQuery) => {
    const params = new URLSearchParams();
    if (query) {
      Object.entries(query).forEach(([k, v]) => {
        if (v !== undefined && v !== "") params.set(k, String(v));
      });
    }
    const qs = params.toString();
    return req<PaginatedResponse<Book>>("/books" + (qs ? "?" + qs : ""));
  },
  get: (id: number) => req<Book>(`/books/${id}`),
  create: (data: Partial<Book> & { authors: string[]; genres: string[] }) =>
    req<Book>("/books", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: Partial<Book> & { authors?: string[]; genres?: string[] }) =>
    req<Book>(`/books/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => req<void>(`/books/${id}`, { method: "DELETE" }),
};

// --- Authors & Genres ---
export const authors = {
  list: () => req<{ id: number; name: string }[]>("/authors"),
};
export const genres = {
  list: () => req<{ id: number; name: string }[]>("/genres"),
};

// --- Library entries ---
export interface EntryQuery {
  status?: ReadingStatus;
  q?: string;
  page?: number;
  limit?: number;
}

export const library = {
  list: (profileId: number, query?: EntryQuery) => {
    const params = new URLSearchParams();
    if (query) {
      Object.entries(query).forEach(([k, v]) => {
        if (v !== undefined && v !== "") params.set(k, String(v));
      });
    }
    const qs = params.toString();
    return req<PaginatedResponse<LibraryEntry>>(
      `/profiles/${profileId}/library` + (qs ? "?" + qs : "")
    );
  },
  add: (profileId: number, bookId: number, status?: ReadingStatus, rating?: number, notes?: string) =>
    req<LibraryEntry>(`/profiles/${profileId}/library`, {
      method: "POST",
      body: JSON.stringify({ book_id: bookId, status, rating, notes }),
    }),
  update: (profileId: number, bookId: number, data: { status?: ReadingStatus; rating?: number; notes?: string }) =>
    req<LibraryEntry>(`/profiles/${profileId}/library/${bookId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  remove: (profileId: number, bookId: number) =>
    req<void>(`/profiles/${profileId}/library/${bookId}`, { method: "DELETE" }),
};

// --- ISBN ---
export const isbn = {
  lookup: (code: string) => req<IsbnResult>(`/isbn/${code}`),
};
```

---

## Step 7 — Stores

### `src/lib/stores/profile.ts`

```ts
import { writable } from "svelte/store";
import type { Profile } from "../types";

const STORAGE_KEY = "orgalivro_active_profile";

function createProfileStore() {
  const stored = localStorage.getItem(STORAGE_KEY);
  const initial: Profile | null = stored ? JSON.parse(stored) : null;

  const { subscribe, set } = writable<Profile | null>(initial);

  return {
    subscribe,
    set(profile: Profile | null) {
      if (profile) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(profile));
      } else {
        localStorage.removeItem(STORAGE_KEY);
      }
      set(profile);
    },
  };
}

export const activeProfile = createProfileStore();
```

### `src/lib/stores/toast.ts`

```ts
import { writable } from "svelte/store";

export type ToastType = "success" | "error" | "info";

interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

let _id = 0;

function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([]);

  function show(message: string, type: ToastType = "info") {
    const id = ++_id;
    update((toasts) => [...toasts, { id, message, type }]);
    setTimeout(() => {
      update((toasts) => toasts.filter((t) => t.id !== id));
    }, 3000);
  }

  return {
    subscribe,
    success: (msg: string) => show(msg, "success"),
    error: (msg: string) => show(msg, "error"),
    info: (msg: string) => show(msg, "info"),
  };
}

export const toast = createToastStore();
```

---

## Step 8 — `src/lib/utils.ts`

```ts
export const STATUS_LABELS: Record<string, string> = {
  want_to_read: "Want to Read",
  reading: "Reading",
  read: "Read",
};

export const STATUS_COLORS: Record<string, string> = {
  want_to_read: "bg-gray-100 text-gray-700",
  reading: "bg-blue-100 text-blue-700",
  read: "bg-green-100 text-green-700",
};

export function debounce<T extends (...args: any[]) => any>(fn: T, ms: number): T {
  let timer: ReturnType<typeof setTimeout>;
  return ((...args: any[]) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), ms);
  }) as T;
}

export function starRating(rating: number | null): string {
  if (!rating) return "—";
  return "★".repeat(rating) + "☆".repeat(5 - rating);
}
```

---

## Step 9 — Layout Components

### `src/components/layout/Navbar.svelte`

```svelte
<script lang="ts">
  import { activeProfile } from "$lib/stores/profile";
  import { profiles as profilesApi } from "$lib/api";
  import type { Profile } from "$lib/types";
  import { onMount } from "svelte";

  let allProfiles: Profile[] = $state([]);
  let open = $state(false);

  onMount(async () => {
    allProfiles = await profilesApi.list();
  });

  function select(p: Profile) {
    activeProfile.set(p);
    open = false;
  }
</script>

<nav class="bg-white border-b border-gray-200 px-4 h-14 flex items-center gap-4">
  <a href="#/" class="font-bold text-lg text-indigo-600">Orgalivro</a>

  <div class="flex gap-3 ml-4 text-sm">
    <a href="#/" class="hover:text-indigo-600">Library</a>
    <a href="#/catalog" class="hover:text-indigo-600">Catalog</a>
    <a href="#/books/add" class="hover:text-indigo-600">Add Book</a>
  </div>

  <div class="ml-auto relative">
    <button
      onclick={() => (open = !open)}
      class="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-200 hover:bg-gray-50 text-sm"
    >
      {#if $activeProfile}
        <span>{$activeProfile.name}</span>
      {:else}
        <span class="text-gray-400">Select profile</span>
      {/if}
      <span class="text-gray-400">▾</span>
    </button>

    {#if open}
      <div class="absolute right-0 mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-lg z-50">
        {#each allProfiles as p}
          <button
            onclick={() => select(p)}
            class="w-full text-left px-4 py-2 hover:bg-gray-50 text-sm
              {$activeProfile?.id === p.id ? 'font-semibold text-indigo-600' : ''}"
          >
            {p.name}
          </button>
        {/each}
        <hr class="my-1" />
        <a
          href="#/profiles"
          onclick={() => (open = false)}
          class="block px-4 py-2 text-sm text-gray-500 hover:bg-gray-50"
        >
          Manage profiles
        </a>
      </div>
    {/if}
  </div>
</nav>
```

### `src/components/layout/Toast.svelte`

```svelte
<script lang="ts">
  import { toast } from "$lib/stores/toast";

  const TYPE_CLASSES = {
    success: "bg-green-600",
    error: "bg-red-600",
    info: "bg-gray-800",
  };
</script>

<div class="fixed bottom-4 right-4 flex flex-col gap-2 z-50">
  {#each $toast as t (t.id)}
    <div
      class="px-4 py-2 rounded-lg text-white text-sm shadow-lg {TYPE_CLASSES[t.type]}"
    >
      {t.message}
    </div>
  {/each}
</div>
```

---

## Step 10 — Book Components

### `src/components/book/BookStatusBadge.svelte`

```svelte
<script lang="ts">
  import { STATUS_LABELS, STATUS_COLORS } from "$lib/utils";
  let { status }: { status: string } = $props();
</script>

<span class="px-2 py-0.5 rounded-full text-xs font-medium {STATUS_COLORS[status] ?? ''}">
  {STATUS_LABELS[status] ?? status}
</span>
```

### `src/components/book/BookCard.svelte`

```svelte
<script lang="ts">
  import type { Book } from "$lib/types";
  let { book }: { book: Book } = $props();
</script>

<a href="#/books/{book.id}" class="block bg-white rounded-xl border border-gray-200 hover:shadow-md transition p-3 flex gap-3">
  {#if book.cover_url}
    <img src={book.cover_url} alt={book.title} class="w-16 h-24 object-cover rounded" />
  {:else}
    <div class="w-16 h-24 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs">No cover</div>
  {/if}
  <div class="flex-1 min-w-0">
    <p class="font-semibold text-sm leading-tight line-clamp-2">{book.title}</p>
    {#if book.authors?.length}
      <p class="text-xs text-gray-500 mt-1">{book.authors.map(a => a.name).join(", ")}</p>
    {/if}
    {#if book.year}
      <p class="text-xs text-gray-400 mt-1">{book.year}</p>
    {/if}
  </div>
</a>
```

### `src/components/book/BookForm.svelte`

```svelte
<script lang="ts">
  import type { IsbnResult } from "$lib/types";

  interface Props {
    prefill?: IsbnResult;
    onsubmit: (data: {
      title: string;
      isbn13: string;
      cover_url: string;
      description: string;
      year: number;
      language: string;
      authors: string[];
      genres: string[];
    }) => Promise<void>;
  }

  let { prefill, onsubmit }: Props = $props();

  let title = $state(prefill?.title ?? "");
  let isbn13 = $state(prefill?.isbn13 ?? "");
  let cover_url = $state(prefill?.cover_url ?? "");
  let description = $state(prefill?.description ?? "");
  let year = $state(prefill?.year ?? 0);
  let language = $state(prefill?.language ?? "");
  let authorsInput = $state((prefill?.authors ?? []).join(", "));
  let genresInput = $state((prefill?.genres ?? []).join(", "));
  let submitting = $state(false);
  let error = $state("");

  $effect(() => {
    if (prefill) {
      title = prefill.title ?? "";
      isbn13 = prefill.isbn13 ?? "";
      cover_url = prefill.cover_url ?? "";
      description = prefill.description ?? "";
      year = prefill.year ?? 0;
      language = prefill.language ?? "";
      authorsInput = (prefill.authors ?? []).join(", ");
      genresInput = (prefill.genres ?? []).join(", ");
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    error = "";
    submitting = true;
    try {
      await onsubmit({
        title,
        isbn13,
        cover_url,
        description,
        year: Number(year),
        language,
        authors: authorsInput.split(",").map(s => s.trim()).filter(Boolean),
        genres: genresInput.split(",").map(s => s.trim()).filter(Boolean),
      });
    } catch (err: any) {
      error = err.message;
    } finally {
      submitting = false;
    }
  }
</script>

<form onsubmit={handleSubmit} class="space-y-4">
  {#if error}
    <p class="text-red-600 text-sm">{error}</p>
  {/if}
  <div>
    <label class="block text-sm font-medium mb-1">Title *</label>
    <input bind:value={title} required class="w-full border rounded-lg px-3 py-2 text-sm" />
  </div>
  <div>
    <label class="block text-sm font-medium mb-1">Authors (comma-separated)</label>
    <input bind:value={authorsInput} class="w-full border rounded-lg px-3 py-2 text-sm" />
  </div>
  <div class="grid grid-cols-2 gap-3">
    <div>
      <label class="block text-sm font-medium mb-1">ISBN-13</label>
      <input bind:value={isbn13} class="w-full border rounded-lg px-3 py-2 text-sm" />
    </div>
    <div>
      <label class="block text-sm font-medium mb-1">Year</label>
      <input type="number" bind:value={year} class="w-full border rounded-lg px-3 py-2 text-sm" />
    </div>
  </div>
  <div>
    <label class="block text-sm font-medium mb-1">Cover URL</label>
    <input bind:value={cover_url} class="w-full border rounded-lg px-3 py-2 text-sm" />
  </div>
  <div>
    <label class="block text-sm font-medium mb-1">Language</label>
    <input bind:value={language} placeholder="en" class="w-full border rounded-lg px-3 py-2 text-sm" />
  </div>
  <div>
    <label class="block text-sm font-medium mb-1">Genres (comma-separated)</label>
    <input bind:value={genresInput} class="w-full border rounded-lg px-3 py-2 text-sm" />
  </div>
  <div>
    <label class="block text-sm font-medium mb-1">Description</label>
    <textarea bind:value={description} rows="4" class="w-full border rounded-lg px-3 py-2 text-sm"></textarea>
  </div>
  <button
    type="submit"
    disabled={submitting}
    class="w-full bg-indigo-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50"
  >
    {submitting ? "Saving…" : "Save Book"}
  </button>
</form>
```

### `src/components/book/BookDetail.svelte`

```svelte
<script lang="ts">
  import type { Book, LibraryEntry, ReadingStatus } from "$lib/types";
  import BookStatusBadge from "./BookStatusBadge.svelte";
  import { starRating } from "$lib/utils";

  let { book, entry, onAddToLibrary, onUpdateEntry }: {
    book: Book;
    entry: LibraryEntry | null;
    onAddToLibrary: (status: ReadingStatus) => Promise<void>;
    onUpdateEntry: (data: { status?: ReadingStatus; rating?: number; notes?: string }) => Promise<void>;
  } = $props();
</script>

<div class="flex gap-6">
  {#if book.cover_url}
    <img src={book.cover_url} alt={book.title} class="w-32 h-48 object-cover rounded-lg shadow" />
  {/if}
  <div class="flex-1">
    <h1 class="text-2xl font-bold">{book.title}</h1>
    {#if book.authors?.length}
      <p class="text-gray-600 mt-1">{book.authors.map(a => a.name).join(", ")}</p>
    {/if}
    <div class="flex gap-2 mt-2 flex-wrap">
      {#if book.year}<span class="text-sm text-gray-500">{book.year}</span>{/if}
      {#if book.language}<span class="text-sm text-gray-500">· {book.language.toUpperCase()}</span>{/if}
      {#if book.isbn13}<span class="text-sm text-gray-500">· ISBN {book.isbn13}</span>{/if}
    </div>
    {#if book.genres?.length}
      <div class="flex gap-1 mt-2 flex-wrap">
        {#each book.genres as g}
          <span class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded text-xs">{g.name}</span>
        {/each}
      </div>
    {/if}

    <div class="mt-4 space-y-2">
      {#if entry}
        <div class="flex items-center gap-2">
          <BookStatusBadge status={entry.status} />
          {#if entry.rating}
            <span class="text-sm">{starRating(entry.rating)}</span>
          {/if}
        </div>
        <select
          value={entry.status}
          onchange={(e) => onUpdateEntry({ status: (e.target as HTMLSelectElement).value as ReadingStatus })}
          class="border rounded px-2 py-1 text-sm"
        >
          <option value="want_to_read">Want to Read</option>
          <option value="reading">Reading</option>
          <option value="read">Read</option>
        </select>
      {:else}
        <p class="text-sm text-gray-500">Not in your library</p>
        <div class="flex gap-2">
          <button onclick={() => onAddToLibrary("want_to_read")} class="text-sm px-3 py-1.5 border rounded-lg hover:bg-gray-50">Want to Read</button>
          <button onclick={() => onAddToLibrary("reading")} class="text-sm px-3 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700">Reading</button>
          <button onclick={() => onAddToLibrary("read")} class="text-sm px-3 py-1.5 bg-green-600 text-white rounded-lg hover:bg-green-700">Read</button>
        </div>
      {/if}
    </div>
  </div>
</div>
{#if book.description}
  <p class="mt-6 text-sm text-gray-700 leading-relaxed">{book.description}</p>
{/if}
```

---

## Step 11 — Library Components

### `src/components/library/LibraryFilters.svelte`

```svelte
<script lang="ts">
  import { debounce } from "$lib/utils";

  let { onFilter }: {
    onFilter: (filters: { status?: string; q?: string }) => void;
  } = $props();

  let status = $state("");
  let q = $state("");

  const emitDebounced = debounce(() => {
    onFilter({ status: status || undefined, q: q || undefined });
  }, 300);

  $effect(() => {
    status; q; // track
    emitDebounced();
  });
</script>

<div class="flex gap-3 items-center">
  <input
    bind:value={q}
    placeholder="Search books…"
    class="border rounded-lg px-3 py-1.5 text-sm flex-1"
  />
  <select bind:value={status} class="border rounded-lg px-2 py-1.5 text-sm">
    <option value="">All statuses</option>
    <option value="want_to_read">Want to Read</option>
    <option value="reading">Reading</option>
    <option value="read">Read</option>
  </select>
</div>
```

### `src/components/library/EntryCard.svelte`

```svelte
<script lang="ts">
  import type { LibraryEntry, ReadingStatus } from "$lib/types";
  import BookStatusBadge from "../book/BookStatusBadge.svelte";
  import { starRating } from "$lib/utils";

  let { entry, onUpdate, onRemove }: {
    entry: LibraryEntry;
    onUpdate: (data: { status?: ReadingStatus }) => Promise<void>;
    onRemove: () => Promise<void>;
  } = $props();
</script>

<div class="bg-white border border-gray-200 rounded-xl p-3 flex gap-3">
  <a href="#/books/{entry.book.id}">
    {#if entry.book.cover_url}
      <img src={entry.book.cover_url} alt={entry.book.title} class="w-14 h-20 object-cover rounded" />
    {:else}
      <div class="w-14 h-20 bg-gray-100 rounded"></div>
    {/if}
  </a>
  <div class="flex-1 min-w-0">
    <a href="#/books/{entry.book.id}" class="font-semibold text-sm line-clamp-2 hover:text-indigo-600">
      {entry.book.title}
    </a>
    {#if entry.book.authors?.length}
      <p class="text-xs text-gray-500 mt-0.5">{entry.book.authors.map(a => a.name).join(", ")}</p>
    {/if}
    <div class="flex items-center gap-2 mt-2 flex-wrap">
      <BookStatusBadge status={entry.status} />
      {#if entry.rating}
        <span class="text-xs">{starRating(entry.rating)}</span>
      {/if}
    </div>
    <div class="flex gap-2 mt-2">
      <select
        value={entry.status}
        onchange={(e) => onUpdate({ status: (e.target as HTMLSelectElement).value as ReadingStatus })}
        class="border rounded px-2 py-0.5 text-xs"
      >
        <option value="want_to_read">Want to Read</option>
        <option value="reading">Reading</option>
        <option value="read">Read</option>
      </select>
      <button onclick={onRemove} class="text-xs text-red-400 hover:text-red-600">Remove</button>
    </div>
  </div>
</div>
```

---

## Step 12 — ISBN Lookup Component

### `src/components/isbn/IsbnLookup.svelte`

```svelte
<script lang="ts">
  import { isbn as isbnApi } from "$lib/api";
  import type { IsbnResult } from "$lib/types";

  let { onResult }: { onResult: (result: IsbnResult) => void } = $props();

  let value = $state("");
  let loading = $state(false);
  let error = $state("");

  async function lookup() {
    if (!value.trim()) return;
    loading = true;
    error = "";
    try {
      const result = await isbnApi.lookup(value.trim());
      onResult(result);
      value = "";
    } catch (err: any) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>

<div class="flex gap-2 items-start">
  <div class="flex-1">
    <div class="flex gap-2">
      <input
        bind:value
        placeholder="ISBN-13 (e.g. 9780261102354)"
        class="flex-1 border rounded-lg px-3 py-2 text-sm"
        onkeydown={(e) => e.key === "Enter" && lookup()}
      />
      <button
        onclick={lookup}
        disabled={loading}
        class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 disabled:opacity-50"
      >
        {loading ? "…" : "Lookup"}
      </button>
    </div>
    {#if error}
      <p class="text-red-500 text-xs mt-1">{error}</p>
    {/if}
  </div>
</div>
```

---

## Step 13 — Profile Components

### `src/components/profile/ProfileForm.svelte`

```svelte
<script lang="ts">
  let { oncreate }: { oncreate: (name: string) => Promise<void> } = $props();

  let name = $state("");
  let error = $state("");
  let loading = $state(false);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim()) return;
    loading = true;
    error = "";
    try {
      await oncreate(name.trim());
      name = "";
    } catch (err: any) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>

<form onsubmit={handleSubmit} class="flex gap-2 items-start">
  <div class="flex-1">
    <input
      bind:value={name}
      placeholder="Profile name"
      required
      class="w-full border rounded-lg px-3 py-2 text-sm"
    />
    {#if error}
      <p class="text-red-500 text-xs mt-1">{error}</p>
    {/if}
  </div>
  <button
    type="submit"
    disabled={loading}
    class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 disabled:opacity-50"
  >
    {loading ? "Creating…" : "Create"}
  </button>
</form>
```

### `src/components/profile/ProfileList.svelte`

```svelte
<script lang="ts">
  import type { Profile } from "$lib/types";
  import { activeProfile } from "$lib/stores/profile";

  let { profiles, onDelete, onSelect }: {
    profiles: Profile[];
    onDelete: (id: number) => Promise<void>;
    onSelect: (p: Profile) => void;
  } = $props();
</script>

<ul class="space-y-2">
  {#each profiles as p (p.id)}
    <li class="flex items-center gap-3 bg-white border border-gray-200 rounded-xl px-4 py-3">
      <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-700 font-bold text-sm">
        {p.name[0].toUpperCase()}
      </div>
      <span class="flex-1 font-medium text-sm">{p.name}</span>
      {#if $activeProfile?.id === p.id}
        <span class="text-xs text-indigo-600 font-medium">Active</span>
      {:else}
        <button onclick={() => onSelect(p)} class="text-xs text-gray-500 hover:text-indigo-600">Set active</button>
      {/if}
      <button onclick={() => onDelete(p.id)} class="text-xs text-red-400 hover:text-red-600 ml-2">Delete</button>
    </li>
  {/each}
</ul>
```

---

## Step 14 — Pages

### `src/pages/ProfilesPage.svelte`

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { profiles as profilesApi } from "$lib/api";
  import type { Profile } from "$lib/types";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import ProfileForm from "$lib/components/profile/ProfileForm.svelte";
  import ProfileList from "$lib/components/profile/ProfileList.svelte";

  let profileList: Profile[] = $state([]);

  onMount(async () => {
    profileList = await profilesApi.list();
  });

  async function create(name: string) {
    const p = await profilesApi.create(name);
    profileList = [...profileList, p];
    toast.success(`Profile "${p.name}" created`);
  }

  async function del(id: number) {
    await profilesApi.delete(id);
    profileList = profileList.filter(p => p.id !== id);
    if ($activeProfile?.id === id) activeProfile.set(null);
    toast.info("Profile deleted");
  }

  function select(p: Profile) {
    activeProfile.set(p);
    toast.success(`Switched to ${p.name}`);
  }
</script>

<div class="max-w-xl mx-auto py-8 px-4 space-y-6">
  <h1 class="text-2xl font-bold">Profiles</h1>
  <ProfileForm oncreate={create} />
  <ProfileList profiles={profileList} onDelete={del} onSelect={select} />
</div>
```

### `src/pages/AddBookPage.svelte`

```svelte
<script lang="ts">
  import { push } from "svelte-spa-router";
  import { books as booksApi, library } from "$lib/api";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import type { IsbnResult } from "$lib/types";
  import IsbnLookup from "$lib/components/isbn/IsbnLookup.svelte";
  import BookForm from "$lib/components/book/BookForm.svelte";

  let prefill: IsbnResult | undefined = $state(undefined);

  async function handleSubmit(data: any) {
    const book = await booksApi.create(data);
    if ($activeProfile) {
      await library.add($activeProfile.id, book.id, "want_to_read");
    }
    toast.success("Book added!");
    push(`/books/${book.id}`);
  }
</script>

<div class="max-w-xl mx-auto py-8 px-4 space-y-6">
  <h1 class="text-2xl font-bold">Add Book</h1>
  <div class="bg-white border border-gray-200 rounded-xl p-4">
    <p class="text-sm font-medium mb-2">Lookup by ISBN</p>
    <IsbnLookup onResult={(r) => (prefill = r)} />
  </div>
  <div class="bg-white border border-gray-200 rounded-xl p-4">
    <BookForm {prefill} onsubmit={handleSubmit} />
  </div>
</div>
```

### `src/pages/CatalogPage.svelte`

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { books as booksApi } from "$lib/api";
  import type { Book } from "$lib/types";
  import { debounce } from "$lib/utils";
  import BookCard from "$lib/components/book/BookCard.svelte";

  let bookList: Book[] = $state([]);
  let total = $state(0);
  let q = $state("");
  let loading = $state(false);

  const loadBooks = debounce(async () => {
    loading = true;
    try {
      const res = await booksApi.list({ q: q || undefined });
      bookList = res.data;
      total = res.total;
    } finally {
      loading = false;
    }
  }, 300);

  onMount(() => loadBooks());
  $effect(() => { q; loadBooks(); });
</script>

<div class="max-w-3xl mx-auto py-8 px-4 space-y-4">
  <div class="flex items-center gap-3">
    <h1 class="text-2xl font-bold flex-1">Catalog</h1>
    <span class="text-sm text-gray-500">{total} books</span>
    <a href="#/books/add" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700">+ Add</a>
  </div>
  <input bind:value={q} placeholder="Search…" class="w-full border rounded-lg px-3 py-2 text-sm" />
  {#if loading}
    <p class="text-center text-gray-400 py-8">Loading…</p>
  {:else if bookList.length === 0}
    <p class="text-center text-gray-400 py-8">No books found.</p>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      {#each bookList as book (book.id)}
        <BookCard {book} />
      {/each}
    </div>
  {/if}
</div>
```

### `src/pages/LibraryPage.svelte`

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { library } from "$lib/api";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import type { LibraryEntry, ReadingStatus } from "$lib/types";
  import LibraryFilters from "$lib/components/library/LibraryFilters.svelte";
  import EntryCard from "$lib/components/library/EntryCard.svelte";

  let entries: LibraryEntry[] = $state([]);
  let total = $state(0);
  let filters: { status?: string; q?: string } = $state({});

  $effect(() => {
    if (!$activeProfile) {
      push("/profiles");
      return;
    }
    load();
  });

  async function load() {
    if (!$activeProfile) return;
    const res = await library.list($activeProfile.id, { status: filters.status as ReadingStatus, q: filters.q });
    entries = res.data;
    total = res.total;
  }

  async function update(entry: LibraryEntry, data: { status?: ReadingStatus }) {
    if (!$activeProfile) return;
    await library.update($activeProfile.id, entry.book.id, data);
    await load();
  }

  async function remove(entry: LibraryEntry) {
    if (!$activeProfile) return;
    await library.remove($activeProfile.id, entry.book.id);
    toast.info("Removed from library");
    await load();
  }
</script>

<div class="max-w-3xl mx-auto py-8 px-4 space-y-4">
  <div class="flex items-center gap-2">
    <h1 class="text-2xl font-bold flex-1">{$activeProfile?.name}'s Library</h1>
    <span class="text-sm text-gray-500">{total} books</span>
  </div>
  <LibraryFilters onFilter={(f) => { filters = f; load(); }} />
  {#if entries.length === 0}
    <p class="text-center text-gray-400 py-12">No books yet. <a href="#/books/add" class="text-indigo-600">Add one!</a></p>
  {:else}
    <div class="space-y-3">
      {#each entries as entry (entry.id)}
        <EntryCard
          {entry}
          onUpdate={(data) => update(entry, data)}
          onRemove={() => remove(entry)}
        />
      {/each}
    </div>
  {/if}
</div>
```

### `src/pages/BookPage.svelte`

```svelte
<script lang="ts">
  import { onMount } from "svelte";
  import { books as booksApi, library } from "$lib/api";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import type { Book, LibraryEntry, ReadingStatus } from "$lib/types";
  import BookDetail from "$lib/components/book/BookDetail.svelte";

  export let params: { id: string };

  let book: Book | null = $state(null);
  let entry: LibraryEntry | null = $state(null);

  onMount(async () => {
    book = await booksApi.get(Number(params.id));
    if ($activeProfile) {
      const res = await library.list($activeProfile.id, {});
      entry = res.data.find(e => e.book.id === book!.id) ?? null;
    }
  });

  async function addToLibrary(status: ReadingStatus) {
    if (!$activeProfile || !book) return;
    entry = await library.add($activeProfile.id, book.id, status);
    toast.success("Added to library!");
  }

  async function updateEntry(data: { status?: ReadingStatus; rating?: number; notes?: string }) {
    if (!$activeProfile || !book) return;
    entry = await library.update($activeProfile.id, book.id, data);
    toast.success("Updated!");
  }
</script>

<div class="max-w-3xl mx-auto py-8 px-4">
  {#if !book}
    <p class="text-center text-gray-400 py-12">Loading…</p>
  {:else}
    <BookDetail {book} {entry} onAddToLibrary={addToLibrary} onUpdateEntry={updateEntry} />
  {/if}
</div>
```

---

## Step 15 — `src/App.svelte`

```svelte
<script lang="ts">
  import Router from "svelte-spa-router";
  import Navbar from "./components/layout/Navbar.svelte";
  import Toast from "./components/layout/Toast.svelte";
  import LibraryPage from "./pages/LibraryPage.svelte";
  import CatalogPage from "./pages/CatalogPage.svelte";
  import BookPage from "./pages/BookPage.svelte";
  import AddBookPage from "./pages/AddBookPage.svelte";
  import ProfilesPage from "./pages/ProfilesPage.svelte";

  const routes = {
    "/": LibraryPage,
    "/catalog": CatalogPage,
    "/books/add": AddBookPage,
    "/books/:id": BookPage,
    "/profiles": ProfilesPage,
  };
</script>

<div class="min-h-screen bg-gray-50">
  <Navbar />
  <main>
    <Router {routes} />
  </main>
  <Toast />
</div>
```

---

## Step 16 — Vite alias config (optional but convenient)

Add `resolve.alias` to `vite.config.ts` so `$lib` resolves to `src/lib`:

```ts
import { resolve } from "path";

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  resolve: {
    alias: {
      "$lib": resolve(__dirname, "./src/lib"),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
```

Also update `tsconfig.json` to include the path mapping:

```json
{
  "compilerOptions": {
    "paths": {
      "$lib": ["./src/lib"],
      "$lib/*": ["./src/lib/*"]
    }
  }
}
```

---

## Development Order

Implement in this order to validate incrementally:

1. `types.ts` → `api.ts` → `utils.ts`
2. Stores (`profile.ts`, `toast.ts`)
3. Navbar + Toast (layout shell)
4. ProfilesPage (end-to-end: create/select/delete profiles)
5. AddBookPage + IsbnLookup + BookForm (add first book)
6. CatalogPage + BookCard
7. LibraryPage + LibraryFilters + EntryCard
8. BookPage + BookDetail

---

## Running the Dev Server

```bash
cd frontend
pnpm dev
# → http://localhost:5173
```

Make sure the backend is running on `:8080` first.

---

## Build for Production

```bash
pnpm build        # outputs to dist/
pnpm preview      # preview the production build
```

The `dist/` folder is a fully static site. Serve it alongside the Go binary and set `ALLOWED_ORIGINS` appropriately.
