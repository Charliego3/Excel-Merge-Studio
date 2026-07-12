<script lang="ts">
    import { index2column } from "$lib/index";
    import type { Sheet } from "../../../bindings/merger/utility";
    import * as Empty from "$lib/components/ui/empty/index.js";
    import { Grid2x2X } from "@lucide/svelte";

    let {
        classes = "overflow-auto text-[10px]",
        style = "",
        border = true,
        sheet,
    }: {
        classes?: string;
        style?: string;
        border?: boolean;
        sheet: Sheet | undefined | null;
    } = $props();
</script>

{#if sheet?.Columns ?? 0 > 0}
    <div class={`${classes} ${border ? "border-t border-t-[#EBEBEB]" : ""}`} {style}>
        <table class="table-fixed">
            <thead class="sticky top-0 h-5">
                <tr>
                    <th></th>
                    {#each Array(sheet?.Columns ?? 0) as _, index}
                        <th>{index2column(index)}</th>
                    {/each}
                </tr>
            </thead>
            <tbody>
                {#each sheet?.Data as row, index}
                    <tr>
                        <td>{index + 1}</td>
                        {#each row as cell}
                            {#if !cell.Skip}
                                <td
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
