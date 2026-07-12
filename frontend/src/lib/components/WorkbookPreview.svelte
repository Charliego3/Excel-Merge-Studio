<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index.js";
    import ScrollArea from "$lib/components/ui/scroll-area/scroll-area.svelte";
    import type { Sheet } from "../../../bindings/merger/utility";
    import SheetPreview from '$lib/components/SheetPreview.svelte';

    let {
        sheets,
        headerHeight,
        checked = false,
        onRowSelected,
        onColSelected,
    }: {
        sheets: Sheet[],
        headerHeight: number;
        checked?: boolean
        onRowSelected?: (row: number) => void;
        onColSelected?: (col: number) => void;
    } = $props();
</script>

<div class="flex-1 h-full overflow-hidden">
    <Tabs.Root class="h-full" value={sheets[0]?.ID}>
        {#if sheets.length > 1}
            <ScrollArea orientation="horizontal" class="pt-1 px-1" scrollbarXClasses="hidden">
                <Tabs.List>
                    {#each sheets as sheet}
                        <Tabs.Trigger class="text-[10px]" value={sheet.ID}>
                            {sheet.Name}
                        </Tabs.Trigger>
                    {/each}
                </Tabs.List>
            </ScrollArea>
        {/if}
        {#each sheets as sheet}
            <Tabs.Content class={`h-full ${sheets.length > 1 ? "-mt-1" : ""}`} value={sheet.ID}>
                <SheetPreview
                    {checked}
                    {onColSelected}
                    {onRowSelected}
                    border={sheets.length > 1}
                    style={`height: calc(100vh ${sheets.length > 1 ? "- 44px" : "+ 4px"} - ${headerHeight}px)`}
                    sheet={sheet} />
            </Tabs.Content>
        {/each}
    </Tabs.Root>
</div>
