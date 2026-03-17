import type {
  Profile,
  Book,
  LibraryEntry,
  PaginatedResponse,
  IsbnResult,
  ReadingStatus,
} from "./types";

const BASE = (import.meta.env.VITE_API_BASE as string) || "/api/v1";

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error: string }).error ?? res.statusText);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

// --- Profiles ---
export const profilesApi = {
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
  owner_profile_id?: number;
  page?: number;
  limit?: number;
}

export const booksApi = {
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
  create: (data: {
    title: string;
    isbn13?: string;
    cover_url?: string;
    description?: string;
    year?: number;
    language?: string;
    authors: string[];
    genres: string[];
    owner_profile_id?: number;
  }) => req<Book>("/books", { method: "POST", body: JSON.stringify(data) }),
  transferOwner: (id: number, profileId: number) =>
    req<Book>(`/books/${id}/owner`, {
      method: "PUT",
      body: JSON.stringify({ profile_id: profileId }),
    }),
  update: (
    id: number,
    data: {
      title?: string;
      isbn13?: string;
      cover_url?: string;
      description?: string;
      year?: number;
      language?: string;
      authors?: string[];
      genres?: string[];
    }
  ) =>
    req<Book>(`/books/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => req<void>(`/books/${id}`, { method: "DELETE" }),
};

// --- Authors & Genres ---
export const authorsApi = {
  list: () => req<{ id: number; name: string }[]>("/authors"),
};
export const genresApi = {
  list: () => req<{ id: number; name: string }[]>("/genres"),
};

// --- Library entries ---
export interface EntryQuery {
  status?: ReadingStatus;
  q?: string;
  page?: number;
  limit?: number;
}

export const libraryApi = {
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
  add: (
    profileId: number,
    bookId: number,
    status?: ReadingStatus,
    rating?: number,
    notes?: string
  ) =>
    req<LibraryEntry>(`/profiles/${profileId}/library`, {
      method: "POST",
      body: JSON.stringify({ book_id: bookId, status, rating, notes }),
    }),
  update: (
    profileId: number,
    bookId: number,
    data: { status?: ReadingStatus; rating?: number; notes?: string }
  ) =>
    req<LibraryEntry>(`/profiles/${profileId}/library/${bookId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  remove: (profileId: number, bookId: number) =>
    req<void>(`/profiles/${profileId}/library/${bookId}`, {
      method: "DELETE",
    }),
};

// --- ISBN ---
export const isbnApi = {
  lookup: (code: string) => req<IsbnResult>(`/isbn/${code}`),
};
