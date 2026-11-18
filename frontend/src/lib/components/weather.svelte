<script lang="ts">
  import { weather } from "#lib/stores/weather.svelte.ts";

  $effect(() => weather.subscribe());

  const icon = $derived(weatherIcon(weather.current?.weatherCode ?? -1));

  function weatherIcon(code: number): string {
    if (code === 0) return "☀️";
    if (code <= 2 && code >= 1) return "🌤️";
    if (code === 45 || code === 48) return "🌫️";
    if (code === 3) return "☁️";
    if (code >= 71 && code <= 77) return "❄️";
    if ((code >= 51 && code <= 67) || (code >= 80 && code <= 82)) return "🌧️";
    if (code === 85 || code === 86) return "🌨️";
    if (code >= 95 && code <= 99) return "⛈️";
    return "🌤️";
  }
</script>

<div
  class="flex items-center gap-[0.3vw] whitespace-nowrap text-[2.4vh] font-semibold text-muted"
  role="img"
  aria-label={weather.current
    ? `Trenutna temperatura u Daruvaru: ${weather.current.temperature} °C`
    : "Vremenski podaci trenutačno nisu dostupni"}
>
  {#if weather.current}
    <span class="text-[2.8vh]" aria-hidden="true">{icon}</span>
    <span>{weather.current.temperature}°C</span>
  {:else if weather.loaded}
    <span>—°C</span>
  {/if}
</div>
