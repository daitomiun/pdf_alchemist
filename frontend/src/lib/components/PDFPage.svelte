<script lang="ts">
	import type { PDFDocumentProxy } from "pdfjs-dist";
	import { renderPage } from "$lib/pdfjs";
	interface Props {
		pdf: PDFDocumentProxy;
		pageNum: number;
		scale: number;
	}
	let { pdf, pageNum, scale }: Props = $props();
</script>

<div class="page">
	<h3>Page {pageNum}</h3>
	<canvas use:renderPage={{ pdf, pageNum, scale }}></canvas>
</div>
