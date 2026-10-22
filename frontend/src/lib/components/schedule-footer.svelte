<script lang="ts">
  import RotationIndicator from "$lib/components/rotation-indicator.svelte";

  type Props = { page: number; pages: number; rotation: number; text?: string; speed?: string };

  const MESSAGES = [
    "Dobrodošli u Tehničku školu Daruvar",
    "Znatiželja, upornost i točnost",
    "Znanje koje danas steknete sutra pokreće strojeve, mreže i ideje",
    "Nema glupih pitanja – nastavnici su tu zbog vas",
    "Mali koraci svaki dan vode do velikih rezultata",
    "Ugodan dan u školi i puno uspjeha na satu!",
  ];
  const SEPARATOR = "\u00a0\u00a0✦\u00a0\u00a0";

  let { page, pages, rotation, text = MESSAGES.join(SEPARATOR), speed }: Props = $props();

  // A longer text needs a longer pass, or it would scroll too fast to read.
  const duration = $derived(speed ?? `${Math.max(30, Math.round(text.length * 0.16))}s`);
</script>

<div
  class="flex h-(--row-h) w-full items-stretch overflow-hidden bg-bar text-row font-bold tracking-[0.05em] text-accent uppercase"
>
  <RotationIndicator {page} {pages} duration={rotation} />

  <!-- duration overrides .ticker's own 50s pass; pass `speed` to force a pace -->
  <div class="ticker flex-1 self-center">
    <div style:animation-duration={duration}>{text}</div>
  </div>
</div>