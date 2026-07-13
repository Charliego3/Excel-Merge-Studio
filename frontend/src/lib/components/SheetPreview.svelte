<script lang="ts">
    import { index2column } from "$lib/index";
    import type { Sheet } from "../../../bindings/merger/utility";
    import * as Empty from "$lib/components/ui/empty/index.js";
    import { Grid2x2X } from "@lucide/svelte";
    import { Checkbox } from "$lib/components/ui/checkbox/index.js";
    import { Label } from "./ui/label";

    let {
        classes = "overflow-auto text-[10px]",
        style = "",
        border = true,
        sheet = $bindable<Sheet | undefined | null>(null),
        checked = false,
        onRowSelected,
        onColSelected,
    }: {
        classes?: string;
        style?: string;
        border?: boolean;
        sheet: Sheet | undefined | null;
        checked?: boolean;
        onRowSelected?: (row: number) => void;
        onColSelected?: (col: number) => void;
    } = $props();

    function handleColCheckboxChange() {}

    function handleRowCheckboxChange() {}
</script>

{#if sheet?.Columns ?? 0 > 0}
    <div
        class={`${classes} ${border ? "border-t border-t-[#EBEBEB]" : ""}`}
        {style}
    >
        <table class="table-fixed">
            <thead class="sticky top-0 h-5">
                <tr>
                    <th class="relative">
                        {#if checked}
                            <Checkbox
                                onCheckedChange={handleColCheckboxChange}
                                class="w-3 h-3 absolute right-1 top-1"
                            />
                            <Checkbox
                                onCheckedChange={handleRowCheckboxChange}
                                class="w-3 h-3 absolute left-1 bottom-1"
                            />
                        {/if}
                    </th>
                    {#each Array(sheet?.Columns ?? 0) as _, index}
                        <th>
                            <div class="flex gap-2 items-center justify-center">
                                {#if checked}
                                    <Checkbox
                                        id={`col-${index}`}
                                        onCheckedChange={onColSelected
                                            ? () => onColSelected(index)
                                            : undefined}
                                    />
                                {/if}
                                <Label
                                    for={`col-${index}`}
                                    class="font-bold text-[12px]"
                                    >{index2column(index)}: {index + 1}</Label
                                >
                            </div>
                        </th>
                    {/each}
                </tr>
            </thead>
            <tbody>
                {#each sheet?.Data as row, rowIndex}
                    <tr>
                        <td>
                            <div class="flex gap-2 items-center justify-center">
                                {#if checked}
                                    <Checkbox
                                        id={`row-${rowIndex}`}
                                        onCheckedChange={onRowSelected
                                            ? () => onRowSelected(rowIndex)
                                            : undefined}
                                    />
                                {/if}
                                <Label
                                    class="text-[11px] font-bold"
                                    for={`row-${rowIndex}`}
                                    >{rowIndex + 1}</Label
                                >
                            </div>
                        </td>
                        {#each row?.Data as cell, colIndex}
                            {#if !cell.Skip}
                                <td
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
