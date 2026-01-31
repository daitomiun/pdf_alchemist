<script lang="ts">
  let { card } = $props();
</script>

<div class="pdf-card">
  {card}
</div>

<style>
  .pdf-card {
    width: 80px;
    height: 100%;
    background-color: white;
    color: black;
  }
</style>
