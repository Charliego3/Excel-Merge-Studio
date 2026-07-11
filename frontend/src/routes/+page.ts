import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
    return {
        current: -1,
        steps: [
            { name: "文件", description: "3 个工作簿" },
            { name: "主键", description: "客户编号" },
            { name: "列映射", description: "12 列已选" },
            { name: "预览", description: "待确认" },
        ],
    };
};
