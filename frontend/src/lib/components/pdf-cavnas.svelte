<script lang="ts">
	import { loadPDF, renderPage } from "$lib/pdfjs";
	import { Effect } from "effect";
	import type { PDFDocumentProxy } from "pdfjs-dist";

	let files = $state<FileList | null>(null);

	let pdf = $state<PDFDocumentProxy>();
	let fileBuffer: ArrayBuffer;
	let numPages = 0;
	$effect(() => {
		if (files && files.length > 0) {
			const file = files[0];
			async function load() {
				fileBuffer = await Effect.runPromise(
					Effect.promise<ArrayBuffer>(() => file.arrayBuffer()),
				);
				pdf = await Effect.runPromise(loadPDF(fileBuffer));
				numPages = pdf.numPages;
			}
			load();
		}
	});

	const pageNumbers = $derived(
		Array.from({ length: numPages }, (_, i) => i + 1),
	);
</script>

{#if pdf == undefined}
	<input id="pdf-upload" bind:files type="file" accept=".pdf" />
{/if}
{#if files && pdf != undefined}
	<div class="pdf-carousel">
		{#each pageNumbers as pageNum}
			<div class="page">
				<h3>Page {pageNum}</h3>
				<canvas use:renderPage={{ pdf, pageNum, scale: 1 }}></canvas>
			</div>
		{/each}
	</div>
{/if}
