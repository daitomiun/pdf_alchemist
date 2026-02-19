<script lang="ts">
	import { renderPage } from "$lib/pdfjs";
	import { pdfManager } from "$lib/state/state.svelte";

	let pdf = $derived(pdfManager.current);
	let currentPage = $derived(pdfManager.currentPage);
	$inspect(currentPage);
	$inspect(pdf);
</script>

{#if pdf != undefined}
	<div class="pdf-viewer">
		<canvas
			use:renderPage={{
				pdf: pdf,
				pageNum: currentPage,
				scale: 2,
			}}
		></canvas>
	</div>
{:else}
	<div class="error">could not render the page</div>
{/if}
