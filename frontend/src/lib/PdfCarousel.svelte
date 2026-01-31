<script lang="ts">
  import PdfCard from "./PdfCard.svelte";
  let cards = Array(20)
    .fill(null)
    .map((_, i) => i + 1);
</script>

<div class="pdf-editor">
  {#each cards as card (card)}
    <PdfCard {card}></PdfCard>
  {/each}
</div>

<style>
  .pdf-editor {
    display: grid;
    grid-auto-flow: column;
    gap: 5px;
    width: 700px;
    height: 120px;
    margin: 1rem 2rem;
    overflow-y: scroll;
    grid-area: 4 / 3 / 6 / 6;
  }
</style>
