<script lang="ts">
  import { formatClock, formatDate } from "$lib/helpers/format";
  import { clock } from "$lib/stores/clock.svelte";
</script>

<div class="text-right">
  <div class="text-clock font-extrabold tabular-nums">{formatClock(clock.now)}</div>
  <div class="-mt-[0.2vh] text-date font-semibold text-muted">{formatDate(clock.now)}</div>
</div>
