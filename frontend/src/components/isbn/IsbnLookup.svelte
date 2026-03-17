<script lang="ts">
  import { isbnApi } from "$lib/api";
  import type { IsbnResult } from "$lib/types";

  let { onResult }: { onResult: (result: IsbnResult) => void } = $props();

  let value = $state("");
  let loading = $state(false);
  let error = $state("");

  async function lookup() {
    const code = value.trim();
    if (!code) return;
    loading = true;
    error = "";
    try {
      const result = await isbnApi.lookup(code);
      onResult(result);
      value = "";
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Lookup failed";
    } finally {
      loading = false;
    }
  }
</script>

<div class="space-y-1">
  <div class="flex gap-2">
    <input
      bind:value
      placeholder="ISBN-13 (e.g. 9780261102354)"
      class="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
      onkeydown={(e) => e.key === "Enter" && lookup()}
    />
    <button
      onclick={lookup}
      disabled={loading}
      class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 disabled:opacity-50 transition"
    >
      {loading ? "…" : "Lookup"}
    </button>
  </div>
  {#if error}
    <p class="text-red-500 text-xs">{error}</p>
  {/if}
</div>
