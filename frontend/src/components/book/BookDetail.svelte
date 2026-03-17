<script lang="ts">
  import type { Book, LibraryEntry, Profile, ReadingStatus } from "$lib/types";
  import BookStatusBadge from "./BookStatusBadge.svelte";
  import { starRating, STATUS_COLORS } from "$lib/utils";
  import { profilesApi } from "$lib/api";

  let {
    book,
    entry = null,
    activeProfileId = null,
    onAddToLibrary,
    onUpdateEntry,
    onTransferOwner,
  }: {
    book: Book;
    entry?: LibraryEntry | null;
    activeProfileId?: number | null;
    onAddToLibrary: (status: ReadingStatus) => Promise<void>;
    onUpdateEntry: (data: {
      status?: ReadingStatus;
      rating?: number;
      notes?: string;
    }) => Promise<void>;
    onTransferOwner?: (profileId: number) => Promise<void>;
  } = $props();

  let showTransfer = $state(false);
  let allProfiles: Profile[] = $state([]);
  let selectedTransferProfile = $state(0);

  async function openTransfer() {
    allProfiles = await profilesApi.list();
    selectedTransferProfile = 0;
    showTransfer = true;
  }

  async function confirmTransfer() {
    if (!selectedTransferProfile || !onTransferOwner) return;
    await onTransferOwner(selectedTransferProfile);
    showTransfer = false;
  }

  const isOwner = $derived(
    book.owner_profile_id != null && activeProfileId === book.owner_profile_id
  );

  let editNotes = $state(false);
  let notesValue = $state("");

  $effect(() => {
    notesValue = entry?.notes ?? "";
  });

  async function saveNotes() {
    await onUpdateEntry({ notes: notesValue });
    editNotes = false;
  }
</script>

