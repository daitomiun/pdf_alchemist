<script lang="ts">
	import { loadPDF, renderPage } from "$lib/pdfjs";
	import { Effect } from "effect";
	import type { PDFDocumentProxy } from "pdfjs-dist";

	let fileBuffer: ArrayBuffer;
	async function getFile(event: Event) {
		const target = event.target as HTMLInputElement;
		const fileList = target.files;
		if (fileList && fileList.length > 0) {
			const file = fileList[0];
			fileBuffer = await Effect.runPromise(
				Effect.promise<ArrayBuffer>(() => file.arrayBuffer()),
			);
		}
	}

	let pdf: PDFDocumentProxy;
	let numPages = 0;
	async function load() {
		pdf = await Effect.runPromise(loadPDF(fileBuffer));
		numPages = pdf.numPages;
	}
	$: pageNumbers = Array.from({ length: numPages }, (_, i) => i + 1);
</script>

<button on:click={load}>Load PDF</button>
<input id="pdf-upload" on:change={getFile} type="file" accept=".pdf" />
<div class="col-span-full">
	<label for="cover-photo" class="block text-sm/6 font-medium text-gray-900"
		>Cover photo</label
	>
	<div
		class="mt-2 flex justify-center rounded-lg border border-dashed border-gray-900/25 px-6 py-10"
	>
		<div class="text-center">
			<svg
				viewBox="0 0 24 24"
				fill="currentColor"
				data-slot="icon"
				aria-hidden="true"
				class="mx-auto size-12 text-gray-300"
			>
				<path
					d="M1.5 6a2.25 2.25 0 0 1 2.25-2.25h16.5A2.25 2.25 0 0 1 22.5 6v12a2.25 2.25 0 0 1-2.25 2.25H3.75A2.25 2.25 0 0 1 1.5 18V6ZM3 16.06V18c0 .414.336.75.75.75h16.5A.75.75 0 0 0 21 18v-1.94l-2.69-2.689a1.5 1.5 0 0 0-2.12 0l-.88.879.97.97a.75.75 0 1 1-1.06 1.06l-5.16-5.159a1.5 1.5 0 0 0-2.12 0L3 16.061Zm10.125-7.81a1.125 1.125 0 1 1 2.25 0 1.125 1.125 0 0 1-2.25 0Z"
					clip-rule="evenodd"
					fill-rule="evenodd"
				/>
			</svg>
			<div class="mt-4 flex text-sm/6 text-gray-600">
				<label
					for="file-upload"
					class="relative cursor-pointer rounded-md bg-transparent font-semibold text-indigo-600 focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-indigo-600 hover:text-indigo-500"
				>
					<span>Upload a file</span>
					<input
						id="file-upload"
						type="file"
						name="file-upload"
						accept=".pdf"
						class="sr-only"
					/>
				</label>
				<p class="pl-1">or drag and drop</p>
			</div>
			<p class="text-xs/5 text-gray-600">PNG, JPG, GIF up to 10MB</p>
		</div>
	</div>
</div>

{#each pageNumbers as pageNum}
	<div class="page">
		<h3>Page {pageNum}</h3>
		<canvas use:renderPage={{ pdf, pageNum }}></canvas>
	</div>
{/each}
