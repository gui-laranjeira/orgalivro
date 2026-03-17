<script lang="ts">
  import type { LibraryEntry, ReadingStatus } from "$lib/types";
  import BookStatusBadge from "../book/BookStatusBadge.svelte";
  import { starRating } from "$lib/utils";

  let {
    entry,
    onUpdate,
    onRemove,
  }: {
    entry: LibraryEntry;
    onUpdate: (data: { status?: ReadingStatus }) => Promise<void>;
    onRemove: () => Promise<void>;
  } = $props();
</script>

<div class="bg-white border border-gray-200 rounded-xl p-3 flex gap-3">
  <a href="#/books/{entry.book.id}" class="flex-shrink-0">
    {#if entry.book.cover_url}
      <img
        src={entry.book.cover_url}
        alt={entry.book.title}
        class="w-14 h-20 object-cover rounded"
      />
    {:else}
      <div class="w-14 h-20 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs">
        ?
      </div>
    {/if}
  </a>

  <div class="flex-1 min-w-0">
    <a
      href="#/books/{entry.book.id}"
      class="font-semibold text-sm line-clamp-2 hover:text-indigo-600 text-gray-900"
    >
      {entry.book.title}
    </a>
    {#if entry.book.authors?.length}
      <p class="text-xs text-gray-500 mt-0.5">
        {entry.book.authors.map((a) => a.name).join(", ")}
      </p>
    {/if}

    <div class="flex items-center gap-2 mt-2 flex-wrap">
      <BookStatusBadge status={entry.status} />
      {#if entry.rating}
        <span class="text-xs text-yellow-500">{starRating(entry.rating)}</span>
      {/if}
    </div>

    <div class="flex gap-2 mt-2 items-center">
      <select
        value={entry.status}
        onchange={(e) =>
          onUpdate({
            status: (e.target as HTMLSelectElement).value as ReadingStatus,
          })}
        class="border border-gray-200 rounded px-2 py-0.5 text-xs focus:outline-none"
      >
        <option value="want_to_read">Want to Read</option>
        <option value="reading">Reading</option>
        <option value="read">Read</option>
      </select>
      <button
        onclick={onRemove}
        class="text-xs text-red-400 hover:text-red-600 ml-1"
      >
        Remove
      </button>
    </div>
  </div>
</div>
