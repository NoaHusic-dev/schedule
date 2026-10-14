<script lang="ts">
  import ScheduleFooter from "$lib/components/schedule-footer.svelte";
  import ScheduleHeader from "$lib/components/schedule-header.svelte";
  import ScheduleRow from "$lib/components/schedule-row.svelte";
  import TopBar from "$lib/components/top-bar.svelte";
  import { CLASS_PAGES, ROWS_PER_PAGE } from "$lib/helpers/classes";
  import { isoDate } from "$lib/helpers/format";
  import { clock } from "$lib/stores/clock.svelte";
  import { schedules } from "$lib/stores/schedules.svelte";

  const ROTATION_SPEED = 8000;

  let pageIndex = $state(0);
  const classes = $derived(CLASS_PAGES[pageIndex]);
  const isToday = $derived(schedules.date === isoDate(clock.now));

  $effect(() => clock.subscribe());
  $effect(() => schedules.subscribe());

  $effect(() => {
    const interval = setInterval(() => {
      pageIndex = (pageIndex + 1) % CLASS_PAGES.length;
    }, ROTATION_SPEED);
    return () => clearInterval(interval);
  });
</script>

<div class="flex h-full flex-col">
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
