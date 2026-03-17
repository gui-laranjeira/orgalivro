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

export interface BookReader {
  profile_id: number;
  profile_name: string;
  status: ReadingStatus;
}

export interface Book {
  id: number;
  title: string;
  isbn13: string;
  cover_url: string;
  description: string;
  year: number;
  language: string;
  owner_profile_id: number | null;
  owner_profile_name: string;
  authors: Author[];
  genres: Genre[];
  readers: BookReader[];
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
