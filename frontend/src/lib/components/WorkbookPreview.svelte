<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index.js";
    import ScrollArea from "$lib/components/ui/scroll-area/scroll-area.svelte";
    import type { Sheet } from "../../../bindings/merger/utility";
    import SheetPreview from '$lib/components/SheetPreview.svelte';
    import { X } from "@lucide/svelte";
    import { RemoveSheet } from "../../../bindings/merger/services/workbook";

    let {
        tabBorder = false,
        border = false,
        sheets = $bindable<(Sheet | null)[] | null>(null),
        headerHeight,
        checked = false,
        selectedSheet = $bindable<string>(undefined),
        selectedWorkbook = $bindable<string>(undefined),
    }: {
        tabBorder?: boolean;
        border?: boolean;
        sheets: (Sheet | null)[] | null,
        headerHeight: number;
        checked?: boolean
        selectedSheet: string;
        selectedWorkbook: string;
    } = $props();

    function removeSheet(sheetName: string) {
        selectedSheet = sheets?.[0]?.Name ?? "";
        RemoveSheet({ Workbook: selectedWorkbook, Sheet: sheetName });
    }
</script>

<div class="flex-1 h-full overflow-hidden">
    <Tabs.Root class={`h-full rounded-none ${tabBorder && (sheets && sheets?.length > 1) ? 'border-t' : ''}`} bind:value={selectedSheet}>
        {#if sheets && sheets?.length > 1}
            <ScrollArea orientation="horizontal" class="-mb-1 py-0.5" scrollbarXClasses="hidden">
                <Tabs.List class="bg-white pl-2 pr-3">
                    {#each sheets as sheet}
                        {#if sheet}
                            <Tabs.Trigger class="group text-[10px] data-active:bg-blue-200/60" value={sheet.Name}>
                                <span>{sheet.Name}</span>
                                {#if checked}
                                    <button onclick={(e) => {
                                        e.stopPropagation();
                                        e.preventDefault();
                                        removeSheet(sheet.Name);
                                    }} onpointerdown={(e) => {
                                        e.stopPropagation();
                                        e.preventDefault();
                                    }} class="cursor-pointer">
                                        <X color="gray" class={`w-3 h-3 ${sheet.Name === selectedSheet ? 'inline' : 'hidden'} group-hover:text-red-700 group-hover:cursor-pointer group-hover:inline`} />
                                    </button>
                                {/if}
                            </Tabs.Trigger>
                        {/if}
                    {/each}
                </Tabs.List>
            </ScrollArea>
        {/if}
        {#each sheets as sheet}
            {#if sheets && sheet}
                <Tabs.Content class={`h-full ${sheets?.length > 1 ? "-mt-1" : ""}`} value={sheet.Name}>
                    <SheetPreview
                        {checked}
                        border={border || sheets?.length > 1}
                        style={`height: calc(100vh ${sheets?.length > 1 ? "- 40px" : "+ 1px"} - ${headerHeight}px)`}
                        sheet={sheet} />
                </Tabs.Content>
            {/if}
        {/each}
    </Tabs.Root>
</div>
