<script lang="ts">
    import WorkbookPreview from "$lib/components/WorkbookPreview.svelte";
    import type { PageProps } from "./$types";
    import { GetWorkbook } from "../../../bindings/merger/services/workbook";
    import SettingAction from "$lib/components/SettingAction.svelte";
    import { Events } from "@wailsio/runtime";
    import type { Setting } from "../../../bindings/merger/utility/models";

    let { data }: PageProps = $props();
    let headerHeight = $state(0);
    let toolbarHeight = $state(0);
    let selectedWorkbook = $derived(data.metas?.[0]?.ID ?? "");
    let workbookInfo = $derived(await GetWorkbook(selectedWorkbook));
    let sheets = $derived(workbookInfo.Sheets);
    let selectedSheet: string = $derived((() => workbookInfo?.Sheets?.[0]?.Name ?? "")());

    Events.On("workbook:sheet:setting", (e) => {
        const data: Setting = e.data;
        sheets = sheets?.map((sheet) =>
            sheet && sheet.Name === data.Sheet
                ? { ...sheet, Header: data.Rows?.[0] ?? 0 }
                : sheet,
        ) ?? [];
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
        <div class="p-2 flex gap-2" bind:clientHeight={toolbarHeight}>
            <select bind:value={selectedWorkbook} class="native-select">
                {#if data.metas}
                    {#each data.metas as meta}
                        <option value={meta.ID}>{meta.Name}</option>
                    {/each}
                {:else}
                    <option value="">Select workbook</option>
                {/if}
            </select>
            <SettingAction {selectedWorkbook} {selectedSheet} />
        </div>

        <WorkbookPreview {selectedWorkbook} tabBorder bind:selectedSheet checked border sheets={sheets} headerHeight={headerHeight + toolbarHeight} />
    </div>
</div>
