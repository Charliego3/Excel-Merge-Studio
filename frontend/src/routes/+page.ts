import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
    return {
        file: "选择工作簿、Sheet 和列",
        sheets: [],
    };
};
