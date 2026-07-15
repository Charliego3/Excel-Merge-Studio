<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { getCurrentTableSelected } from "$lib/index.js";
    import { SetHeader, SetMain, GetMain, DeleteColsAndRows } from "../../../bindings/merger/services/setting";
    import { WorkbooksMeta } from "../../../bindings/merger/services/workbook";

    let { selectedWorkbook, selectedSheet } = $props();
    let main = $state(await GetMain());
    let workbooks = $state(await WorkbooksMeta());
    let currentWorkbook = $derived(workbooks?.find(w => w.ID === main.Workbook));
    let description = $derived.by(() => {
        let desc = "工作薄: " + (currentWorkbook?.Name ?? "");
        desc += "\n工作表: " + (main.Sheet ?? "");
        return desc;
    });

    function setHeader() {
        SetHeader(getCurrentTableSelected(selectedSheet));
    }

    function deleteColsAndRows() {
        DeleteColsAndRows(getCurrentTableSelected(selectedSheet))
    }

    function setMain() {
        SetMain({ Workbook: selectedWorkbook, Sheet: selectedSheet });
    }
</script>

<div class="text-gray-500 flex justify-center gap-2">
    <Button title={description} onclick={setMain} size="xs" variant="outline">设为主表 (Sheet)</Button>
    <Button onclick={setHeader} size="xs" variant="outline" title="表头只能选择一行">设为表头</Button>
    <Button onclick={setHeader} size="xs" variant="outline" title="主键列只能选择一列">设为主键列</Button>
    <Button onclick={deleteColsAndRows} size="xs" variant="outline">删除行和列</Button>
</div>
