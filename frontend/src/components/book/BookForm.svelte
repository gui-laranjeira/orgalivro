<script lang="ts">
  import type { IsbnResult } from "$lib/types";

  interface BookFormData {
    title: string;
    isbn13: string;
    cover_url: string;
    description: string;
    year: number;
    language: string;
    authors: string[];
    genres: string[];
  }

  let {
    prefill = undefined,
    onsubmit,
  }: {
    prefill?: IsbnResult;
    onsubmit: (data: BookFormData) => Promise<void>;
  } = $props();

  let title = $state("");
  let isbn13 = $state("");
  let cover_url = $state("");
  let description = $state("");
  let year = $state<number | "">("");
  let language = $state("");
  let authorsInput = $state("");
  let genresInput = $state("");
  let submitting = $state(false);
  let error = $state("");

  $effect(() => {
    if (prefill) {
      title = prefill.title ?? "";
      isbn13 = prefill.isbn13 ?? "";
      cover_url = prefill.cover_url ?? "";
      description = prefill.description ?? "";
      year = prefill.year ?? "";
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
        year: Number(year) || 0,
        language,
        authors: authorsInput
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
        genres: genresInput
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
      });
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      submitting = false;
    }
  }
</script>

<form onsubmit={handleSubmit} class="space-y-4">
  {#if error}
    <p class="text-red-600 text-sm bg-red-50 px-3 py-2 rounded-lg">{error}</p>
  {/if}

  <div>
    <label for="bf-title" class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
    <input
      id="bf-title"
      bind:value={title}
      required
      class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    />
  </div>

  <div>
    <label for="bf-authors" class="block text-sm font-medium text-gray-700 mb-1"
      >Authors (comma-separated)</label
    >
    <input
      id="bf-authors"
      bind:value={authorsInput}
      placeholder="Frank Herbert, Ursula K. Le Guin"
      class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    />
  </div>

  <div class="grid grid-cols-2 gap-3">
    <div>
      <label for="bf-isbn" class="block text-sm font-medium text-gray-700 mb-1">ISBN-13</label>
      <input
        id="bf-isbn"
        bind:value={isbn13}
        placeholder="9780261102354"
        class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
      />
    </div>
    <div>
      <label for="bf-year" class="block text-sm font-medium text-gray-700 mb-1">Year</label>
      <input
        id="bf-year"
        type="number"
        bind:value={year}
        placeholder="1965"
        class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
      />
    </div>
  </div>

  <div>
    <label for="bf-cover" class="block text-sm font-medium text-gray-700 mb-1">Cover URL</label>
    <input
      id="bf-cover"
      bind:value={cover_url}
      placeholder="https://…"
      class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    />
  </div>

  <div>
    <label for="bf-lang" class="block text-sm font-medium text-gray-700 mb-1">Language</label>
    <input
      id="bf-lang"
      bind:value={language}
      placeholder="en"
      class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    />
  </div>

  <div>
    <label for="bf-genres" class="block text-sm font-medium text-gray-700 mb-1"
      >Genres (comma-separated)</label
    >
    <input
      id="bf-genres"
      bind:value={genresInput}
      placeholder="Fantasy, Science Fiction"
      class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    />
  </div>

  <div>
    <label for="bf-desc" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
    <textarea
      id="bf-desc"
      bind:value={description}
      rows={4}
      class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
    ></textarea>
  </div>

  <button
    type="submit"
    disabled={submitting}
    class="w-full bg-indigo-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition"
  >
    {submitting ? "Saving…" : "Save Book"}
  </button>
</form>
