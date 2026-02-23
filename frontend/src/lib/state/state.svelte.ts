import type { PDFDocumentProxy } from "pdfjs-dist";

class PdfManager {
	current = $state<PDFDocumentProxy>();
	currentPage = $state(1);
	pageNumArr = $state<number[]>([]);

	setPdf(doc: PDFDocumentProxy) {
		this.current = doc;
	}
	setCurrentPage(pageNum: number) {
		this.currentPage = pageNum;
	}
	setInitialPageArr() {
		this.pageNumArr = Array.from({ length: this.current!.numPages }, (_, i) => i + 1)
	}

	delete(page: number) {
		console.log(this.pageNumArr)
		this.pageNumArr.splice(this.pageNumArr.indexOf(page), 1)
		console.log(this.pageNumArr)
	}
	swap(animatingCards: Set<number>, draggingCard: number | null, card: number, dragDuration: number) {
		if (
			draggingCard === null ||
			draggingCard === card ||
			animatingCards.has(card)
		)
			return;

		const a = this.pageNumArr.indexOf(draggingCard);
		const b = this.pageNumArr.indexOf(card);
		if (a < 0 || b < 0) return;

		const next = [...this.pageNumArr];
		[next[a], next[b]] = [next[b], next[a]];
		this.pageNumArr = next;

		animatingCards.add(card);
		setTimeout(() => animatingCards.delete(card), dragDuration);
	}
}

export const pdfManager = new PdfManager()
