<script lang="ts">
  import { toast } from "$lib/stores/toast";

  const TYPE_CLASSES: Record<string, string> = {
    success: "bg-green-600",
    error: "bg-red-600",
    info: "bg-gray-800",
  };
</script>

<div class="fixed bottom-4 right-4 flex flex-col gap-2 z-50">
  {#each $toast as t (t.id)}
    <div class="px-4 py-2 rounded-lg text-white text-sm shadow-lg {TYPE_CLASSES[t.type]}">
      {t.message}
    </div>
  {/each}
</div>
