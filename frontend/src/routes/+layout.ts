export const ssr = false;
export const prerender = true;

import type { PageLoad } from "./$types";
import {
    FileSpreadsheet,
    KeyRound,
    TableColumnsSplit,
    Fullscreen,
} from "@lucide/svelte";

export const load: PageLoad = async () => {
    return {
        steps: [
            {
                name: "文件",
                path: "/",
                description: "3 个工作簿",
                icon: FileSpreadsheet,
                iconColor: "oklch(82.8% 0.189 84.429)",
                iconBgColor: "oklch(96.2% 0.059 95.617)",
            },
            {
                name: "主键",
                path: "/upload",
                description: "客户编号",
                icon: KeyRound,
                iconColor: "oklch(70.4% 0.191 22.216)",
                iconBgColor: "oklch(93.6% 0.032 17.717)",
            },
            {
                name: "列映射",
                path: "/",
                description: "12 列已选",
                icon: TableColumnsSplit,
                iconColor: "oklch(74% 0.238 322.16)",
                iconBgColor: "oklch(95.2% 0.037 318.852)",
            },
            {
                name: "预览",
                path: "/",
                description: "待确认",
                icon: Fullscreen,
                iconColor: "oklch(67.3% 0.182 276.935)",
                iconBgColor: "oklch(93% 0.034 272.788)",
            },
        ],
    };
};
