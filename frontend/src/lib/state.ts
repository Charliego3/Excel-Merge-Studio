import { setContext, getContext } from "svelte";

export type State = {
    main_index: number;
    work_index: number;
};

let key = Symbol("state");
export const sidebarWidthKey = Symbol("sidebarWidth");

export function setStateContext(state: State): void {
    setContext(key, state);
}

export function getStateContext(): State {
    return getContext(key);
}
