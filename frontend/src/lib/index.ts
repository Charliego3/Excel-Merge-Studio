// place files you want to import through the `$lib` alias in this folder.
export function index2column(index: number): string {
    const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    let result = "";
    while (index >= 0) {
        result = letters[index % 26] + result;
        index = Math.floor(index / 26) - 1;
    }
    return result;
}

export interface TableSelected {
    sheet?: string | null;
    rows: number[];
    cols: number[];
}

/// 获取当前表格选中的行和列
export function getCurrentTableSelected(sheet: string): TableSelected {
    let table: HTMLTableElement | null = document.querySelector(`table[data-sheet="${sheet}"]`);
    if (!table) return { sheet: undefined, rows: [], cols: [] };
    return {
        sheet: table.getAttribute('data-sheet'),
        rows: getSelected(table, 'tbody tr td input'),
        cols: getSelected(table, 'thead tr th input.thead-col'),
    };
}

function getSelected(table: HTMLTableElement, selectors: string): number[] {
    return [...table.querySelectorAll(selectors)]
        .map((el, index) => {
            if (el instanceof HTMLInputElement && el.checked) {
                return index;
            }
            return null;
        })
        .filter((index): index is number => index !== null);
}
