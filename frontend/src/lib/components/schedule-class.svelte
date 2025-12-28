<script lang="ts">
  import { fade } from 'svelte/transition';
  import { scheduleStore } from "$lib/stores/schedules.svelte";
  import { getPeriodsFromRange } from "$util/timetable";
  import { formatClassName, formatGroup } from "$util/class";

  let { name }: { name: string } = $props();

  let classData = $derived(scheduleStore.classes[name]);
  let schedule = $derived(classData?.schedule ?? []);
  
  let tick = $state(0);

  $effect(() => {
    const interval = setInterval(() => {
      tick++;
    }, 5000);
    return () => clearInterval(interval);
  });

  function getLessonData(periodIndex: number) {
    const period = schedule[periodIndex];
    if (!period || period.lessons.length === 0) return null;
    
    const lessonIdx = tick % period.lessons.length;
    const lesson = period.lessons[lessonIdx];

    return {
      period,
      lesson,
      transitionKey: `${periodIndex}-${lesson.subject.short}-${lesson.groups?.join('')}`
    };
  }

  let current = $derived(getLessonData(0));
  let next = $derived(getLessonData(1));
</script>

{#snippet lessonSlot(data: any)}
  {#if data}
    <div class="col-time">{getPeriodsFromRange(data.period.start, data.period.end)}</div>
    <div class="col-room">{data.lesson.room}</div>
    
    <div class="col-subject">
      {#key data.transitionKey}
        <div in:fade={{ duration: 200 }} out:fade={{ duration: 200 }}>
          <span class="group">{formatGroup(data.lesson.groups)}</span> 
          {data.lesson.subject.name}
        </div>
      {/key}
    </div>

    <div class="col-teacher">{data.lesson.teacher.short}</div>
  {:else}
    <div class="col-empty"></div>
  {/if}
{/snippet}

<div class="grid-row">
  {#if scheduleStore.isLoading && !classData}
    <div class="loading-state">...</div>
  {:else if classData}
    <div class="col-class">{formatClassName(name)}</div>

    {@render lessonSlot(current)}

    <div class="col-divider line"></div>

    {@render lessonSlot(next)}

    <div class="col-class right-align">{formatClassName(name)}</div>
  {/if}
</div>
