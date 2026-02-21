import { Effect } from "effect";
import * as PDFJS from "pdfjs-dist";
import pdfjsWorker from "pdfjs-dist/build/pdf.worker.min.mjs?url";
import type { RenderParameters } from "pdfjs-dist/types/src/display/api";

PDFJS.GlobalWorkerOptions.workerSrc = pdfjsWorker;

export function loadPDF(src: string | URL | ArrayBuffer): Effect.Effect<PDFJS.PDFDocumentProxy> {
	const loadingTask = PDFJS.getDocument(src);
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
		scale: number;
		containerWidth: number;
	}
) {
	let currentParams = params;

	const render = async () => {
		const { pdf, pageNum, scale } = currentParams;
		const page = await Effect.runPromise(getPage(pdf, pageNum));
		const unscaledViewport = page.getViewport({ scale: 1 });

		const fitScale = currentParams.containerWidth / unscaledViewport.width;
		const internalScale = fitScale * scale;
		const viewport = page.getViewport({ scale: internalScale });

		const context = node.getContext("2d");
		if (!context) return;

		node.height = viewport.height;
		node.width = viewport.width;

		node.style.width = `${currentParams.containerWidth}px`;
		node.style.height = `${unscaledViewport.height * fitScale}px`;

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
		update(newParams: { pdf: PDFJS.PDFDocumentProxy, pageNum: number, scale: number, containerWidth: number }) {
			currentParams = newParams;
			render();
		},
		destroy() { }
	};
}


