<script lang="ts">
    import "../layout.css";
    import WorkBook from "$lib/components/WorkBook.svelte";
    import FlowStep from "$lib/components/FlowStep.svelte";
    import { getStateContext, setStateContext, sidebarWidthKey, type State } from "$lib/state";
    import { Events } from "@wailsio/runtime";
    import { setContext } from "svelte";
    import { WorkbooksMeta } from "../../bindings/merger/services/workbook";
    import type { WorkbookMeta } from "../../bindings/merger/utility";
    import { ModeWatcher } from "mode-watcher";

    let isMac = navigator.userAgent.includes("Mac");

    let { data, children }: { data: {
        steps: FlowStep[];
        workbooks: (WorkbookMeta | null)[] | null
        fullscreen: boolean;
    }; children: any } = $props();
    let steps: any = $derived(data.steps);
    let workbooks: (WorkbookMeta | null)[] | null = $derived(data.workbooks);
    let fullscreen: boolean = $derived(data.fullscreen);

    let sidebarWidth = $state(0);
    setContext(sidebarWidthKey, () => sidebarWidth);

    Events.On("common:WindowFullscreen", () => (fullscreen = true));
    Events.On("common:WindowUnFullscreen", () => (fullscreen = false));
    Events.On("workbooks:updated", async () => workbooks = await WorkbooksMeta());

    let appState: State = $state({ main_index: 0, work_index: -1 });
    let context = getStateContext();
    if (!context) {
        setStateContext(appState);
    }
</script>

<ModeWatcher />
<div class="flex w-full h-full">
    <div bind:clientWidth={sidebarWidth}
        class={`border-r select-none flex-none border-r-gray-300 w-64 flex flex-col bg-[#FBFBFA]`}
    >
        {#if isMac && !fullscreen}
            <div class="flex-none h-10"></div>
        {/if}

        <div class="flex flex-col p-2 gap-2">
            {#each steps as step, index}
                <FlowStep {step} {index} />
            {/each}
        </div>

        <div class="text-gray-500 flex items-center gap-1 mt-2 mx-4">
            <span class="text-sm">工作簿</span>
            {#if workbooks && workbooks.length > 0}
                <span class="text-[11px]">({workbooks.length})</span>
            {/if}
        </div>
        {#if workbooks === null || workbooks?.length === 0}
            <div
                class={`h-full border border-gray-300 m-2 flex items-center justify-center ${isMac ? "rounded-2xl" : "rounded-lg"}`}
            >
                <span class="text-gray-400 text-[11px]">无工作薄</span>
            </div>
        {:else}
            <div class="overflow-y-auto overflow-x-hidden px-2 pb-2 flex flex-col gap-2 mt-1">
                {#each workbooks as book, index}
                    <WorkBook {book} {index} />
                {/each}
            </div>
        {/if}
    </div>

    <div class="w-full h-full bg-white" style={`width: calc(100vw - ${sidebarWidth}px);`}>
        {@render children()}
    </div>
</div>
