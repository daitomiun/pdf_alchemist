import { ActionType, type Command } from "$lib/types/pdf";
import type { PDFDocumentProxy } from "pdfjs-dist";
import { generate } from "short-uuid";


class PdfManager {
	current = $state<PDFDocumentProxy>();
	currentPage = $state(1);
	pageNumArr = $state<number[]>([]);
	pendingCuts = $state<Map<string, number[]>>();

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
				this.pendingCuts = new Map();
				this.pageNumArr = Array.from({ length: this.current!.numPages }, (_, i) => i + 1);
				break;
			case ActionType.SPLIT:
				const groupedPages = new Set(
					[...this.pendingCuts!.values()].flat()
				)
				const availablePages = this.pageNumArr.filter((p) => !groupedPages.has(p))
				// INFO: All the available pages that are not in the split groups  
				console.log(availablePages)

				const cut = availablePages.indexOf(this.pageNumArr[cmd.splitAt - 1]);
				console.log(`cut --> ${cut}`)
				const groupList = availablePages.slice(0, cut);
				console.log(`groupList -> ${groupList}`)
				// TODO: 1. get the splitAt (ie: nextNeighbor) and page id
				// 2. From the splitAt read backwards on the page list (ie: list -> [1,2,3,4,5]; splitAt -> 3; split group [1,2] )
				// 3. set the push with a new string id map
				// 4. call and execute the action from the user action change
				// NOTE: the UI should show a hightlight background color showing the distinction, after each render the list will update the styles for the desired group
				this.pendingCuts!.set(cmd.groupId, groupList)
				console.log(this.pendingCuts)
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

	split(splitAt: number) {
		const cmd: Command = {
			type: ActionType.SPLIT,
			splitAt: splitAt,
			groupId: generate(),
		}
		this.execute(cmd);
	}
}

export const pdfManager = new PdfManager()

