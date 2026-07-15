<script lang="ts">
    import { Plus } from "@lucide/svelte";
    import { ShowFilePicker } from "../../bindings/merger/services/reader";
    import { type Sheet } from "../../bindings/merger/utility";
    import * as Empty from "$lib/components/ui/empty/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Sheet as SheetIcon } from "@lucide/svelte";
    import WorkbookPreview from "$lib/components/WorkbookPreview.svelte";
    import * as Kbd from "$lib/components/ui/kbd/index.js";
    import type { PageProps } from "./$types";
    import { onDestroy } from "svelte";
    import SettingAction from "$lib/components/SettingAction.svelte";
    import Loading from "$lib/components/Loading.svelte";
    import { Events } from "@wailsio/runtime";
    import { setCurrentWorkbookId, clearCurrentWorkbookId } from "$lib/index.js";
    import type { Main, Setting } from "../../bindings/merger/utility/models";

    let { data }: PageProps = $props();
    let file: string = $derived(data.file);
    let sheets: Sheet[] = $state((() => data.sheets)());
    let headerHeight: number = $state(0);
    let toolbarHeight: number = $state(0);
    let loading: boolean = $state(false);
    let selectedSheet: string = $state("");
    let workbookId: string = $derived(data.id);

    Events.On("workbook:read:start", () => loading = true);
    Events.On("workbook:sheet:setting", (e) => {
        const data: Setting = e.data;
        sheets = sheets.map((sheet) =>
            sheet.Name === data.Sheet ? { ...sheet, Header: data.Rows?.[0] ?? 0 } : sheet,
        );
    });

    Events.On("workbooks:sheet:removed", (e) => {
        const main: Main = e.data;
        const index = sheets.findIndex((sheet) => sheet.Name === main.Sheet);
        selectedSheet = index > 1 ? sheets[index - 1]?.Name : sheets[0]?.Name;
        sheets = sheets.filter((sheet) => sheet.Name !== main.Sheet);
    });

    onDestroy(() => clearCurrentWorkbookId());

    function onWorkbookLoaded() {
        ShowFilePicker().then((data: any) => {
            workbookId = data.id;
            file = data.file;
            sheets = data.sheets;
            selectedSheet = data.sheets[0]?.Name ?? "";
            setCurrentWorkbookId(data.id);
        })
        .catch((e) => console.log(e))
        .finally(() => (loading = false));
    }
</script>

<div class="w-full h-full flex flex-col">
    <div
        bind:clientHeight={headerHeight}
        class="flex-none flex justify-between items-center border-b border-b-gray-300 select-none w-full"
    >
        <div class="flex flex-col py-2 px-4 gap-1 w-full">
            <span class="font-bold">
                文件与工作簿
                {#if sheets.length > 0}
                    <span class="text-[11px] font-normal text-gray-500 ml-1"
                        >(选择表头，所需列，隐藏行/列)</span
                    >
                {/if}
            </span>
            <span class="text-[11px] text-gray-500">{file}</span>
        </div>
        {#if sheets.length > 0}
            <button
                title="选择其他文件"
                onclick={onWorkbookLoaded}
                class="p-2 cursor-pointer border border-gray-300 rounded-xl mr-4"
            >
                <Plus size={18} />
            </button>
        {/if}
    </div>
    {#if loading}
        <div class="h-full flex justify-center items-center">
            <Loading {headerHeight} />
        </div>
    {:else if sheets.length > 0}
        <div class="p-2" bind:clientHeight={toolbarHeight}>
            <SettingAction selectedWorkbook={workbookId} {selectedSheet} />
        </div>
        <WorkbookPreview
            tabBorder
            selectedWorkbook={workbookId}
            border
            checked
            {sheets}
            bind:selectedSheet
            headerHeight={headerHeight + toolbarHeight}
        />
    {:else}
        <Empty.Root class="select-none">
            <Empty.Header>
                <Empty.Media variant="icon">
                    <SheetIcon size={30} color="#127A65" />
                </Empty.Media>
                <Empty.Title>选择文件</Empty.Title>
                <Empty.Description class="text-[12px]/6">
                    请选择一个表格文件进行打开和查看。支持的文件格式：
                    <Kbd.Group>
                        <Kbd.Root>.xlsx</Kbd.Root>
                        <Kbd.Root>.xls</Kbd.Root>
                        和
                        <Kbd.Root>.csv</Kbd.Root>
                    </Kbd.Group>
                </Empty.Description>
            </Empty.Header>
            <Empty.Content>
                <Button onclick={onWorkbookLoaded}>选择文件</Button>
            </Empty.Content>
        </Empty.Root>
    {/if}
</div>
