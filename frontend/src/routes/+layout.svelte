<script lang="ts">
    import "../layout.css";
    import type { Flow } from "$lib/types";
    import WorkBook from "$lib/componments/WorkBook.svelte";
    import FlowStep from "$lib/componments/FlowStep.svelte";
    import { setStateContext, type State } from "$lib/state";
    import { Window, Events } from "@wailsio/runtime";
    import { onMount } from "svelte";
    import { Workbooks } from "../../bindings/merger/services/workbook";
    import type { Workbook } from "../../bindings/merger/utility";

    let fullscreen: boolean = $state(false);
    Events.On("common:WindowFullscreen", () => (fullscreen = true));
    Events.On("common:WindowUnFullscreen", () => (fullscreen = false));
    let appState: State = $state({ main_index: 0, work_index: -1 });
    setStateContext(appState);
    let isMac = navigator.userAgent.includes("Mac");
    let { data, children }: { data: Flow; children: any } = $props();
    let flow = $state((() => data)());
    let workbooks: (Workbook | null)[] | null = $state([]);

    onMount(async () => {
        fullscreen = await Window.IsFullscreen();
        workbooks = await Workbooks();
    });
</script>

<div class="flex w-full h-full">
    <div
        class={`border-r select-none flex-none border-r-gray-300 min-w-60 flex flex-col bg-[#FBFBFA]`}
    >
        {#if isMac && !fullscreen}
            <div class="flex-none h-10"></div>
        {/if}

        <div class="flex flex-col p-2 gap-2">
            {#each flow.steps as step: FlowStep, index}
                <FlowStep {step} {index} />
            {/each}
        </div>

        <div class="text-gray-500 flex items-center gap-1 mt-2 mx-4">
            <span class="text-sm">工作簿</span>
            {#if workbooks && workbooks.length > 0}
                <span class="text-[11px]">({workbooks.length})</span>
            {/if}
        </div>
        {#if workbooks && workbooks.length === 0}
            <div
                class={`h-full border border-gray-300 m-2 flex items-center justify-center ${isMac ? "rounded-2xl" : "rounded-lg"}`}
            >
                <span>请选择文件d223</span>
            </div>
        {:else}
            <div class="overflow-y-auto px-2 pb-2 flex flex-col gap-2 mt-1">
                {#each workbooks as book, index}
                    <WorkBook {book} {index} />
                {/each}
            </div>
        {/if}
    </div>

    <div class="w-full h-full bg-white">
        {@render children()}
    </div>
</div>
