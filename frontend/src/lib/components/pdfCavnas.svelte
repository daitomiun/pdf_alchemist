<script lang="ts">
	import { loadPDF } from "$lib/pdfjs";
	import { Effect } from "effect";
	import PDFCarousel from "./PDFCarousel.svelte";

	import { pdfManager } from "$lib/state/state.svelte";
	import PDFViewer from "./PDFViewer.svelte";

	let files = $state<FileList | null>(null);

	let fileBuffer: ArrayBuffer;
	let numPages = 0;
	$effect(() => {
		if (files && files.length > 0) {
			const file = files[0];
			async function load() {
				fileBuffer = await Effect.runPromise(
					Effect.promise<ArrayBuffer>(() => file.arrayBuffer()),
				);
				let pdf = await Effect.runPromise(loadPDF(fileBuffer));
				pdfManager.setPdf(pdf);
				numPages = pdf.numPages;
			}
			load();
		}
	});

	$inspect(pdfManager.currentPage);

	const pageNumbers = $derived(
		Array.from({ length: numPages }, (_, i) => i + 1),
	);
</script>

{#if pdfManager.current == undefined}
	<input id="pdf-upload" bind:files type="file" accept=".pdf" />
{/if}
{#if files && pdfManager.current != undefined}
	<div class="container-grid">
		<div class="new-elements"></div>
		<PDFViewer></PDFViewer>
		<PDFCarousel pdf={pdfManager.current} scale={0.8} {pageNumbers}
		></PDFCarousel>
	</div>
{/if}

<style>
	.container-grid {
		display: grid;
		grid-template-columns: repeat(5, 1fr);
		grid-template-rows: repeat(5, 1fr);
		grid-column-gap: 0px;
		grid-row-gap: 0px;
	}

	.new-elements {
		grid-area: 1 / 3 / 4 / 6;
		border: 4px solid black;
		background-color: grey;
	}
</style>
