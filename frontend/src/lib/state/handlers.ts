import { ActionType, type Command } from "$lib/types/pdf";

type CommandHandler = (pages: number[], cmd: Command) => number[];

export const commandHandlers = new Map<ActionType, CommandHandler>([
	[ActionType.DELETE, (pages, cmd) => {

		pages.filter((page) => page == cmd.pageId)



		return []
	}],
	[ActionType.SWAP, (pages, cmd) => {
		return []
	}],
]);

