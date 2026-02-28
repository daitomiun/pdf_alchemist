import { ActionType, type Command } from "$lib/types/pdf";
import type { PDFDocumentProxy } from "pdfjs-dist";

class PdfManager {
	current = $state<PDFDocumentProxy>();
	currentPage = $state(1);
	pageNumArr = $state<number[]>([]);

	#undoStack: Command[] = []
	#redoStack: Command[] = []

	private execute(cmd: Command) {
		this.applyCommand(cmd);
		this.#undoStack.push(cmd);
		console.log(this.#undoStack)
		this.#redoStack = [];
	}

	private applyCommand(cmd: Command) {
		switch (cmd.type) {
			case ActionType.DELETE:
				this.pageNumArr = this.pageNumArr.filter((page) => page != cmd.pageId)
				console.log(this.pageNumArr)
				break;
			case ActionType.SWAP:
				if (cmd.pageIdA < 0 || cmd.pageIdB < 0) return;
				const next = [...this.pageNumArr];
				[next[cmd.pageIdA], next[cmd.pageIdB]] = [next[cmd.pageIdB], next[cmd.pageIdA]];
				this.pageNumArr = next;
				break;
			case ActionType.RESET:
				this.#undoStack = [];
				this.#redoStack = [];
				this.pageNumArr = Array.from({ length: this.current!.numPages }, (_, i) => i + 1);
				break;
		}
	}


	setPdf(doc: PDFDocumentProxy) {
		this.current = doc;
	}
	setCurrentPage(pageNum: number) {
		this.currentPage = pageNum;
	}
	setInitialPageArr() {
		this.execute({ type: ActionType.RESET })
	}

	delete(page: number) {
		const cmd: Command = {
			type: ActionType.DELETE,
			pageId: page,
			previousIndex: this.pageNumArr.indexOf(page)
		}
		this.execute(cmd)
	}

	swap(animatingCards: Set<number>, draggingCard: number | null, card: number, dragDuration: number) {
		if (
			draggingCard === null ||
			draggingCard === card ||
			animatingCards.has(card)
		)
			return;

		const cmd: Command = {
			type: ActionType.SWAP,
			pageIdA: this.pageNumArr.indexOf(draggingCard),
			pageIdB: this.pageNumArr.indexOf(card),
		}
		this.execute(cmd);

		animatingCards.add(card);
		setTimeout(() => animatingCards.delete(card), dragDuration);
	}
}

export const pdfManager = new PdfManager()

