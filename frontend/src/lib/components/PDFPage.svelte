<script lang="ts">
	import type { PDFDocumentProxy } from "pdfjs-dist";
	import { renderPage } from "$lib/pdfjs";
	import { pdfManager } from "$lib/state/state.svelte";
	interface Props {
		pdf: PDFDocumentProxy;
		pageNum: number;
		scale: number;
	}
	let { pdf, pageNum, scale }: Props = $props();
</script>

<div class="page">
	<canvas
		onclick={() => pdfManager.setCurrentPage(pageNum)}
		use:renderPage={{ pdf, pageNum, scale }}
	></canvas>
</div>
