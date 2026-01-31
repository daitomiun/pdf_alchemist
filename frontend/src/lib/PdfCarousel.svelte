<script lang="ts">
  import PdfCard from "./PdfCard.svelte";
  let cards = Array(20)
    .fill(null)
    .map((_, i) => i + 1);
</script>

<div class="carousel">
  {#each cards as card (card)}
    <PdfCard {card}></PdfCard>
  {/each}
</div>

<style>
  .carousel {
    display: grid;
    grid-auto-flow: column;
    gap: 5px;
    width: 50px;
    margin: 1rem 2rem;
  }
</style>



