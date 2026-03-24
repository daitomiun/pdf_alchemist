import { ActionType, type Command } from "$lib/types/pdf";
import type { PDFDocumentProxy } from "pdfjs-dist";
import { generate } from "short-uuid";


class PdfManager {
	pdf = $state<Pdf>({
		pages: [],
		carouselScale: 0.8,
		viewerScale: 2,
		startCut: -1,
		endCut: -1,
		currentPage: 1,
		pendingCuts: []
	});

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
				cmd.page.groupIds.forEach((id) => {
					const pageIdx = this.pdf.pages.indexOf(cmd.page);
					const groupIdx = this.pdf.pages[pageIdx].groupIds.indexOf(id);
					this.pdf.pages[pageIdx].groupIds.splice(groupIdx, 1);
				});
				this.pdf.pages = this.pdf.pages.filter((page) => page.pageNum != cmd.page.pageNum)
				break;
			case ActionType.SWAP:
				const pageIdA = this.pdf.pages.indexOf(cmd.pageA)
				const pageIdB = this.pdf.pages.indexOf(cmd.pageB)
				if (pageIdA < 0 || pageIdB < 0) return;
				const next = [...this.pdf.pages];
				const groupIdPageA = next[pageIdA].groupIds
				const groupIdPageB = next[pageIdB].groupIds
				next[pageIdA].groupIds = groupIdPageB;
				next[pageIdB].groupIds = groupIdPageA;
				[next[pageIdA], next[pageIdB]] = [next[pageIdB], next[pageIdA]];

				this.pdf.pages = next;
				break;
			case ActionType.RESET:
				this.#undoStack = [];
				this.#redoStack = [];
				this.pdf.pages = Array.from(
					{ length: this.pdf.proxy!.numPages },
					(_, i) => ({
						id: i,
						pageNum: i + 1,
						groupIds: [],
						scale: this.pdf.carouselScale,
					})
				);
				this.pdf.startCut = -1;
				this.pdf.endCut = -1;
				this.pdf.pendingCuts = []

				break;
			case ActionType.SPLIT:
				this.pdf.startCut = -1;
				this.pdf.endCut = -1;
				const from = Math.min(cmd.startCut, cmd.endCut);
				const to = Math.max(cmd.startCut, cmd.endCut);

				this.pdf.pages.forEach((p, index) => {
					if (index >= from && index < to) {
						if (!p.groupIds.includes(cmd.groupId)) {
							p.groupIds.push(cmd.groupId)
						}
					}
				})
				if (!this.pdf.pendingCuts.includes(cmd.groupId)) {
					this.pdf.pendingCuts.push(cmd.groupId)
				}
				console.log(this.pdf.pages)
				break;
		}
	}


	setPdf(doc: PDFDocumentProxy) {
		this.pdf.proxy = doc;
		this.pdf.pages = Array.from(
			{ length: this.pdf.proxy.numPages },
			(_, i) => ({
				id: i,
				pageNum: i + 1,
				groupIds: [],
				scale: this.pdf.carouselScale,
			})
		);
	}
	setCurrent(page: number) {
		this.pdf.currentPage = page;
	}
	resetChanges() {
		this.execute({ type: ActionType.RESET })
	}

	delete(page: Page) {
		const cmd: Command = {
			type: ActionType.DELETE,
			page: page,
			previousIndex: this.pdf.pages.indexOf(page)
		}
		this.execute(cmd)
	}

	swap(animatingCards: Set<Page>, draggingCard: Page | null, card: Page, dragDuration: number) {
		if (
			draggingCard === null ||
			draggingCard === card ||
			animatingCards.has(card)
		)
			return;

		const cmd: Command = {
			type: ActionType.SWAP,
			pageA: draggingCard,
			pageB: card,
		}
		this.execute(cmd);

		animatingCards.add(card);
		setTimeout(() => animatingCards.delete(card), dragDuration);
	}

	split(cut: number) {
		this.setCut(cut)
		console.log(`cut -> ${cut}`)
		if (this.pdf.startCut === -1 || this.pdf.endCut === -1) return;

		if (this.pdf.startCut === this.pdf.endCut) {
			this.pdf.startCut = -1;
			this.pdf.endCut = -1;
			return;
		}
		const cmd: Command = {
			type: ActionType.SPLIT,
			startCut: this.pdf.startCut,
			endCut: this.pdf.endCut,
			groupId: generate(),
		}
		console.log(`start ${this.pdf.startCut}, end -> ${this.pdf.endCut}`);
		this.execute(cmd);
	}

	private setCut(cut: number) {
		if (this.pdf.startCut === -1) {
			this.pdf.startCut = cut;
		} else if (this.pdf.endCut === -1) {
			this.pdf.endCut = cut;
		}
	}

	export() {
		console.log("placeholder!!!! for now >:)")
	}
}


export const pdfManager = new PdfManager()

export type Pdf = {
	proxy?: PDFDocumentProxy;
	pages: Page[];
	pendingCuts: string[];
	carouselScale: number;
	viewerScale: number;
	currentPage: number;
	startCut: number;
	endCut: number;
}

export type Page = {
	id: number;
	pageNum: number;
	groupIds: string[];
	scale: number;
}

