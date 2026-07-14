<script lang="ts">
    import * as ButtonGroup from "$lib/components/ui/button-group/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { getCurrentTableSelected } from "$lib/index.js";
    import { SetHeader } from "../../../bindings/merger/services/setting";

    let { selectedSheet = $bindable("") } = $props();

    function settingHeader() {
        SetHeader(getCurrentTableSelected(selectedSheet));
    }

    function hiddenRows() {

    }

    function hiddenCols() {

    }
</script>

<ButtonGroup.Root class="text-gray-500">
    <ButtonGroup.Root>
        <Button onclick={settingHeader} size="xs" variant="outline" title="表头只能选择一行">设为表头</Button>
        <Button onclick={hiddenRows} size="xs" variant="outline">隐藏行</Button>
    </ButtonGroup.Root>
    <ButtonGroup.Root>
        <Button onclick={hiddenCols} size="xs" variant="outline">隐藏列</Button>
    </ButtonGroup.Root>
</ButtonGroup.Root>
