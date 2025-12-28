<script lang="ts">
  import Header from "$comp/ui/header.svelte";
  import ScheduleHeader from "$comp/ui/schedule-header.svelte";
  import ScheduleFooter from "$comp/ui/schedule-footer.svelte";
  import ScheduleClass from "$comp/schedule-class.svelte";
</script>

<div class="app-container">
  <Header />

  <main class="grid-container">
    <ScheduleHeader />

    <ScheduleClass name="1 EL" />
    <ScheduleClass name="2 EL" />
    <ScheduleClass name="3 EL" />
    <ScheduleClass name="4EL" />

    <ScheduleClass name="1 MT" />
    <ScheduleClass name="2 MT" />
    <ScheduleClass name="3 MT" />
    <ScheduleClass name="4 MT" />

    <ScheduleClass name="1 PT" />
    <ScheduleClass name="1 PT" />
    <ScheduleClass name="2 PT" />
    <ScheduleClass name="3 PT" />
    <ScheduleClass name="4 PT" />

    <ScheduleClass name="1 RT" />
    <ScheduleClass name="2 RT" />
    <ScheduleClass name="3 RT" />
    <ScheduleClass name="4 RT" />

    <ScheduleFooter />
  </main>
</div>
