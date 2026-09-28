<script lang="ts">
  import { periodStatus } from "$lib/helpers/timetable";
  import { clock } from "$lib/stores/clock.svelte";

  let { isToday }: { isToday: boolean } = $props();

  const status = $derived(periodStatus(clock.now, isToday));
</script>

<div class="min-w-[29vw] text-center">
  <div class="flex items-baseline justify-center gap-[0.75vw] whitespace-nowrap">
    <span class="text-status font-extrabold tracking-[0.02em] text-accent uppercase">{status.title}</span>
    {#if status.detail}
      <span class="text-status-detail font-semibold text-ink/80 tabular-nums">{status.detail}</span>
    {/if}
  </div>
  {#if status.progress !== null}
    <div class="mx-auto mt-[0.8vh] h-[0.55vh] w-[23vw] overflow-hidden rounded-full bg-line">
      <div class="h-full bg-accent transition-[width] duration-1000 ease-linear" style:width="{status.progress * 100}%"></div>
    </div>
  {/if}
</div>
