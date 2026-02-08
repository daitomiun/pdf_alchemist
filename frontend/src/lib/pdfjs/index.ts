import { Effect } from "effect";
import * as PDFJS from "pdfjs-dist";
import pdfjsWorker from "pdfjs-dist/build/pdf.worker.min.mjs?url";
import type { RenderParameters } from "pdfjs-dist/types/src/display/api";

// Configure worker
PDFJS.GlobalWorkerOptions.workerSrc = pdfjsWorker;

export function loadPDF(url: string): Effect.Effect<PDFJS.PDFDocumentProxy> {
	const loadingTask = PDFJS.getDocument(url);
	return Effect.promise<PDFJS.PDFDocumentProxy>(() => loadingTask.promise)
}

export function getPage(pdf: PDFJS.PDFDocumentProxy, pageIdx: number): Effect.Effect<PDFJS.PDFPageProxy> {
	return Effect.promise<PDFJS.PDFPageProxy>(() => pdf.getPage(pageIdx));
}

export function renderPage(
	node: HTMLCanvasElement,
	params: {
		pdf: PDFJS.PDFDocumentProxy;
		pageNum: number;
	}
) {
	const render = async () => {
		const { pdf, pageNum } = params;
		const page = await Effect.runPromise(getPage(pdf, pageNum));
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


