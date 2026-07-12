<script lang="ts">
    import type { FlowStep } from "$lib/types";
    import { getStateContext } from "$lib/state";

    let { step, index }: { step: FlowStep; index: number } = $props();
    let state = getStateContext();
</script>

<a
    class={`h-14 rounded-xl p-2 flex justify-between border items-center ${state.main_index == index ? "border-[#98cdbc] bg-green-50" : "border-gray-300"}`}
    href={step.path}
    onclick={() => {
        state.main_index = index;
        state.work_index = -1;
    }}
>
    <div class="flex gap-3 items-center">
        <div
            class="p-2 rounded-full"
            style={`background-color: ${step.iconBgColor}`}
        >
            <step.icon size={16} color={step.iconColor} />
        </div>
        <div class="flex flex-col gap-1">
            <span
                class={`font-semibold text-sm ${state.main_index != index ? "text-gray-500" : ""}`}
                >{step.name}</span
            >
            <span class="text-gray-400 text-[10px]">{step.description}</span>
        </div>
    </div>
    {#if state.main_index == index}
        <div
            class="rounded-full bg-[#F4E6DD] px-1.5 py-0.5 text-[10px] text-[#B85F35] font-semibold"
        >
            当前
        </div>
    {/if}
</a>
