<script lang="ts">
  import { onMount } from "svelte";
  import { activeProfile } from "$lib/stores/profile";
  import { profilesApi } from "$lib/api";
  import type { Profile } from "$lib/types";

  let allProfiles: Profile[] = $state([]);
  let open = $state(false);

  onMount(async () => {
    try {
      allProfiles = await profilesApi.list();
    } catch {
      // backend may not be running yet
    }
  });

  function select(p: Profile) {
    activeProfile.set(p);
    open = false;
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as HTMLElement;
    if (!target.closest("[data-dropdown]")) {
      open = false;
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

<nav class="bg-white border-b border-gray-200 px-4 h-14 flex items-center gap-4 sticky top-0 z-40">
  <a href="#/" class="font-bold text-lg text-indigo-600">Orgalivro</a>

  <div class="flex gap-4 ml-4 text-sm">
    <a href="#/" class="text-gray-600 hover:text-indigo-600 font-medium">My Books</a>
    <a href="#/catalog" class="text-gray-600 hover:text-indigo-600 font-medium">Library</a>
    <a href="#/books/add" class="text-gray-600 hover:text-indigo-600 font-medium">+ Add Book</a>
  </div>

  <div class="ml-auto relative" data-dropdown>
    <button
      onclick={(e) => { e.stopPropagation(); open = !open; }}
      class="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-200 hover:bg-gray-50 text-sm"
    >
      {#if $activeProfile}
        <span class="w-6 h-6 rounded-full bg-indigo-100 text-indigo-700 font-bold text-xs flex items-center justify-center">
          {$activeProfile.name[0].toUpperCase()}
        </span>
        <span>{$activeProfile.name}</span>
      {:else}
        <span class="text-gray-400">Select profile</span>
      {/if}
      <span class="text-gray-400 text-xs">▾</span>
    </button>

    {#if open}
      <div class="absolute right-0 mt-1 w-52 bg-white border border-gray-200 rounded-lg shadow-lg z-50">
        {#each allProfiles as p (p.id)}
          <button
            onclick={() => select(p)}
            class="w-full text-left px-4 py-2 hover:bg-gray-50 text-sm flex items-center gap-2
              {$activeProfile?.id === p.id ? 'font-semibold text-indigo-600' : 'text-gray-700'}"
          >
            <span class="w-6 h-6 rounded-full bg-indigo-100 text-indigo-700 font-bold text-xs flex items-center justify-center">
              {p.name[0].toUpperCase()}
            </span>
            {p.name}
            {#if $activeProfile?.id === p.id}
              <span class="ml-auto text-xs">✓</span>
            {/if}
          </button>
        {/each}
        <hr class="my-1" />
        <a
          href="#/profiles"
          onclick={() => (open = false)}
          class="block px-4 py-2 text-sm text-gray-500 hover:bg-gray-50"
        >
          Manage profiles…
        </a>
      </div>
    {/if}
  </div>
</nav>
