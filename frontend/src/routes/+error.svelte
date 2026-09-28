<script lang="ts">
  import { page } from "$app/state";
</script>

<svelte:head><title>Pogreška {page.status}</title></svelte:head>

<div class="flex h-full flex-col items-center justify-center gap-[2vh] text-center">
  <p class="text-subtitle font-semibold tracking-[0.2em] text-muted uppercase">Pogreška {page.status}</p>
  <h1 class="text-clock font-extrabold">Stranica nije pronađena.</h1>
  <a href="/" class="text-row text-accent underline">Natrag na raspored</a>
</div>
