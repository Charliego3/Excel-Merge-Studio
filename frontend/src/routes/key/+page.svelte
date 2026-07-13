<script lang="ts">
    import WorkbookPreview from "$lib/components/WorkbookPreview.svelte";
    import * as NativeSelect from "$lib/components/ui/native-select/index.js";
    import type { PageProps } from "./$types";
    import { GetWorkbook } from "../../../bindings/merger/services/workbook";

    let { data }: PageProps = $props();
    let headerHeight = $state(0);
    let toolbarHeight = $state(0);
    let selectedWorkbook = $derived.by(() => {
        if (data.metas) return data.metas[0]?.ID;
        return "";
    });
    let workbookInfo = $derived.by(async () => {
        if (!selectedWorkbook) return;
        return await GetWorkbook(selectedWorkbook);
    });
</script>

<div class="flex flex-col w-full h-full">
    <div
        bind:clientHeight={headerHeight}
        class="flex justify-between items-center border-b border-b-gray-300 select-none w-full"
    >
        <div class="flex flex-col py-2 px-4 gap-1 w-full">
            <span class="font-bold">主键列 & 表头</span>
            <span class="text-[11px] text-gray-500">设置列 (表头、隐藏、选择)</span>
        </div>
    </div>

    <div class="w-full h-full">
        <div class="p-2" bind:clientHeight={toolbarHeight}>
            <NativeSelect.Root size="sm" value={selectedWorkbook}>
                {#if data.metas}
                    {#each data.metas as meta}
                        <NativeSelect.Option value={meta.ID}>{meta.Name}</NativeSelect.Option>
                    {/each}
                {:else}
                    <NativeSelect.Option value="">Select workbook</NativeSelect.Option>
                {/if}
            </NativeSelect.Root>
        </div>

        <WorkbookPreview tabBorder checked border sheets={(await workbookInfo)?.Sheets ?? []} headerHeight={headerHeight + toolbarHeight} />
    </div>
</div>
