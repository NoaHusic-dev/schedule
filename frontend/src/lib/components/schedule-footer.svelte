<script lang="ts">
  type Props = { text?: string; speed?: string };
  let { text = "Dobrodošli u Tehničku školu Daruvar!", speed = "20s" }: Props = $props();
</script>

<div
  class="relative flex h-(--row-h) w-full items-center overflow-hidden bg-accent text-row font-bold tracking-[0.05em] text-black uppercase"
>
  <div class="marquee absolute whitespace-nowrap" style:--duration={speed}>{text}</div>
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
