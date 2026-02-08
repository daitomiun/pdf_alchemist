import { Effect } from "effect";
import * as PDFJS from "pdfjs-dist";
import pdfjsWorker from "pdfjs-dist/build/pdf.worker.min.mjs?url";
import type { RenderParameters } from "pdfjs-dist/types/src/display/api";

// Configure worker
PDFJS.GlobalWorkerOptions.workerSrc = pdfjsWorker;

export default function loadPDFWithEffect(node: HTMLCanvasElement, url: string) {
	const render = async () => {
		const loadingTask = PDFJS.getDocument(url);
		const pdf = await Effect.runPromise(
			Effect.promise<PDFJS.PDFDocumentProxy>(() => loadingTask.promise)
		);
		const page = await Effect.runPromise(
			Effect.promise<PDFJS.PDFPageProxy>(() => pdf.getPage(3))
		);
		const scale = 1;
		const viewport = page.getViewport({ scale });

		const context = node.getContext("2d");
		if (!context) return;

		node.height = viewport.height;
		node.width = viewport.width;

		const renderContext: RenderParameters = {
			canvas: node,
			canvasContext: context,
			viewport,
		};

		const render = Effect.promise(() => page.render(renderContext).promise);
		await Effect.runPromise(render);
	}

	render();
	return {
		destroy() { }
	};
}
