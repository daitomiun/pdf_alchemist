import type { Page } from "$lib/state/state.svelte";

export enum ActionType {
	DELETE = "delete",
	SPLIT = "split",
	ADD = "add",
	SWAP = "swap",
	RESET = "reset",
}

export type Command =
	| { type: ActionType.DELETE; page: Page; previousIndex: number }

	| { type: ActionType.SWAP; pageIdA: number; pageIdB: number }
	| { type: ActionType.SPLIT; groupId: string; startCut: number, endCut: number }
	| { type: ActionType.ADD; sourceDoc: string; sourcePage: number; insertAt: number }
	| { type: ActionType.RESET };

