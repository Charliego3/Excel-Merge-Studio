<script lang="ts">
    import * as ButtonGroup from "$lib/components/ui/button-group/index.js";
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

<ButtonGroup.Root class="text-gray-500">
    <ButtonGroup.Root>
        <Button title={description} onclick={setMain} size="xs" variant="outline">设为主表 (Sheet)</Button>
    </ButtonGroup.Root>
    <ButtonGroup.Root>
        <Button onclick={setHeader} size="xs" variant="outline" title="表头只能选择一行">设为表头</Button>
    </ButtonGroup.Root>
    <ButtonGroup.Root>
        <Button onclick={deleteColsAndRows} size="xs" variant="outline">删除行和列</Button>
    </ButtonGroup.Root>
</ButtonGroup.Root>
