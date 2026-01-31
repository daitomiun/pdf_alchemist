<script lang="ts">
  import PdfCarousel from "./lib/PdfCarousel.svelte";
</script>

<main>
  <div class="app-layout">
    <div class="pdf-view">test</div>
    <div class="file-attachment">test</div>
    <PdfCarousel></PdfCarousel>
  </div>
</main>

<style>
  .app-layout {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    grid-template-rows: repeat(5, 1fr);
  }
  .pdf-view {
    grid-area: 1 / 1 / 6 / 3;
    background-color: grey;
    border-color: black 3px;
    height: 1000px;
  }
  .file-attachment {
    background-color: white;
    border-color: black 3px;
    height: 700px;
    grid-area: 1 / 3 / 4 / 6;
  }
</style>
