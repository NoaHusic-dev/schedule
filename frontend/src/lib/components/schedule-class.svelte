<script lang="ts">
  import { fade } from "svelte/transition";
  import { scheduleStore } from "$lib/stores/schedules.svelte";
  import { clock } from "$lib/stores/clock.svelte";
  import { getPeriodsFromRange } from "$util/timetable";
  import { formatClassName, formatGroup } from "$util/class";

  let { name }: { name: string } = $props();

  let classData = $derived(scheduleStore.classes[name]);
  let schedule = $derived(classData?.schedule ?? []);

  // 5-second rotation ticker
  let tick = $state(0);
  $effect(() => {
    const interval = setInterval(() => {
      tick++;
    }, 5000);
    return () => clearInterval(interval);
  });

  /**
   * Calculates exactly what to show in Current and Next slots.
   * Updates reactively whenever the clock ticks or the 5s interval fires.
   */
  let displayData = $derived.by(() => {
    if (!schedule.length) return null;

    const currentId = clock.currentPeriodId;
    // If no school period is active, assume start of day or show nothing
    if (!currentId) return null;

    // 1. Find Current Active Lessons
    const activeEntries = schedule.filter((item) => {
      const periods = getPeriodsFromRange(item.start, item.end);

      // FIX: If "All Hours", show it UNLESS it is the 8th period
      if (periods === "Svi sati") return currentId !== 8;

      return periods.split(",").map(Number).includes(currentId);
    });

    if (activeEntries.length === 0) return null;

    // Flatten to handle "1 SE" split entries case
    const currentPool = activeEntries.flatMap((entry) =>
      entry.lessons.map((l) => ({ entry, lesson: l })),
    );

    // Pick the current lesson based on rotation
    const curIndex = tick % currentPool.length;
    const currentItem = currentPool[curIndex];

    // Check for "Svi sati" - if true, no next lesson needed
    const curPeriods = getPeriodsFromRange(
      currentItem.entry.start,
      currentItem.entry.end,
    );
    if (curPeriods === "Svi sati") {
      return {
        current: { ...currentItem, key: `${currentItem.lesson.subject.short}` },
        next: null,
      };
    }

    // 2. Find Next Lesson relative to the CHOSEN Current Lesson
    // We look for classes starting AFTER the current one ends.
    const nextPool = schedule
      .filter((item) => item.start >= currentItem.entry.end) // Starts after current ends
      .sort((a, b) => a.start.localeCompare(b.start)); // Earliest first

    let nextItem = null;
    if (nextPool.length > 0) {
      // Find the group of classes starting at the soonest time
      const nextStartTime = nextPool[0].start;
      const candidates = nextPool.filter((i) => i.start === nextStartTime);

      const flattenedNext = candidates.flatMap((entry) =>
        entry.lessons.map((l) => ({ entry, lesson: l })),
      );

      // Rotate next lesson using the same tick (synced)
      if (flattenedNext.length > 0) {
        nextItem = flattenedNext[tick % flattenedNext.length];
      }
    }

    return {
      current: {
        ...currentItem,
        // Key is CONTENT-based. If subject/group doesn't change, NO transition.
        key: `${currentItem.lesson.subject.short}-${currentItem.lesson.groups?.join("")}`,
      },
      next: nextItem
        ? {
            ...nextItem,
            key: `${nextItem.lesson.subject.short}-${nextItem.lesson.groups?.join("")}`,
          }
        : null,
    };
  });
</script>

{#snippet lessonSlot(data: any)}
  {#if data}
    <div class="col-time">
      {getPeriodsFromRange(data.entry.start, data.entry.end)}
    </div>
    <div class="col-room">{data.lesson.room}</div>

    <div class="col-subject">
      {#key data.key}
        <div
          in:fade={{ duration: 200 }}
          out:fade={{ duration: 200 }}
          class="lesson-inner"
        >
          <span class="group">{formatGroup(data.lesson.groups)}</span>
          {data.lesson.subject.name}
        </div>
      {/key}
    </div>

    <div class="col-teacher">{data.lesson.teacher.short}</div>
  {:else}
    <div class="col-time"></div>
    <div class="col-room"></div>
    <div class="col-subject"></div>
    <div class="col-teacher"></div>
  {/if}
{/snippet}

<div class="grid-row">
  {#if scheduleStore.isLoading && !classData}
    <div class="loading">...</div>
  {:else if classData}
    <div class="col-class">{formatClassName(name)}</div>

    {@render lessonSlot(displayData?.current)}

    <div class="col-divider line"></div>

    {@render lessonSlot(displayData?.next)}

    <div class="col-class right-align">{formatClassName(name)}</div>
  {/if}
</div>
