<script lang="ts">
	import PDFPage from "./PDFPage.svelte";
	import type { PDFDocumentProxy } from "pdfjs-dist";
	import { flip } from "svelte/animate";

	type props = {
		pdf: PDFDocumentProxy;
		scale: number;
		pageNumbers: number[];
	};
	let { pdf, scale, pageNumbers }: props = $props();

	const dragDuration = 300;
	let draggingCard: number | null = null;
	let animatingCards = new Set<number>();

	function swapWith(card: number) {
		if (
			draggingCard === null ||
			draggingCard === card ||
			animatingCards.has(card)
		)
			return;

		const a = pageNumbers.indexOf(draggingCard);
		const b = pageNumbers.indexOf(card);
		if (a < 0 || b < 0) return;

		const next = [...pageNumbers];
		[next[a], next[b]] = [next[b], next[a]];
		pageNumbers = next;

		animatingCards.add(card);
		setTimeout(() => animatingCards.delete(card), dragDuration);
	}
</script>

<div class="pdf-carousel">
	{#each pageNumbers as pageNum (pageNum)}
		<div class="card" animate:flip={{ duration: dragDuration }}>
			<div
				role="list"
				class="pdf-widget"
				draggable="true"
				ondragstart={() => (draggingCard = pageNum)}
				ondragend={() => (draggingCard = null)}
				ondragenter={() => swapWith(pageNum)}
			>
				{pageNum}
			</div>
			<PDFPage {pdf} {pageNum} {scale}></PDFPage>
		</div>
	{/each}
</div>

<style>
	.pdf-carousel {
		display: flex;
		width: 700px;
		border: 5px black;
	}
	.card {
		width: 100%;
		height: 100%;
	}
	.pdf-widget {
		width: 100%;
		height: 20px;
		background-color: gray;
	}
</style>
