<script lang="ts">
  import { clock } from "$lib/stores/clock.svelte";
</script>

<div class="clock-display">
  {clock.formatted}
</div>
