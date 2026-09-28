<script lang="ts">
  import { formatClock } from "$lib/helpers/format";
  import { clock } from "$lib/stores/clock.svelte";
</script>

<div class="text-clock font-extrabold">{formatClock(clock.now)}</div>
