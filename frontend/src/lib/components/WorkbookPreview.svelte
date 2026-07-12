<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index.js";
    import ScrollArea from "$lib/components/ui/scroll-area/scroll-area.svelte";
    import type { Sheet } from "../../../bindings/merger/utility";
    import SheetPreview from '$lib/components/SheetPreview.svelte';

    let { sheets, headerHeight }: { sheets: Sheet[], headerHeight: number } = $props();
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
            <Tabs.Content class="h-full -mt-1" value={sheet.ID}>
                <SheetPreview
                    border={sheets.length > 1}
                    style={`height: calc(100vh ${sheets.length > 1 ? "- 44px" : "+ 4px"} - ${headerHeight}px)`}
                    sheet={sheet} />
            </Tabs.Content>
        {/each}
    </Tabs.Root>
</div>
