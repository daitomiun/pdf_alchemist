<script lang="ts">
	import PDFPage from "./PDFPage.svelte";
	import type { PDFDocumentProxy } from "pdfjs-dist";

	type props = {
		pdf: PDFDocumentProxy;
		scale: number;
		pageNumbers: number[];
	};
	let { pdf, scale, pageNumbers }: props = $props();
</script>

<div class="pdf-carousel">
	{#each pageNumbers as pageNum}
		<PDFPage {pdf} {pageNum} {scale}></PDFPage>
	{/each}
</div>
