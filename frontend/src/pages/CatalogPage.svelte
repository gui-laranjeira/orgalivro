<script lang="ts">
  import { onMount } from "svelte";
  import { booksApi, genresApi } from "$lib/api";
  import type { Book } from "$lib/types";
  import { debounce } from "$lib/utils";
  import BookCard from "../components/book/BookCard.svelte";

  let bookList: Book[] = $state([]);
  let total = $state(0);
  let q = $state("");
  let genre = $state("");
  let genres: { id: number; name: string }[] = $state([]);
  let loading = $state(false);

  const loadBooks = debounce(async () => {
    loading = true;
    try {
      const res = await booksApi.list({ q: q || undefined, genre: genre || undefined });
      bookList = res.data;
      total = res.total;
    } catch {
      // ignore
    } finally {
      loading = false;
    }
  }, 300);

  onMount(async () => {
    genres = await genresApi.list();
    loadBooks();
  });

  $effect(() => {
    void q;
    void genre;
    loadBooks();
  });
</script>

<div class="max-w-3xl mx-auto py-8 px-4 space-y-4">
  <div class="flex items-center gap-3">
    <h1 class="text-2xl font-bold text-gray-900 flex-1">Catalog</h1>
    <span class="text-sm text-gray-500">{total} book{total !== 1 ? "s" : ""}</span>
    <a
      href="#/books/add"
      class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 transition"
    >
      + Add
    </a>
  </div>

  <div class="flex gap-2">
    <input
      bind:value={q}
      placeholder="Search by title or description…"
      class="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    />
    <select
      bind:value={genre}
      class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-white"
    >
      <option value="">All categories</option>
      {#each genres as g}
        <option value={g.name}>{g.name}</option>
      {/each}
    </select>
  </div>

  {#if loading}
    <p class="text-center text-gray-400 py-12">Loading…</p>
  {:else if bookList.length === 0}
    <div class="text-center py-12 space-y-2">
      <p class="text-gray-400">No books found.</p>
      <a href="#/books/add" class="text-indigo-600 text-sm hover:underline"
        >Add your first book</a
      >
    </div>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      {#each bookList as book (book.id)}
        <BookCard {book} />
      {/each}
    </div>
  {/if}
</div>
