<script lang="ts">
    import { Plus } from "@lucide/svelte";
    import { ShowFilePicker } from "../../bindings/merger/services/reader";
    import { type Sheet, type Workbook } from "../../bindings/merger/utility";
    import { index2column } from "$lib/index";

    let file: string = $state("");
    let sheets: Sheet[] = $state([]);

    function onWorkbookLoaded() {
        ShowFilePicker().then((data: any) => {
            file = data.file;
            sheets = data.sheets;
        }).catch(e => console.log(e));
    }
</script>

<div class="w-full h-full flex flex-col">
    <div
        class="flex-none flex justify-between items-center border-b border-b-gray-300 select-none w-full"
    >
        <div class="flex flex-col py-2 px-4 gap-1 w-full">
            <span class="font-bold">文件与主表</span>
            <span class="text-[11px] text-gray-500">选择工作簿、Sheet 和列</span
            >
        </div>
        <button
            onclick={onWorkbookLoaded}
            class="p-2 cursor-pointer border border-gray-300 rounded-xl mr-4"
        >
            <Plus size={18} />
        </button>
    </div>
    {#each sheets as sheet}
        <div class="flex-1 p-2 w-full h-full overflow-hidden">
            <div
                class="table-container overflow-auto w-full h-full rounded text-[10px]"
            >
                <table class="table-fixed w-full">
                    <caption
                        class="sticky top-0 h-5 bg-black text-white"
                        >{file} -- {sheet?.Name}</caption
                    >
                    <thead class="sticky top-5 h-5 bg-amber-800">
                        <tr class="text-white">
                            {#each Array(sheet?.Columns ?? 0) as _, index}
                                <th>
                                    <div
                                        class="flex items-center justify-center"
                                    >
                                        {index2column(index)}
                                    </div>
                                </th>
                            {/each}
                        </tr>
                    </thead>
                    <tbody>
                        {#each sheet?.Data as row}
                            <tr>
                                {#each row as cell}
                                    {#if !cell.Skip}
                                        <td
                                            colspan={cell.IsMerged
                                                ? cell.ColSpan
                                                : 1}
                                            rowspan={cell.IsMerged
                                                ? cell.RowSpan
                                                : 1}
                                        >
                                            <div
                                                class={`${cell.IsMerged ? "flex items-center justify-center" : ""}`}
                                            >
                                                <span>{cell.Value}</span
                                                >
                                            </div>
                                        </td>
                                    {/if}
                                {/each}
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    {/each}
</div>
