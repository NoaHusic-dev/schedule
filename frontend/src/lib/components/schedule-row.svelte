<script lang="ts">
  import { fade } from "svelte/transition";
  import { formatClassName, formatGroup } from "$lib/helpers/classes";
  import { currentPeriodId, formatPeriods, periodIds, slotsFor, type SlotView } from "$lib/helpers/timetable";
  import { clock } from "$lib/stores/clock.svelte";
  import { schedules } from "$lib/stores/schedules.svelte";

  type Props = {
    name: string;
    /** False when the server's schedule is for another day: then the row shows no lessons rather
     *  than another day's as current. */
    isToday: boolean;
    /** The row's place on the page, which staggers its entrance. */
    index: number;
  };
  let { name, isToday, index }: Props = $props();

  const classData = $derived(schedules.get(name));
  const row = $derived(
    isToday && classData ? slotsFor(classData.schedule, currentPeriodId(clock.minutes), clock.tick) : null,
  );
</script>

{#snippet lessonSlot(data: SlotView | null | undefined)}
  {#if data}
    <div class="text-center font-bold whitespace-nowrap tabular-nums">{formatPeriods(periodIds(data.slot.start, data.slot.end))}</div>
    <div class="room-badge">
      {#key data.key}
        <div class="board-ellipsis" in:fade={{ duration: 200 }}>{data.lesson.room}</div>
      {/key}
    </div>
    <div class="overflow-hidden px-[1vw]">
      {#key data.key}
        <div class="board-ellipsis" in:fade={{ duration: 200 }} out:fade={{ duration: 200 }}>
          {#if data.lesson.groups?.length}<span class="text-group">{formatGroup(data.lesson.groups)}</span>{/if}
          {data.lesson.subject.name}
        </div>
      {/key}
    </div>
    <div class="mr-[1em] overflow-hidden text-right font-normal text-teacher">
      {#key data.key}
        <div class="board-ellipsis" in:fade={{ duration: 200 }} out:fade={{ duration: 200 }}>{data.lesson.teacher.short}</div>
      {/key}
    </div>
  {:else}
    <div></div>
    <div></div>
    <div></div>
    <div></div>
  {/if}
{/snippet}

<!-- The class badge takes its colour from the end of `name` (RT, EL, ...) via layout.css, so it needs no per-row style. -->
<div class="board-row px-[1vw]" style:--i={index}>
  {#if !schedules.loaded}
    <div>...</div>
  {:else if classData}
    <div class="class-badge" data-razred={name}>{formatClassName(name)}</div>

    {@render lessonSlot(row?.current)}

    <div class="h-full w-0.5 bg-line"></div>

    {@render lessonSlot(row?.next)}

    <div class="class-badge justify-self-end" data-razred={name}>{formatClassName(name)}</div>
  {/if}
</div>