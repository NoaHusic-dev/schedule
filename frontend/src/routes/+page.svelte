<script lang="ts">
  import { onMount } from 'svelte';
  import Header from "$comp/ui/header.svelte";
  import ScheduleHeader from "$comp/ui/schedule-header.svelte";
  import ScheduleFooter from "$comp/ui/schedule-footer.svelte";
  import ScheduleClass from "$comp/schedule-class.svelte";

const sets = [
    [
      "1 EL", "2 EL", "3 EL", "4EL", 
      "1 MT", "2 MT", "3 MT", "4 MT", 
      "1 PT", "2 PT", "3 PT", "4 PT", 
      "1 RT", "2 RT", "3 RT", "4 RT"
    ],
    [
      "1 K", "2 K", "3 K", "1 P", 
      "2 P", "3 P", "1 S", "2 S", 
      "3 S", "1 SE", "2SA", "2 SP", 
      "3 SP", "1M", "2M", "3M"
    ]
  ];

  let currentIndex = 0;
  const ROTATION_SPEED = 10000;

  $: currentClasses = sets[currentIndex];

  onMount(() => {
    const interval = setInterval(() => {
      currentIndex = (currentIndex + 1) % sets.length;
    }, ROTATION_SPEED);

    return () => clearInterval(interval);
  });
</script>

<div class="app-container">
  <Header />

  <main class="grid-container">
    <ScheduleHeader />

    {#each currentClasses as className, index}
      <ScheduleClass name={className} />
    {/each}

    <ScheduleFooter />
  </main>
</div>
