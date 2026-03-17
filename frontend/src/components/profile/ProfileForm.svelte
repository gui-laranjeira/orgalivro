<script lang="ts">
  let { oncreate }: { oncreate: (name: string) => Promise<void> } = $props();

  let name = $state("");
  let error = $state("");
  let loading = $state(false);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim()) return;
    loading = true;
    error = "";
    try {
      await oncreate(name.trim());
      name = "";
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Failed to create profile";
    } finally {
      loading = false;
    }
  }
</script>

<form onsubmit={handleSubmit} class="space-y-2">
  <div class="flex gap-2 items-start">
    <div class="flex-1">
      <input
        bind:value={name}
        placeholder="Profile name"
        required
        class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
      />
      {#if error}
        <p class="text-red-500 text-xs mt-1">{error}</p>
      {/if}
    </div>
    <button
      type="submit"
      disabled={loading}
      class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 disabled:opacity-50 transition"
    >
      {loading ? "Creating…" : "Create"}
    </button>
  </div>
</form>
