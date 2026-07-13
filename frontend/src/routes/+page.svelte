<script lang="ts">
    import { Plus } from "@lucide/svelte";
    import { ShowFilePicker } from "../../bindings/merger/services/reader";
    import { type Sheet } from "../../bindings/merger/utility";
    import * as Empty from "$lib/components/ui/empty/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Sheet as SheetIcon } from "@lucide/svelte";
    import WorkbookPreview from "$lib/components/WorkbookPreview.svelte";
    import * as Kbd from "$lib/components/ui/kbd/index.js";
    import type { PageProps } from './$types';
    import { onDestroy } from "svelte";
    import * as ButtonGroup from "$lib/components/ui/button-group/index.js";

    let { data }: PageProps = $props();
    let file: string = $derived(data.file);
    let sheets: Sheet[] = $derived(data.sheets);
    let headerHeight: number = $state(0);
    let toolbarHeight: number = $state(0);

    let loading: boolean = $state(false);

    onDestroy(() => {
        localStorage.removeItem("currentId");
    });

    function onWorkbookLoaded() {
        loading = true;
        ShowFilePicker()
            .then((data: any) => {
                file = data.file;
                sheets = data.sheets;
                localStorage.setItem("currentId", data.id);
            })
            .catch((e) => console.log(e))
            .finally(() => loading = false);
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
                    <span class="text-[11px] font-normal text-gray-500 ml-1">(选择表头，所需列，隐藏行/列)</span>
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
        <div class="flex justify-center items-center">
            Loading...
        </div>
    {:else if sheets.length > 0}
        <div class="p-2 text-gray-500" bind:clientHeight={toolbarHeight}>
            <ButtonGroup.Root>
                <ButtonGroup.Root>
                    <Button size="xs" variant="outline" title="表头只能选择一行">设为表头</Button>
                    <Button size="xs" variant="outline">隐藏行</Button>
                </ButtonGroup.Root>
                <ButtonGroup.Root>
                    <Button size="xs" variant="outline">隐藏列</Button>
                </ButtonGroup.Root>
            </ButtonGroup.Root>
        </div>
        <WorkbookPreview border checked {sheets} headerHeight={headerHeight + toolbarHeight} />
    {:else}
        <Empty.Root>
            <Empty.Header>
                <Empty.Media variant="icon">
                    <SheetIcon size={30} color="#127A65" />
                </Empty.Media>
                <Empty.Title>未选择任何文件</Empty.Title>
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
