<script lang="ts">
  import type { Book } from "$lib/types";
  import { STATUS_COLORS } from "$lib/utils";
  let { book }: { book: Book } = $props();
</script>

<a
  href="#/books/{book.id}"
  class="block bg-white rounded-xl border border-gray-200 hover:shadow-md transition p-3 flex gap-3"
>
  {#if book.cover_url}
    <img
      src={book.cover_url}
      alt={book.title}
      class="w-16 h-24 object-cover rounded flex-shrink-0"
    />
  {:else}
    <div
      class="w-16 h-24 bg-gray-100 rounded flex items-center justify-center text-gray-400 text-xs flex-shrink-0"
    >
      No cover
    </div>
  {/if}
  <div class="flex-1 min-w-0">
    <p class="font-semibold text-sm leading-tight line-clamp-2 text-gray-900">
      {book.title}
    </p>
    {#if book.authors?.length}
      <p class="text-xs text-gray-500 mt-1">
        {book.authors.map((a) => a.name).join(", ")}
      </p>
    {/if}
    {#if book.year}
      <p class="text-xs text-gray-400 mt-1">{book.year}</p>
    {/if}
    {#if book.genres?.length}
      <div class="flex gap-1 mt-2 flex-wrap">
        {#each book.genres.slice(0, 2) as g}
          <span class="px-1.5 py-0.5 bg-gray-100 text-gray-500 rounded text-xs">
            {g.name}
          </span>
        {/each}
      </div>
    {/if}
    {#if book.readers?.length}
      <div class="flex gap-1 mt-2 flex-wrap">
        {#each book.readers as reader}
          <span class="px-1.5 py-0.5 rounded text-xs font-medium {STATUS_COLORS[reader.status]}">
            {reader.profile_name}
          </span>
        {/each}
      </div>
    {/if}
  </div>
</a>
