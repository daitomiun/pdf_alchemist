import { Cause, Effect, Exit } from "effect";
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
	let activeTask: PDFJS.RenderTask | null = null;
	let currentId = 0;

	const render = async () => {
		const { pdf, pageNum, scale, containerWidth } = currentParams;

		const callId = ++currentId;
		if (activeTask) {
			activeTask.cancel();
			activeTask = null;
		}

		const page = await Effect.runPromise(getPage(pdf, pageNum));

		if (callId !== currentId) return; // INFO: stops svelte from rerendering the page, skipping the 'Use different canvas or ensure previous operations were cancelled or completed' error

		const unscaledViewport = page.getViewport({ scale: 1 });

		const fitScale = currentParams.containerWidth / unscaledViewport.width;
		const internalScale = fitScale * scale;
		const viewport = page.getViewport({ scale: internalScale });

		const context = node.getContext("2d");
		if (!context) return;

		node.height = viewport.height;
		node.width = viewport.width;

		node.style.width = `${containerWidth}px`;
		node.style.height = `${unscaledViewport.height * fitScale}px`;

		const renderContext: RenderParameters = {
			canvas: node,
			canvasContext: context,
			viewport,
		};
		const renderTask = page.render(renderContext);
		activeTask = renderTask;

		const renderEffect = Effect.promise(() => renderTask.promise);

		const result = await Effect.runPromiseExit(renderEffect);

		activeTask = null;

		if (Exit.isFailure(result)) {
			const cause = result.cause;
			if (Cause.isFailType(cause)) {
				const error: any = cause.error;
				if (error.name === "RenderingCancelledException") {
					console.log("Safe: Previous render was cancelled.");
				} else {
					console.error("The render failed with:", error);
				}
			}
			else if (Cause.isDieType(cause)) {
				console.error("The effect 'died' (crashed):", cause.defect);
			}
		}
	}
	render();
	return {
		update(newParams: { pdf: PDFJS.PDFDocumentProxy, pageNum: number, scale: number, containerWidth: number }) {
			currentParams = newParams;
			render();
		},
		destroy() {
			currentId++;
			if (activeTask) activeTask.cancel();
		}
	};
}


