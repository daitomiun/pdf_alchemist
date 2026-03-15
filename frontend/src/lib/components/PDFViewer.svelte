<script lang="ts">
	import { renderPage } from "$lib/pdfjs";
	import { pdfManager } from "$lib/state/state.svelte";

	let pdf = $derived(pdfManager.pdf.proxy);
	let currentPage = $derived(pdfManager.pdf.currentPage);
	let containerWidth = $state(0);

	$inspect(currentPage);

	$effect(() => {
		const isPageVisible = pdfManager.pdf.pages.some(
			(p) => p.pageNum === currentPage,
		);
		console.log(`is visible -> ${isPageVisible}`);
		if (!isPageVisible) {
			console.log(pdfManager.pdf.pages);
			console.log(pdfManager.pdf.pages.findIndex((p) => p.pageNum === 1));
			pdfManager.setCurrent(pdfManager.pdf.pages[0].pageNum);
		}
	});
</script>

{#if pdf != undefined && pdfManager.pdf.pages.length > 0}
	<div class="pdf-padding">
		<div class="pdf-viewer" bind:clientWidth={containerWidth}>
			<canvas
				use:renderPage={{
					pdf: pdf,
					pageNum: currentPage,
					scale: pdfManager.pdf.viewerScale,
					containerWidth: containerWidth,
				}}
			></canvas>
		</div>
	</div>
{:else}
	<div class="error">There's no more pages :(</div>
{/if}

<style>
	.pdf-padding {
		width: 100%;
		height: 100%;
		padding: 3.5rem;
		background-color: lightseagreen;
		display: flex;
		align-items: center;
		justify-content: center;
		grid-area: 1 / 1 / 6 / 3;
	}
	.pdf-viewer {
		border: 5px solid black;
		width: 100%;
	}
</style>
