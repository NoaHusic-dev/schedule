<script lang="ts">
  import { scheduleStore } from "$lib/stores/schedules.svelte";

  let { name }: { name: string } = $props();

  let classData = $derived(scheduleStore.classes[name]);
</script>

<div class="grid-row">
  {#if scheduleStore.isLoading && !classData}
    <p>Loading schedule...</p>
  {:else if classData}
    <div class="col-class">{name}</div>
    <div class="col-time">{classData.schedule[0].start}</div>
    <div class="col-room">{classData.schedule[0].lessons[0].room}</div>
    <div class="col-subject">{classData.schedule[0].lessons[0].subject.name}</div>
    <div class="col-teacher">{classData.schedule[0].lessons[0].teacher.short}</div>

    <div class="col-divider line"></div>

    <div class="col-time">{classData.schedule[1].start}</div>
    <div class="col-room">{classData.schedule[1].lessons[0].room}</div>
    <div class="col-subject">{classData.schedule[1].lessons[0].subject.name}</div>
    <div class="col-teacher">{classData.schedule[1].lessons[0].teacher.short}</div>
    <div class="col-class right-align">{name}</div>
  {:else}
    <p></p>
  {/if}
</div>
