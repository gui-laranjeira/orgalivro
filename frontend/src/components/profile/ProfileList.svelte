<script lang="ts">
  import type { Profile } from "$lib/types";
  import { activeProfile } from "$lib/stores/profile";

  let {
    profiles,
    onDelete,
    onSelect,
  }: {
    profiles: Profile[];
    onDelete: (id: number) => Promise<void>;
    onSelect: (p: Profile) => void;
  } = $props();
</script>

{#if profiles.length === 0}
  <p class="text-gray-400 text-sm text-center py-4">No profiles yet.</p>
{:else}
  <ul class="space-y-2">
    {#each profiles as p (p.id)}
      <li
        class="flex items-center gap-3 bg-white border border-gray-200 rounded-xl px-4 py-3"
      >
        <div
          class="w-9 h-9 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-700 font-bold text-sm flex-shrink-0"
        >
          {p.name[0].toUpperCase()}
        </div>
        <span class="flex-1 font-medium text-sm text-gray-800">{p.name}</span>
        {#if $activeProfile?.id === p.id}
          <span class="text-xs text-indigo-600 font-medium">Active</span>
        {:else}
          <button
            onclick={() => onSelect(p)}
            class="text-xs text-gray-500 hover:text-indigo-600"
          >
            Set active
          </button>
        {/if}
        <button
          onclick={() => onDelete(p.id)}
          class="text-xs text-red-400 hover:text-red-600 ml-1"
        >
          Delete
        </button>
      </li>
    {/each}
  </ul>
{/if}