<div class="flex gap-6 flex-wrap">
  {#if book.cover_url}
    <img
      src={book.cover_url}
      alt={book.title}
      class="w-36 h-52 object-cover rounded-xl shadow flex-shrink-0"
    />
  {:else}
    <div
      class="w-36 h-52 bg-gray-100 rounded-xl flex items-center justify-center text-gray-400 text-sm flex-shrink-0"
    >
      No cover
    </div>
  {/if}

  <div class="flex-1 min-w-0">
    <h1 class="text-2xl font-bold text-gray-900">{book.title}</h1>

    {#if book.authors?.length}
      <p class="text-gray-600 mt-1 text-sm">
        {book.authors.map((a) => a.name).join(", ")}
      </p>
    {/if}

    <div class="flex gap-2 mt-2 flex-wrap text-sm text-gray-500">
      {#if book.year}<span>{book.year}</span>{/if}
      {#if book.language}<span>· {book.language.toUpperCase()}</span>{/if}
      {#if book.isbn13}<span>· ISBN {book.isbn13}</span>{/if}
    </div>

    {#if book.genres?.length}
      <div class="flex gap-1 mt-3 flex-wrap">
        {#each book.genres as g}
          <span class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded text-xs"
            >{g.name}</span
          >
        {/each}
      </div>
    {/if}

    <div class="flex flex-wrap gap-x-4 gap-y-1 mt-3 text-xs text-gray-500 items-center">
      {#if book.owner_profile_name}
        <span>Owner: <span class="font-medium text-gray-700">{book.owner_profile_name}</span></span>
      {/if}
      {#if isOwner}
        <button
          onclick={openTransfer}
          class="text-indigo-600 hover:underline"
        >Transfer ownership</button>
      {/if}
    </div>

    {#if showTransfer}
      <div class="mt-3 flex items-center gap-2 flex-wrap">
        <select
          bind:value={selectedTransferProfile}
          class="border border-gray-300 rounded px-2 py-1 text-sm focus:outline-none"
        >
          <option value={0}>Select profile…</option>
          {#each allProfiles.filter(p => p.id !== book.owner_profile_id) as p}
            <option value={p.id}>{p.name}</option>
          {/each}
        </select>
        <button
          onclick={confirmTransfer}
          disabled={!selectedTransferProfile}
          class="px-3 py-1 bg-indigo-600 text-white text-xs rounded hover:bg-indigo-700 disabled:opacity-40"
        >Confirm</button>
        <button
          onclick={() => (showTransfer = false)}
          class="px-3 py-1 border rounded text-xs hover:bg-gray-50"
        >Cancel</button>
      </div>
    {/if}

    {#if book.readers?.length}
      <div class="flex gap-1.5 mt-3 flex-wrap items-center">
        <span class="text-xs text-gray-400">Reading:</span>
        {#each book.readers as reader}
          <span class="px-2 py-0.5 rounded-full text-xs font-medium {STATUS_COLORS[reader.status]}">
            {reader.profile_name}
          </span>
        {/each}
      </div>
    {/if}

    <div class="mt-5 space-y-3">
      {#if entry}
        <div class="flex items-center gap-3 flex-wrap">
          <BookStatusBadge status={entry.status} />
          {#if entry.rating}
            <span class="text-sm text-yellow-500">{starRating(entry.rating)}</span>
          {/if}
        </div>

        <div class="flex items-center gap-2 flex-wrap">
          <label for="bd-status" class="text-xs text-gray-500">Status:</label>
          <select
            id="bd-status"
            value={entry.status}
            onchange={(e) =>
              onUpdateEntry({
                status: (e.target as HTMLSelectElement).value as ReadingStatus,
              })}
            class="border border-gray-300 rounded px-2 py-1 text-sm focus:outline-none"
          >
            <option value="want_to_read">Want to Read</option>
            <option value="reading">Reading</option>
            <option value="read">Read</option>
          </select>

          <label for="bd-rating" class="text-xs text-gray-500 ml-2">Rating:</label>
          <select
            id="bd-rating"
            value={entry.rating ?? ""}
            onchange={(e) => {
              const v = (e.target as HTMLSelectElement).value;
              onUpdateEntry({ rating: v ? Number(v) : undefined });
            }}
            class="border border-gray-300 rounded px-2 py-1 text-sm focus:outline-none"
          >
            <option value="">—</option>
            {#each [1, 2, 3, 4, 5] as n}
              <option value={n}>{"★".repeat(n)}</option>
            {/each}
          </select>
        </div>

        {#if editNotes}
          <div class="space-y-1">
            <textarea
              bind:value={notesValue}
              rows={3}
              class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
            ></textarea>
            <div class="flex gap-2">
              <button
                onclick={saveNotes}
                class="px-3 py-1 bg-indigo-600 text-white rounded text-xs hover:bg-indigo-700"
                >Save</button
              >
              <button
                onclick={() => (editNotes = false)}
                class="px-3 py-1 border rounded text-xs hover:bg-gray-50">Cancel</button
              >
            </div>
          </div>
        {:else}
          <div class="text-sm">
            {#if entry.notes}
              <p class="text-gray-700 italic">"{entry.notes}"</p>
            {/if}
            <button
              onclick={() => {
                notesValue = entry?.notes ?? "";
                editNotes = true;
              }}
              class="text-xs text-indigo-600 hover:underline mt-1"
            >
              {entry.notes ? "Edit notes" : "Add notes"}
            </button>
          </div>
        {/if}
      {:else}
        <p class="text-sm text-gray-500">Not in your library yet.</p>
        <div class="flex gap-2 flex-wrap">
          <button
            onclick={() => onAddToLibrary("want_to_read")}
            class="text-sm px-3 py-1.5 border border-gray-300 rounded-lg hover:bg-gray-50"
            >Want to Read</button
          >
          <button
            onclick={() => onAddToLibrary("reading")}
            class="text-sm px-3 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >Reading</button
          >
          <button
            onclick={() => onAddToLibrary("read")}
            class="text-sm px-3 py-1.5 bg-green-600 text-white rounded-lg hover:bg-green-700"
            >Read</button
          >
        </div>
      {/if}
    </div>
  </div>
</div>

{#if book.description}
  <div class="mt-8">
    <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2">
      Description
    </h2>
    <p class="text-sm text-gray-700 leading-relaxed">{book.description}</p>
  </div>
{/if}
