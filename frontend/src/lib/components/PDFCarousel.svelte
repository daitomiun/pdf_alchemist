<script lang="ts">
	import { pdfManager } from "$lib/state/state.svelte";
	import PDFPage from "./PDFPage.svelte";
	import type { PDFDocumentProxy } from "pdfjs-dist";
	import { flip } from "svelte/animate";

	type props = {
		pdf: PDFDocumentProxy;
		scale: number;
	};
	let { pdf, scale }: props = $props();

	$inspect(pdfManager.pageNumArr);
	const dragDuration = 300;
	let draggingCard: number | null = null;
	let animatingCards = new Set<number>();
</script>

<div class="pdf-carousel">
	{#each pdfManager.pageNumArr as pageNum (pageNum)}
		<div class="card" animate:flip={{ duration: dragDuration }}>
			<div
				role="list"
				class="pdf-widget"
				draggable="true"
				ondragstart={() => (draggingCard = pageNum)}
				ondragend={() => (draggingCard = null)}
				ondragenter={() =>
					pdfManager.swap(animatingCards, draggingCard, pageNum, dragDuration)}
			>
				<span>{pageNum}</span>
				<button
					onclick={() => pdfManager.delete(pageNum)}
					aria-label="delete-page"
					type="button"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class="size-4"
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
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
</style>
