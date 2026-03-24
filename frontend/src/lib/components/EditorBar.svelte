<script lang="ts">
	import { pdfManager } from "$lib/state/state.svelte";
</script>

<div class="editor-bar">
	<button
		aria-label="reset-changes"
		type="button"
		class="text-white bg-purple-500 box-border border border-transparent hover:bg-warning-strong focus:ring-4 focus:ring-warning-medium shadow-xs font-medium leading-5 rounded-base text-sm px-4 py-5 focus:outline-none"
		onclick={() => pdfManager.export()}
	>
		Undo
	</button>
	<button
		aria-label="reset-changes"
		type="button"
		class="text-white bg-purple-500 box-border border border-transparent hover:bg-warning-strong focus:ring-4 focus:ring-warning-medium shadow-xs font-medium leading-5 rounded-base text-sm px-4 py-5 focus:outline-none"
		onclick={() => pdfManager.export()}
	>
		Do
	</button>
	<span>---</span>
	<button
		aria-label="reset-changes"
		type="button"
		class="text-white bg-blue-500 box-border border border-transparent hover:bg-warning-strong focus:ring-4 focus:ring-warning-medium shadow-xs font-medium leading-5 rounded-base text-sm px-4 py-5 focus:outline-none"
		onclick={() => pdfManager.resetChanges()}
	>
		Reset changes
	</button>
	<button
		aria-label="reset-changes"
		type="button"
		class="text-white bg-purple-500 box-border border border-transparent hover:bg-warning-strong focus:ring-4 focus:ring-warning-medium shadow-xs font-medium leading-5 rounded-base text-sm px-4 py-5 focus:outline-none"
		onclick={() => pdfManager.export()}
	>
		Export Changes
	</button>
</div>
