<script lang="ts">
	import { pdfManager, type Page } from "$lib/state/state.svelte";
	import PDFPage from "./PDFPage.svelte";
	import { flip } from "svelte/animate";

	const dragDuration = 300;
	let draggingCard: Page | null = null;
	let animatingCards = new Set<Page>();

	let pendingCuts = $derived(pdfManager.pdf.pendingCuts);
	let pages = $derived(pdfManager.pdf.pages);
	let startCut = $derived(pdfManager.pdf.startCut);
	let endCut = $derived(pdfManager.pdf.endCut);
	$inspect(pendingCuts);
	//	$inspect(startCut);
	//	$inspect(endCut);
	$inspect(pages);
</script>

<div class="pdf-carousel">
	{#each pages as page, i (page)}
		<div class="pdf-group" animate:flip={{ duration: dragDuration }}>
			<div class="card-row">
				{#if i <= pdfManager.pdf.pages.length && pdfManager.pdf.pages.length > 1}
					<div class="spacer-group">
						<div
							role="button"
							tabindex={i}
							aria-label="split-pages"
							class="split-spacer {i == startCut || i == endCut
								? 'active'
								: 'inactive'}"
							onmousedown={() => pdfManager.split(i)}
						>
							-{i}-
						</div>
					</div>
				{/if}
				<div class="card">
					<div
						role="list"
						class="pdf-widget"
						draggable="true"
						ondragstart={() => (draggingCard = page)}
						ondragend={() => (draggingCard = null)}
						ondragenter={() =>
							pdfManager.swap(animatingCards, draggingCard, page, dragDuration)}
					>
						<span>{page.pageNum}</span>
						<button
							onmousedown={() => pdfManager.delete(page)}
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
					<PDFPage
						pdf={pdfManager.pdf.proxy}
						pageNum={page.pageNum}
						scale={page.scale}
					></PDFPage>
				</div>
				{#if i == pdfManager.pdf.pages.length - 1}
					<div class="spacer-group">
						<div
							role="button"
							tabindex={i + 1}
							aria-label="split-pages"
							class="split-spacer {i + 1 == startCut || i + 1 == endCut
								? 'active'
								: 'inactive'}"
							onmousedown={() => pdfManager.split(i + 1)}
						>
							{i + 1}
						</div>
					</div>
				{/if}
			</div>
			<div class="group-container">
				{#each pdfManager.pdf.pendingCuts as groupId, _ (groupId)}
					{#if page.groupIds.includes(groupId)}
						<div class="group">g</div>
					{:else}
						<div class="group-white">w</div>
					{/if}
				{/each}
			</div>
		</div>
	{/each}
</div>

<style>
	.group-white {
		border-color: 1px solid black;
		background-color: white;
	}
	.card-row {
		display: flex;
		flex-direction: row;
		align-items: center;
	}
	.group-container {
		grid-column: 1 / -1;
	}
	.group {
		background-color: lime;
		border: 1px solid black;
	}
	.pdf-group {
		display: grid;
		grid-template-columns: 1fr 1fr;
		grid-template-rows: auto 1fr;
	}
	.split-spacer {
		height: 100%;
		border: 1px dashed black;
	}
	.spacer-group {
		height: 100%;
	}
	.split-spacer:hover {
		color: green;
	}
	.pdf-carousel {
		grid-area: 4 / 3 / 6 / 6;
		display: flex;
		width: 100%;
		border: 5px black;
		overflow-x: auto;
		overflow-y: hidden;
		white-space: nowrap;
		background-color: aliceblue;
	}
	.inactive {
		color: violet;
	}
	.active {
		color: green;
	}
	.card {
		width: 100%;
		max-width: 100px;
		height: 100%;
		margin: 5px;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}
	.pdf-widget {
		width: 100%;
		height: 20px;
		background-color: gray;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	.pdf-widget:hover {
		cursor: grab;
	}
</style>
