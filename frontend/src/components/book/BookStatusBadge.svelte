<script lang="ts">
  import { STATUS_LABELS, STATUS_COLORS } from "$lib/utils";
  let { status }: { status: string } = $props();
</script>

<span class="px-2 py-0.5 rounded-full text-xs font-medium {STATUS_COLORS[status] ?? 'bg-gray-100 text-gray-700'}">
  {STATUS_LABELS[status] ?? status}
</span>
