<script lang="ts">
  import logo from "$lib/assets/tsd-logo-dark.png";
  import Clock from "$lib/components/clock.svelte";
  import PeriodStatus from "$lib/components/period-status.svelte";
  import Weather from "$lib/components/weather.svelte";

  let { isToday }: { isToday: boolean } = $props();
</script>

<header class="grid h-[10vh] grid-cols-[1fr_auto_1fr] items-center bg-board px-[2vw]">
  <div class="flex items-center gap-[1vw]">
    <!-- Decorative: the school's name is right beside it. -->
    <img src={logo} alt="" class="h-[7vh] w-auto" />
    <div>
      <h1 class="mb-[0.2vh] text-brand font-extrabold tracking-[1px]">TEHNIČKA ŠKOLA DARUVAR</h1>
      <span class="text-subtitle font-semibold tracking-[0.2em] text-muted uppercase">Pregled nastave uživo</span>
    </div>
  </div>

  <PeriodStatus {isToday} />

  <div class="flex items-center gap-[1.25vw] justify-self-end">
    <Weather />
    <Clock />
  </div>
</header>
