<script lang="ts">
  import ScheduleFooter from "#lib/components/schedule-footer.svelte";
  import ScheduleHeader from "#lib/components/schedule-header.svelte";
  import ScheduleRow from "#lib/components/schedule-row.svelte";
  import TopBar from "#lib/components/top-bar.svelte";
  import { CLASS_PAGES, ROWS_PER_PAGE } from "#lib/helpers/classes.ts";
  import { isoDate } from "#lib/helpers/format.ts";
  import { periodStatus } from "#lib/helpers/timetable.ts";
  import { clock } from "#lib/stores/clock.svelte.ts";
  import { schedules } from "#lib/stores/schedules.svelte.ts";

  const ROTATION_SPEED = 8000;

  let pageIndex = $state(0);
  const classes = $derived(CLASS_PAGES[pageIndex]);
  const isToday = $derived(schedules.date === isoDate(clock.now));
  const isBigBreak = $derived(periodStatus(clock.now, isToday).bigBreak);

  $effect(() => clock.subscribe());
  $effect(() => schedules.subscribe());

  $effect(() => {
    const interval = setInterval(() => {
      pageIndex = (pageIndex + 1) % CLASS_PAGES.length;
    }, ROTATION_SPEED);
    return () => clearInterval(interval);
  });
</script>

<div class="relative isolate flex h-full flex-col">
  <div class="break-layer" aria-hidden="true">
    <div class="break-rays"></div>
    <svg class="break-wave break-wave-cyan" viewBox="0 0 2880 320" preserveAspectRatio="none">
      <path d="M0 160 C240 40 480 280 720 160 S1200 40 1440 160 S1920 280 2160 160 S2640 40 2880 160" />
    </svg>
    <svg class="break-wave break-wave-pink" viewBox="0 0 2880 320" preserveAspectRatio="none">
      <path d="M0 160 C240 280 480 40 720 160 S1200 280 1440 160 S1920 40 2160 160 S2640 280 2880 160" />
    </svg>
    <div class="break-flash" class:go={isBigBreak}></div>
  </div>
  <TopBar {isToday} />

  <main class="flex h-[90vh] flex-col">
    <ScheduleHeader />

    {#each classes as name, i (name)}
      <ScheduleRow {name} {isToday} index={i} />
    {/each}
    <!-- A short page keeps the striping and the footer at the bottom. -->
    {#each { length: ROWS_PER_PAGE - classes.length }, i (i)}
      <div class="board-row shrink-0" style:--i={classes.length + i} style:--dep="transparent"></div>
    {/each}

    <ScheduleFooter page={pageIndex} pages={CLASS_PAGES.length} rotation={ROTATION_SPEED} />
  </main>
</div>
