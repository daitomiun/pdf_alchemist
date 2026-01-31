<script lang="ts">
  let { card } = $props();
</script>

<div class="pdf-card">
  {card}
</div>

<style>
  .pdf-card {
    width: 70px;
    height: 120px;
    background-color: white;
    color: black;
  }
</style>
