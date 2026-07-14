<script lang="ts">
    import { index2column } from "$lib/index";
    import type { Sheet } from "../../../bindings/merger/utility";
    import * as Empty from "$lib/components/ui/empty/index.js";
    import { Grid2x2X } from "@lucide/svelte";
    import { Label } from "./ui/label";

    let {
        classes = "overflow-auto text-[10px]",
        style = "",
        border = true,
        sheet = $bindable<Sheet | undefined | null>(null),
        checked = false,
    }: {
        classes?: string;
        style?: string;
        border?: boolean;
        sheet: Sheet | undefined | null;
        checked?: boolean;
    } = $props();

    let table: HTMLTableElement | null = $state(null);

    function handleColCheckboxChange() {
        const row = table?.tHead?.rows[0];
        const switcher: HTMLInputElement | null | undefined =
            row?.querySelector('#colSwitch');

        for (let i = 1; i < (row?.cells?.length ?? 0); i++) {
            const checkbox: HTMLInputElement | null | undefined =
                row?.cells[i].querySelector('input[type="checkbox"]');
            if (checkbox) checkbox.checked = switcher?.checked ? true : false;
        }
    }

    function handleRowCheckboxChange() {
        const switcher: HTMLInputElement | null | undefined =
            table?.tHead?.querySelector('#rowSwitch');

        for (const row of table?.tBodies[0]?.rows ?? []) {
            const td = row.cells[0];
            let checkbox: HTMLInputElement | null = td.querySelector('input[type="checkbox"]');
            if (checkbox) checkbox.checked = switcher?.checked ? true : false;
        }
    }
</script>

{#if sheet?.Columns ?? 0 > 0}
    <div
        class={`${classes} ${border ? "border-t" : ""}`}
        {style}
    >
        <table data-sheet={sheet?.Name} data-workbook={sheet?.WorkbookId} class="table-fixed" bind:this={table}>
            <thead class="sticky top-0 h-5">
                <tr>
                    <th class="relative">
                        {#if checked}
                            <input id="colSwitch" onchange={handleColCheckboxChange} type="checkbox" class="rounded w-3 h-3 absolute right-1 top-1 border-input -mt-0.5"/>
                            <input id="rowSwitch" onchange={handleRowCheckboxChange} type="checkbox" class="rounded w-3 h-3 absolute left-1 bottom-1 border-input -mt-0.5"/>
                        {/if}
                    </th>
                    {#each Array(sheet?.Columns ?? 0) as _, index}
                        <th>
                            <div class="flex gap-2 items-center justify-center">
                                {#if checked}
                                    <input type="checkbox" class="thead-col rounded border-input -mt-0.5"/>
                                {/if}
                                <Label for={`col-${index}`} class="font-bold text-[12px]">{index2column(index)}: {index + 1}</Label>
                            </div>
                        </th>
                    {/each}
                </tr>
            </thead>
            <tbody>
                {#each sheet?.Data as row, rowIndex}
                    <tr data-header={sheet?.Header === rowIndex}>
                        <td data-header={sheet?.Header === rowIndex}>
                            <div class="flex gap-2 items-center justify-center">
                                {#if checked}
                                    <input id={`row-${rowIndex}`} type="checkbox" class="rounded border-input -mt-0.5"/>
                                {/if}
                                <Label class="text-[11px] font-bold" for={`row-${rowIndex}`}>{rowIndex + 1}</Label>
                            </div>
                        </td>
                        {#each row?.Data as cell}
                            {#if !cell.Skip}
                                <td data-header={sheet?.Header === rowIndex}
                                    title={`${cell.IsMerged ? index2column(cell.StartCol - 1) + (cell.StartRow) + ":" + index2column(cell.EndCol - 1) + (cell.EndRow) : index2column(cell.StartCol - 1) + cell.StartRow}`}
                                    colspan={cell.IsMerged ? cell.ColSpan : 1}
                                    rowspan={cell.IsMerged ? cell.RowSpan : 1}
                                >
                                    {#if cell.Value !== ""}
                                        <div
                                            class={`${/^[+-]?(?:\d+\.?\d*|\.\d+)%?$/.test(cell.Value) ? "flex justify-end" : cell.IsMerged ? "flex items-center justify-center" : ""}`}
                                        >
                                            <span>{cell.Value}</span>
                                        </div>
                                    {:else}
                                        <div
                                            class="flex justify-end text-gray-300"
                                        >
                                            --
                                        </div>
                                    {/if}
                                </td>
                            {/if}
                        {/each}
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
{:else}
    <div class="flex h-full items-center justify-center">
        <Empty.Root>
            <Empty.Header>
                <Empty.Media variant="icon">
                    <Grid2x2X size={30} color="#127A65" />
                </Empty.Media>
                <Empty.Title>Sheet 暂无数据</Empty.Title>
                <Empty.Description class="text-[12px]/6">
                    当前没有可显示的数据，请检查文件内容或选择其他表格文件。
                </Empty.Description>
            </Empty.Header>
        </Empty.Root>
    </div>
{/if}
