<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import { booksApi } from "$lib/api";
  import { activeProfile } from "$lib/stores/profile";
  import { debounce } from "$lib/utils";
  import type { Book } from "$lib/types";
  import BookCard from "../components/book/BookCard.svelte";

  let books: Book[] = $state([]);
  let total = $state(0);
  let q = $state("");
  let loading = $state(false);

  const load = debounce(async () => {
    if (!$activeProfile) return;
    loading = true;
    try {
      const res = await booksApi.list({
        owner_profile_id: $activeProfile.id,
        q: q || undefined,
      });
      books = res.data;
      total = res.total;
    } catch {
      // ignore
    } finally {
      loading = false;
    }
  }, 300);

  onMount(() => {
    if (!$activeProfile) {
      push("/profiles");
      return;
    }
    load();
  });

  $effect(() => {
    void $activeProfile;
    void q;
    load();
  });
</script>

<div class="max-w-3xl mx-auto py-8 px-4 space-y-4">
  <div class="flex items-center gap-2">
    <h1 class="text-2xl font-bold text-gray-900 flex-1">
      {$activeProfile ? `${$activeProfile.name}'s Books` : "My Books"}
    </h1>
    <span class="text-sm text-gray-500">{total} book{total !== 1 ? "s" : ""}</span>
    <a
      href="#/books/add"
      class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 transition"
    >+ Add</a>
  </div>

  <input
    bind:value={q}
    placeholder="Search your books…"
    class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
  />

  {#if loading}
    <p class="text-center text-gray-400 py-12">Loading…</p>
  {:else if books.length === 0}
    <div class="text-center py-12 space-y-2">
      <p class="text-gray-400">No books here yet.</p>
      <a href="#/books/add" class="text-indigo-600 text-sm hover:underline">Add your first book!</a>
    </div>
  {:else}
    <div class="space-y-3">
      {#each books as book (book.id)}
        <BookCard {book} />
      {/each}
    </div>
  {/if}
</div>
