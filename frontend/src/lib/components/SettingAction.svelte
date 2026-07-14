<script lang="ts">
    import * as ButtonGroup from "$lib/components/ui/button-group/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { getCurrentTableSelected } from "$lib/index.js";
    import { SetHeader, SetMain, GetMain } from "../../../bindings/merger/services/setting";
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

    function settingHeader() {
        SetHeader(getCurrentTableSelected(selectedSheet));
    }

    function hiddenRows() {

    }

    function hiddenCols() {

    }

    function setMain() {
        main = { Workbook: selectedWorkbook, Sheet: selectedSheet };
        SetMain(main);
    }
</script>

<ButtonGroup.Root class="text-gray-500">
    <ButtonGroup.Root>
        <Button title={description} onclick={setMain} size="xs" variant="outline">设为主表 (Sheet)</Button>
    </ButtonGroup.Root>
    <ButtonGroup.Root>
        <Button onclick={settingHeader} size="xs" variant="outline" title="表头只能选择一行">设为表头</Button>
        <Button onclick={hiddenRows} size="xs" variant="outline">隐藏行</Button>
    </ButtonGroup.Root>
    <ButtonGroup.Root>
        <Button onclick={hiddenCols} size="xs" variant="outline">隐藏列</Button>
    </ButtonGroup.Root>
</ButtonGroup.Root>
