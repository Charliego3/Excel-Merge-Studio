<script lang="ts">
    import { Window } from "@wailsio/runtime";
    import { sidebarWidthKey } from "$lib/state";
    import { getContext, onMount } from "svelte";

    const cellWidth = 100;
    const cellHeight = 40;

    let { headerHeight }: { headerHeight: number } = $props();
    let windowWidth: number = (await Window.Size()).width;
    let windowHeight: number = $state((await Window.Size()).height);
    let sidebarWidth: () => number = getContext(sidebarWidthKey);
    let width: number = $derived(windowWidth - sidebarWidth());
    let height: number = $derived(windowHeight - headerHeight);
    let columns: number = $derived(Math.floor(width / cellWidth));
    let rows: number = $derived(Math.ceil(height / cellHeight));

    onMount(() => {
        window.addEventListener("resize", async () => {
            windowHeight = (await Window.Size()).height;
        });
    });
</script>

<table class="animate-pulse pointer-events-none border-none">
    <tbody>
        {#each Array(rows) as _}
            <tr style={`height: ${cellHeight}px;`}>
                {#each Array(columns) as _}
                    <td style={`width: ${cellWidth}px;`}></td>
                {/each}
            </tr>
        {/each}
    </tbody>
</table>

<style>
    tbody tr:nth-child(odd) {
        background-color: transparent;
    }

    tbody tr:first-child td {
        border-top: none;
        background-color: ghostwhite;
    }
</style>
