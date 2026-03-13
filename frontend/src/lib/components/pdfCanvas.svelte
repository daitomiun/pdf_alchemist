<script lang="ts">
	import { loadPDF } from "$lib/pdfjs";
	import { Effect } from "effect";
	import PDFCarousel from "./PDFCarousel.svelte";

	import { pdfManager } from "$lib/state/state.svelte";
	import PDFViewer from "./PDFViewer.svelte";

	let files = $state<FileList | null>(null);

	let fileBuffer: ArrayBuffer;
	$effect(() => {
		if (files && files.length > 0) {
			const file = files[0];
			async function load() {
				fileBuffer = await Effect.runPromise(
					Effect.promise<ArrayBuffer>(() => file.arrayBuffer()),
				);
				let pdf = await Effect.runPromise(loadPDF(fileBuffer));
				pdfManager.setPdf(pdf);
				pdfManager.setInitialPageArr();
			}
			load();
		}
	});
</script>

{#if pdfManager.pdf.proxy == undefined}
	<input id="pdf-upload" bind:files type="file" accept=".pdf" />
{/if}
{#if files && pdfManager.pdf.proxy != undefined}
	<div class="container-grid">
		<div class="new-elements">
			<button
				aria-label="reset-changes"
				type="button"
				class="text-white bg-orange-500 box-border border border-transparent hover:bg-warning-strong focus:ring-4 focus:ring-warning-medium shadow-xs font-medium leading-5 rounded-base text-sm px-4 py-5 focus:outline-none"
				onclick={() => pdfManager.setInitialPageArr()}
			>
				Reset changes
			</button>
		</div>
		<PDFViewer></PDFViewer>
		<PDFCarousel></PDFCarousel>
	</div>
{/if}

<style>
	.container-grid {
		height: 100vh;
		overflow: hidden;
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		grid-template-rows: repeat(auto-fit, minmax(200px, 1fr));
		grid-column-gap: 0px;
		grid-row-gap: 0px;
	}

	.new-elements {
		grid-area: 1 / 3 / 4 / 6;
		border: 4px solid black;
		background-color: grey;
	}
</style>
