<script lang="ts">
  import { fade } from "svelte/transition";
  import { formatClassName, formatGroup } from "$lib/helpers/classes";
  import { currentPeriodId, getPeriodsFromRange, slotsFor, type SlotView } from "$lib/helpers/timetable";
  import { clock } from "$lib/stores/clock.svelte";
  import { schedules } from "$lib/stores/schedules.svelte";

  type Props = {
    name: string;
    /** False when the server's schedule is for another day: then the row shows no lessons rather
     *  than another day's as current. */
    isToday: boolean;
  };
  let { name, isToday }: Props = $props();

  const classData = $derived(schedules.get(name));
  const row = $derived(
    isToday && classData ? slotsFor(classData.schedule, currentPeriodId(clock.minutes), clock.tick) : null,
  );
</script>

{#snippet lessonSlot(data: SlotView | null | undefined)}
  {#if data}
    <div class="text-center font-mono font-bold">{getPeriodsFromRange(data.slot.start, data.slot.end)}</div>
    <div class="rounded bg-accent text-center font-bold text-black">
      {#key data.key}
        <div in:fade={{ duration: 200 }}>{data.lesson.room}</div>
      {/key}
    </div>
    <div class="overflow-hidden px-[1vw] text-ellipsis whitespace-nowrap">
      {#key data.key}
        <div in:fade={{ duration: 200 }} out:fade={{ duration: 200 }}>
          <span>{formatGroup(data.lesson.groups)}</span>
          {data.lesson.subject.name}
        </div>
      {/key}
    </div>
    <div class="mr-[1em] overflow-hidden text-right font-normal whitespace-nowrap opacity-50">
      {#key data.key}
        <div in:fade={{ duration: 200 }} out:fade={{ duration: 200 }}>{data.lesson.teacher.short}</div>
      {/key}
    </div>
  {:else}
    <div></div>
    <div></div>
    <div></div>
    <div></div>
  {/if}
{/snippet}

<div
  class="grid h-(--row-h) w-full grid-cols-(--board-cols) items-center px-[1vw] text-row odd:bg-board-alt even:bg-board"
>
  {#if !schedules.loaded}
    <div>...</div>
  {:else if classData}
    <div class="font-bold text-accent">{formatClassName(name)}</div>

    {@render lessonSlot(row?.current)}

    <div class="h-full w-0.5 bg-accent opacity-50"></div>

    {@render lessonSlot(row?.next)}

    <div class="font-bold text-right text-accent">{formatClassName(name)}</div>
  {/if}
</div>
