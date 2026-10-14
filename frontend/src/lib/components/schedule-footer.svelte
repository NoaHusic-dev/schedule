<script lang="ts">
  import RotationIndicator from "$lib/components/rotation-indicator.svelte";

  type Props = { page: number; pages: number; rotation: number; text?: string; speed?: string };
  let { page, pages, rotation, text = "Dobrodošli u Tehničku školu Daruvar!", speed = "20s" }: Props = $props();
</script>

<div
  class="flex h-(--row-h) w-full items-stretch overflow-hidden bg-bar text-row font-bold tracking-[0.05em] text-accent uppercase"
>
  <RotationIndicator {page} {pages} duration={rotation} />

  <!-- speed overrides .ticker's own 50s pass, keeping the pace the board had -->
  <div class="ticker flex-1 self-center">
    <div style:animation-duration={speed}>{text}</div>
  </div>
</div>

