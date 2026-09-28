<script lang="ts">
  import "./layout.css";
  import { updated } from "$app/state";
  import favicon from "$lib/assets/favicon.svg";
  import type { LayoutProps } from "./$types";

  let { children }: LayoutProps = $props();

  // Nobody touches the TV: when a new build is deployed (polled every 5 min, see vite.config.ts),
  // reload so it runs the new code. A full reload is the point here, so no goto().
  $effect(() => {
    if (updated.current) location.reload();
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>Tehnička škola Daruvar · Nastava</title>
</svelte:head>

{@render children()}
