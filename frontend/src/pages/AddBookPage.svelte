<script lang="ts">
  import { push } from "svelte-spa-router";
  import { booksApi } from "$lib/api";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import type { IsbnResult } from "$lib/types";
  import IsbnLookup from "../components/isbn/IsbnLookup.svelte";
  import BookForm from "../components/book/BookForm.svelte";

  let prefill: IsbnResult | undefined = $state(undefined);

  async function handleSubmit(data: {
    title: string;
    isbn13: string;
    cover_url: string;
    description: string;
    year: number;
    language: string;
    authors: string[];
    genres: string[];
  }) {
    const book = await booksApi.create({
      ...data,
      owner_profile_id: $activeProfile?.id,
    });
    toast.success("Book added!");
    push(`/books/${book.id}`);
  }
</script>

<div class="max-w-xl mx-auto py-8 px-4 space-y-6">
  <h1 class="text-2xl font-bold text-gray-900">Add Book</h1>

  <div class="bg-white border border-gray-200 rounded-xl p-4 space-y-2">
    <p class="text-sm font-medium text-gray-700">Lookup by ISBN</p>
    <IsbnLookup onResult={(r) => (prefill = r)} />
    {#if prefill}
      <p class="text-xs text-green-600">✓ Pre-filled — review and save below.</p>
    {/if}
  </div>

  <div class="bg-white border border-gray-200 rounded-xl p-4">
    <p class="text-sm font-medium text-gray-700 mb-4">Book details</p>
    <BookForm {prefill} onsubmit={handleSubmit} />
  </div>
</div>
