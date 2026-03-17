<script lang="ts">
  import { debounce } from "$lib/utils";

  let {
    onFilter,
  }: {
    onFilter: (filters: { status?: string; q?: string }) => void;
  } = $props();

  let status = $state("");
  let q = $state("");

  const emitDebounced = debounce(() => {
    onFilter({ status: status || undefined, q: q || undefined });
  }, 300);

  $effect(() => {
    // reactive: re-runs whenever status or q changes
    void status;
    void q;
    emitDebounced();
  });
</script>

<div class="flex gap-3 items-center">
  <input
    bind:value={q}
    placeholder="Search books…"
    class="border border-gray-300 rounded-lg px-3 py-1.5 text-sm flex-1 focus:outline-none focus:ring-2 focus:ring-indigo-500"
  />
  <select
    bind:value={status}
    class="border border-gray-300 rounded-lg px-2 py-1.5 text-sm focus:outline-none"
  >
    <option value="">All statuses</option>
    <option value="want_to_read">Want to Read</option>
    <option value="reading">Reading</option>
    <option value="read">Read</option>
  </select>
</div>
