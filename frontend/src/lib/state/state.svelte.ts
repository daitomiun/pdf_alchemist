import type { PDFDocumentProxy } from "pdfjs-dist";

class PdfManager {
	current = $state<PDFDocumentProxy>();
	currentPage = $state(1);

	setPdf(doc: PDFDocumentProxy) {
		this.current = doc;
	}
	setCurrentPage(pageNum: number) {
		this.currentPage = pageNum;
	}
}

export const pdfManager = new PdfManager()
