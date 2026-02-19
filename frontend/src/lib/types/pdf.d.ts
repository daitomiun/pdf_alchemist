
export type Pdf = {
	Pdf: PDFDocumentProxy;
	scale: number;
	pageNumbers: number[]
};

export enum ActionType {
	Split = "Split",
	Swap = "Swap",
	Delete = "Delete",
	Add = "Add",
}

export enum PdfDefinition {
	low = 0.1,
	medium = 0.5,
	high = 3,
}

export type PdfCommand<TPayload = unknown> = {
	type: ActionType;
	payload: TPayload;
	do(pdf: Pdf): Effect.promise<any>
	undo(pdf: Pdf): Effect.promise<any>
};

export class PdfHistory {
	private done: PdfCommand[] = [];
	private undone: PdfCommand[] = [];
	async execute(pdf: Pdf, cmd: PdfCommand) { }
	async undo(pdf: Pdf) { }
	async redo(pdf: Pdf) { }
}
