
export type Pdf = {
	Pdf: PDFDocumentProxy;
	numPages: number;
};

export enum ActionType {
	Split = "Split",
	Swap = "Swap",
	Delete = "Delete",
	Add = "Add",
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
