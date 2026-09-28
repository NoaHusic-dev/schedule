<script lang="ts">
  import Clock from "$lib/components/clock.svelte";
  import PeriodStatus from "$lib/components/period-status.svelte";

  let { isToday }: { isToday: boolean } = $props();
</script>

<header class="grid h-[10vh] grid-cols-[1fr_auto_1fr] items-center bg-board px-[2vw]">
  <div>
    <h1 class="mb-[0.2vh] text-brand font-extrabold tracking-[1px]">TEHNIČKA ŠKOLA DARUVAR</h1>
    <span class="text-subtitle font-semibold tracking-[0.2em] text-muted uppercase">Pregled nastave uživo</span>
  </div>

  <PeriodStatus {isToday} />

  <div class="justify-self-end">
    <Clock />
  </div>
</header>
