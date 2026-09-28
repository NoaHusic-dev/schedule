<script lang="ts">
  import RotationIndicator from "$lib/components/rotation-indicator.svelte";

  type Props = { page: number; pages: number; rotation: number; text?: string; speed?: string };
  let { page, pages, rotation, text = "Dobrodošli u Tehničku školu Daruvar!", speed = "20s" }: Props = $props();
</script>

<div
  class="flex h-(--row-h) w-full items-stretch overflow-hidden bg-bar text-row font-bold tracking-[0.05em] text-accent uppercase"
>
  <RotationIndicator {page} {pages} duration={rotation} />

  <div class="relative flex flex-1 items-center overflow-hidden">
    <div class="marquee absolute whitespace-nowrap" style:--duration={speed}>{text}</div>
  </div>
</div>

<style>
  .marquee {
    will-change: transform;
    animation: scroll-linear var(--duration) linear infinite;
  }

  @keyframes scroll-linear {
    from {
      transform: translateX(100vw);
    }
    to {
      transform: translateX(-100%);
    }
  }
</style>
