<script lang="ts">
  import { fade } from "svelte/transition";
  import { scheduleStore } from "$lib/stores/schedules.svelte";
  import { clock } from "$lib/stores/clock.svelte"; // Import clock
  import { getPeriodsFromRange } from "$util/timetable";
  import { formatClassName, formatGroup } from "$util/class";

  let { name }: { name: string } = $props();

  let classData = $derived(scheduleStore.classes[name]);
  let schedule = $derived(classData?.schedule ?? []);

  // REMOVED: local tick state and effect interval
  // We now use clock.tick directly

  let displayData = $derived.by(() => {
    if (!schedule.length) return null;

    const currentId = clock.currentPeriodId;
    if (!currentId) return null;

    // 1. Find Current Active Lessons
    const activeEntries = schedule.filter((item) => {
      const periods = getPeriodsFromRange(item.start, item.end);
      if (periods === "Svi sati") return currentId !== 8;
      return periods.split(",").map(Number).includes(currentId);
    });

    if (activeEntries.length === 0) return null;

    const currentPool = activeEntries.flatMap((entry) =>
      entry.lessons.map((l) => ({ entry, lesson: l })),
    );

    // USE GLOBAL TICK HERE
    const curIndex = clock.tick % currentPool.length;
    const currentItem = currentPool[curIndex];

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

    // 2. Find Next Lesson
    const nextPool = schedule
      .filter((item) => item.start >= currentItem.entry.end)
      .sort((a, b) => a.start.localeCompare(b.start));

    let nextItem = null;
    if (nextPool.length > 0) {
      const nextStartTime = nextPool[0].start;
      const candidates = nextPool.filter((i) => i.start === nextStartTime);

      const flattenedNext = candidates.flatMap((entry) =>
        entry.lessons.map((l) => ({ entry, lesson: l })),
      );

      if (flattenedNext.length > 0) {
        // USE GLOBAL TICK HERE AS WELL
        nextItem = flattenedNext[clock.tick % flattenedNext.length];
      }
    }

    return {
      current: {
        ...currentItem,
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
    <div class="col-room">
      {#key data.key}
        <div
          in:fade={{ duration: 200 }}
        >
          {data.lesson.room}
        </div>
      {/key}
    </div>

    <div class="col-subject">
      {#key data.key}
        <div
          in:fade={{ duration: 200 }}
          out:fade={{ duration: 200 }}
        >
          <span class="group">{formatGroup(data.lesson.groups)}</span>
          {data.lesson.subject.name}
        </div>
      {/key}
    </div>

    <div class="col-teacher">
      {#key data.key}
        <div
          in:fade={{ duration: 200 }}
          out:fade={{ duration: 200 }}
        >
          {data.lesson.teacher.short}
        </div>
      {/key}
    </div>
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
