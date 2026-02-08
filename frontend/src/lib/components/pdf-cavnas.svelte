<script lang="ts">
	import { loadPDF, renderPage } from "$lib/pdfjs";
	import { Effect } from "effect";
	import type { PDFDocumentProxy } from "pdfjs-dist";

	let pdf: PDFDocumentProxy;
	let numPages = 0;
	async function load() {
		pdf = await Effect.runPromise(loadPDF("../../static/sample-local-pdf.pdf"));
		numPages = pdf.numPages;
	}
	$: pageNumbers = Array.from({ length: numPages }, (_, i) => i + 1);
</script>

<button on:click={load}>Load PDF</button>

{#each pageNumbers as pageNum}
	<div class="page">
		<h3>Page {pageNum}</h3>
		<canvas use:renderPage={{ pdf, pageNum }}></canvas>
	</div>
{/each}
