<script lang="ts">
  import { fade } from 'svelte/transition';
  import { scheduleStore } from "$lib/stores/schedules.svelte";
  import { clock } from "$lib/stores/clock.svelte";
  import { getPeriodsFromRange } from "$util/timetable";
  import { formatClassName, formatGroup } from "$util/class";

  let { name }: { name: string } = $props();

  let classData = $derived(scheduleStore.classes[name]);
  let schedule = $derived(classData?.schedule ?? []);
  
  let tick = $state(0);
  $effect(() => {
    const interval = setInterval(() => { tick++; }, 5000);
    return () => clearInterval(interval);
  });

  // 1. Helper: Formats a raw schedule entry into the shape your snippet needs
  function formatLessonData(entry: any) {
    if (!entry) return null;
    
    const lessonIdx = tick % entry.lessons.length;
    const lesson = entry.lessons[lessonIdx];

    return {
      entry,
      lesson,
      // Use start time for key instead of ID to ensure uniqueness between blocks
      transitionKey: `${entry.start}-${lesson.subject.short}-${lesson.groups?.join('')}`
    };
  }

  // 2. Logic for CURRENT (Stays Time/Period based)
  let current = $derived.by(() => {
    if (!schedule.length) return null;

    const targetId = clock.currentPeriodId;
    if (!targetId) return null;

    const entry = schedule.find(item => {
      const periodsStr = getPeriodsFromRange(item.start, item.end);
      if (periodsStr === "Svi sati") return true;
      const periodIds = periodsStr.split(',').map(Number);
      return periodIds.includes(targetId);
    });

    if (!entry) return null;

    // Strict time check: Is it actually happening right now?
    const now = clock.time;
    if (now < new Date(entry.start) || now > new Date(entry.end)) {
        return null;
    }

    return formatLessonData(entry);
  });

  // 3. Logic for NEXT (Array based)
  // We simply grab the next object in the list, ensuring we skip the current block entirely.
  let next = $derived.by(() => {
    if (!schedule.length) return null;

    // Case A: We are currently in a lesson
    if (current?.entry) {
      const currentIndex = schedule.indexOf(current.entry);
      // The next distinct block is just index + 1
      return formatLessonData(schedule[currentIndex + 1]);
    }

    // Case B: We are in a break or before school
    // Find the first lesson that starts in the future
    const now = clock.time;
    const upcomingEntry = schedule.find(item => new Date(item.start) > now);
    
    return formatLessonData(upcomingEntry);
  });
</script>

{#snippet lessonSlot(data: any)}
  {#if data}
    <div class="col-time">{getPeriodsFromRange(data.entry.start, data.entry.end)}</div>
    
    <div class="col-room">
      {#key data.transitionKey}
        <div in:fade={{ duration: 200 }} class="lesson-inner">
          {data.lesson.room}
        </div>
      {/key}
    </div>
    
    <div class="col-subject">
      {#key data.transitionKey}
        <div in:fade={{ duration: 200 }} out:fade={{ duration: 200 }} class="lesson-inner">
          <span class="group">{formatGroup(data.lesson.groups)}</span> 
          {data.lesson.subject.name}
        </div>
      {/key}
    </div>

    <div class="col-teacher">
      {#key data.transitionKey}
        <div in:fade={{ duration: 200 }} out:fade={{ duration: 200 }} class="lesson-inner">
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

    {@render lessonSlot(current)}

    <div class="col-divider line"></div>

    {@render lessonSlot(next)}

    <div class="col-class right-align">{formatClassName(name)}</div>
  {/if}
</div>
