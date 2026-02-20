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

	let containerWidth = $state(0);
</script>

<div class="page" bind:clientWidth={containerWidth}>
	<canvas
		onclick={() => pdfManager.setCurrentPage(pageNum)}
		use:renderPage={{ pdf, pageNum, scale, containerWidth }}
	></canvas>
</div>
