<script lang="ts">
  type Props = { page: number; pages: number; duration: number };
  let { page, pages, duration }: Props = $props();
</script>

<!-- Tells readers which page is up and when theirs comes back. The bar restarts with each page. -->
<div class="flex h-full flex-none items-center gap-[0.7vw] bg-bar-alt px-[1.4vw] text-accent tabular-nums">
  <span>{page + 1} / {pages}</span>
  <div class="h-[0.55vh] w-[5.7vw] overflow-hidden rounded-full bg-line">
    {#key page}
      <div class="h-full animate-fill bg-accent" style:--fill-duration="{duration}ms"></div>
    {/key}
  </div>
</div>
