<script lang="ts">
    import type { WorkbookInfo } from '../../../bindings/merger/utility/models';
    import WorkbookPreview from '$lib/components/WorkbookPreview.svelte';

    let { data }: {data: WorkbookInfo} = $props();
    let headerHeight = $state(0);
</script>

<div class="flex flex-col w-full h-full">
    <div bind:clientHeight={headerHeight} class="flex justify-between items-center border-b border-b-gray-300 select-none w-full">
        <div class="flex flex-col py-2 px-4 gap-1 w-full">
            <span class="font-bold">工作簿预览</span>
            <span class="text-[11px] text-gray-500">{data?.FilePath}</span>
        </div>
    </div>

    {#if data?.Sheets}
        <WorkbookPreview sheets={data?.Sheets} headerHeight={headerHeight} />
    {/if}
</div>
