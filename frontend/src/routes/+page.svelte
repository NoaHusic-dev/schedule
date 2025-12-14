<script lang="ts">
  import type { PageData } from "./$types";
  import Header from "$comp/ui/header.svelte";
  import ScheduleHeader from "$comp/ui/schedule-header.svelte";
  import ScheduleClass from "$comp/schedule-class.svelte";

  const { data }: { data: PageData } = $props<{ data: PageData }>();
</script>

<div class="app-container">
  <Header />

  <main class="grid-container">
    <ScheduleHeader />

    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />

    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />

    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />

    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />
    <ScheduleClass />
  </main>
</div>
