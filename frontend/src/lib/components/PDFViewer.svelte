<script lang="ts">
	import { renderPage } from "$lib/pdfjs";
	import { pdfManager } from "$lib/state/state.svelte";

	let pdf = $derived(pdfManager.current);
	let currentPage = $derived(pdfManager.currentPage);
	let containerWidth = $state(0);

	$effect(() => {
		const isPageVisible = pdfManager.pageNumArr.includes(currentPage);
		if (!isPageVisible) {
			pdfManager.setCurrentPage(pdfManager.pageNumArr[0]);
		}
	});
</script>

{#if pdf != undefined && pdfManager.pageNumArr.length > 0}
	<div class="pdf-viewer" bind:clientWidth={containerWidth}>
		<canvas
			use:renderPage={{
				pdf: pdf,
				pageNum: currentPage,
				scale: 2,
				containerWidth: containerWidth,
			}}
		></canvas>
	</div>
{:else}
	<div class="error">could not render the page</div>
{/if}

<style>
	.pdf-viewer {
		grid-area: 1 / 1 / 6 / 3;
		border: 5px solid black;
		background-color: lightseagreen;
		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
