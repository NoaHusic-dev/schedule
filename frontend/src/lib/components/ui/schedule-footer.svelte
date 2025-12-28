<script lang="ts">
  let {
    text = "Dobrodošli u Tehničku školu Daruvar!",
    speed = "20s",
  }: { text?: string; speed?: string } = $props();
</script>

<div class="footer-row scroll-wrapper">
  <div class="marquee-content" style="--duration: {speed}">
    {text}
  </div>
</div>

<style>
  .scroll-wrapper {
    overflow: hidden;
    position: relative;
    display: flex !important;
    align-items: center;
    justify-content: flex-start;
    width: 100%;
  }

  .marquee-content {
    position: absolute;
    white-space: nowrap;
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

  .marquee-content:hover {
    animation-play-state: paused;
  }
</style>
