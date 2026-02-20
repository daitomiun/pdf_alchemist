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
				<button
					aria-label="delete-page"
					type="button"
					class="text-white bg-danger box-border border border-transparent hover:bg-danger-strong focus:ring-4 focus:ring-danger-medium shadow-xs font-medium leading-5 rounded-base text-sm px-4 py-2.5 focus:outline-none"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class="size-6"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M6 18 18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>
			<PDFPage {pdf} {pageNum} {scale}></PDFPage>
		</div>
	{/each}
</div>

<style>
	.pdf-carousel {
		grid-area: 4 / 3 / 6 / 6;
		display: flex;
		width: 100%;
		border: 5px black;
		overflow-x: auto;
		white-space: nowrap;
		background-color: aliceblue;
	}
	.card {
		width: 100%;
		height: 100%;
		margin: 5px;
	}
	.pdf-widget {
		width: 100%;
		height: 20px;
		background-color: gray;
	}
</style>
