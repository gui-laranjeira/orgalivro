<script lang="ts">
  import { onMount } from "svelte";
  import { booksApi, libraryApi } from "$lib/api";
  import { activeProfile } from "$lib/stores/profile";
  import { toast } from "$lib/stores/toast";
  import type { Book, LibraryEntry, ReadingStatus } from "$lib/types";
  import BookDetail from "../components/book/BookDetail.svelte";

  let { params }: { params: { id: string } } = $props();

  let book: Book | null = $state(null);
  let entry: LibraryEntry | null = $state(null);
  let error = $state("");

  onMount(async () => {
    try {
      book = await booksApi.get(Number(params.id));
      if ($activeProfile) {
        const res = await libraryApi.list($activeProfile.id);
        entry = res.data.find((e) => e.book.id === book!.id) ?? null;
      }
    } catch {
      error = "Book not found";
    }
  });

  async function addToLibrary(status: ReadingStatus) {
    if (!$activeProfile || !book) return;
    try {
      entry = await libraryApi.add($activeProfile.id, book.id, status);
      toast.success("Added to library!");
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Failed to add");
    }
  }

  async function updateEntry(data: {
    status?: ReadingStatus;
    rating?: number;
    notes?: string;
  }) {
    if (!$activeProfile || !book) return;
    try {
      entry = await libraryApi.update($activeProfile.id, book.id, data);
      toast.success("Updated!");
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Failed to update");
    }
  }

  async function transferOwner(profileId: number) {
    if (!book) return;
    try {
      book = await booksApi.transferOwner(book.id, profileId);
      toast.success("Ownership transferred!");
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Failed to transfer");
    }
  }
</script>

<div class="max-w-3xl mx-auto py-8 px-4">
  <a href="#/" class="text-sm text-indigo-600 hover:underline mb-4 inline-block"
    >← Back to library</a
  >

  {#if error}
    <p class="text-red-500 mt-4">{error}</p>
  {:else if !book}
    <p class="text-center text-gray-400 py-12">Loading…</p>
  {:else}
    <div class="bg-white border border-gray-200 rounded-xl p-6">
      <BookDetail
        {book}
        {entry}
        activeProfileId={$activeProfile?.id ?? null}
        onAddToLibrary={addToLibrary}
        onUpdateEntry={updateEntry}
        onTransferOwner={transferOwner}
      />
    </div>
  {/if}
</div>
