<script lang="ts">
    import { Sheet } from "@lucide/svelte";
    import { getStateContext } from "$lib/state";
    import { CircleX } from "@lucide/svelte";
    import { goto, invalidateAll } from "$app/navigation";
    import type { WorkbookMeta } from "../../../bindings/merger/utility";
    import { RemoveWorkbook, WorkbooksMeta } from "../../../bindings/merger/services/workbook";

    let { book, index }: { book: WorkbookMeta | null; index: number } = $props();
    let state = getStateContext();
</script>

<button
    onclick={() => {
        state.work_index = index;
        state.main_index = -1;
        goto(`/preview?id=${book?.ID}`);
        localStorage.setItem("currentWorkbookId", book?.ID ?? "");
    }}
    class={`group flex border rounded-lg w-60 hover:cursor-pointer p-2 items-center gap-2 justify-between ${state.work_index == index ? "border-[#98cdbc] bg-green-50" : "border-gray-300"}`}
>
    <div class="flex gap-2 items-center">
        <div class="self-center bg-emerald-100 p-2 rounded-lg">
            <Sheet size={18} color="#127A65" />
        </div>
        <div class="flex flex-col items-start gap-1 text-left">
            <span
                title={book?.FilePath}
                class={`font-medium text-[10px] ${state.work_index != index ? "text-gray-500" : ""}`}
                >{book?.Name}</span
            >
            <span title={book?.SheetNames?.join("\n")} class="line-clamp-2 text-gray-500 text-[8px]">
                {book?.SheetNames?.join(" | ")}
            </span>
        </div>
    </div>
    <a
        href={"#"}
        type="button"
        onclick={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            RemoveWorkbook(book?.ID ?? "").then(async () => {
                WorkbooksMeta().then((books) => {
                    let currentId = localStorage.getItem("currentId");
                    if (!currentId) return;
                    let ids = books?.map(book => book.ID);
                    // console.dir({current: localStorage.getItem("currentId"), names: ids, id: book?.ID})
                    if (!ids?.includes(currentId)) {
                        goto("/", { invalidateAll: true });
                    } else {
                        invalidateAll();
                    }
                });
            });
        }}
        class="invisible group-hover:visible text-gray-400 hover:text-red-700 cursor-pointer"
    >
        <CircleX size={16} />
    </a>
</button>
