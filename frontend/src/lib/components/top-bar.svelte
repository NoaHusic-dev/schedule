<script lang="ts">
  import Clock from "$lib/components/clock.svelte";
</script>

<header class="flex h-[10vh] items-center justify-between bg-board px-[2vw]">
  <div>
    <h1 class="mb-[0.2vh] text-brand font-extrabold tracking-[1px]">TEHNIČKA ŠKOLA DARUVAR</h1>
    <span class="text-subtitle font-semibold tracking-[0.2em] text-muted uppercase">Pregled nastave uživo</span>
  </div>

  <Clock />
</header>
