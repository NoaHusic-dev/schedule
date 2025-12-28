<script lang="ts">
  import { scheduleStore } from "$lib/stores/schedules.svelte";
  import { getPeriodsFromRange } from "$util/timetable";
  import { formatClassName } from "$util/class";
  import { formatGroup } from "$util/class";

  let { name }: { name: string } = $props();

  let classData = $derived(scheduleStore.classes[name]);
</script>

<div class="grid-row">
  {#if scheduleStore.isLoading && !classData}
    <p>Loading schedule...</p>
  {:else if classData}
    <div class="col-class">{formatClassName(name)}</div>
    <div class="col-time">{getPeriodsFromRange(classData.schedule[0]?.start, classData.schedule[0]?.end)}</div>
    <div class="col-room">{classData.schedule[0]?.lessons[0].room}</div>
    <div class="col-subject">{formatGroup(classData.schedule[0]?.lessons[0].groups)} {classData.schedule[0]?.lessons[0].subject.name}</div>
    <div class="col-teacher">{classData.schedule[0]?.lessons[0].teacher.short}</div>

    <div class="col-divider line"></div>

    <div class="col-time">{getPeriodsFromRange(classData.schedule[1]?.start, classData.schedule[1]?.end)}</div>
    <div class="col-room">{classData.schedule[1]?.lessons[0].room}</div>
    <div class="col-subject">{formatGroup(classData.schedule[1]?.lessons[0].groups)} {classData.schedule[1]?.lessons[0].subject.name}</div>
    <div class="col-teacher">{classData.schedule[1]?.lessons[0].teacher.short}</div>
    <div class="col-class right-align">{formatClassName(name)}</div>
  {:else}
    <p></p>
  {/if}
</div>
