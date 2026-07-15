import type { Setting } from "../../bindings/merger/utility";

export function removeSheet(sheets: Sheet[], sheetName: string, selectedSheet: string): {
    sheets: Sheet[];
    sheetName: string;
} {
    const index = sheets?.findIndex((sheet) => sheet?.Name === sheetName) ?? 0;
    sheets = sheets?.filter((sheet) => sheet?.Name !== sheetName) ?? [];
    if (sheetName === selectedSheet) {
        selectedSheet = (index > 1 ? sheets?.[index - 1]?.Name : sheets?.[0]?.Name) ?? "";
    }
    return { sheets, sheetName: selectedSheet };
}

export function index2column(index: number): string {
    const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    let result = "";
    while (index >= 0) {
        result = letters[index % 26] + result;
        index = Math.floor(index / 26) - 1;
    }
    return result;
}

/// 获取当前表格选中的行和列
export function getCurrentTableSelected(sheet: string): Setting {
    let table: HTMLTableElement | null = document.querySelector(`table[data-sheet="${sheet}"]`);
    if (!table) return { Workbook: "", Sheet: "", Rows: [], Cols: [] };
    return {
        Workbook: table.getAttribute('data-workbook') ?? "",
        Sheet: table.getAttribute('data-sheet') ?? "",
        Rows: getSelected(table, 'tbody tr td input'),
        Cols: getSelected(table, 'thead tr th input.thead-col'),
    };
}

export function unselectAll(sheet: string): void {
    let table: HTMLTableElement | null = document.querySelector(`table[data-sheet="${sheet}"]`);
    if (!table) return;
    table.querySelectorAll('thead tr th input').forEach((el) => {
        if (el instanceof HTMLInputElement) el.checked = false;
    });
    table.querySelectorAll('tbody tr td input').forEach((el) => {
        if (el instanceof HTMLInputElement) el.checked = false;
    });
}

export function deleteColsAndRows(sheets: Sheet[], setting: Setting): Sheet[] {
    return sheets.map((sheet) => {
        if (sheet.Name !== setting.Sheet || !sheet.Data) return sheet;
        let maxCellCount = 0;
        let rows = sheet.Data.filter((_, idx) => !setting.Rows?.includes(idx))
            .map((row) => {
                let cols = row.Data.filter((_, idx) => !setting.Cols?.includes(idx));
                if (cols.length > maxCellCount) maxCellCount = cols.length;
                return { ...row, Data: cols, Columns: cols.length };
            })
            .filter((row) => row.Columns > 0);
        return { ...sheet, Data: rows, Columns: maxCellCount };
    });
}

export function deleteCols(sheets: Sheet[], setting: Setting): Sheet[] {
    return sheets.map((sheet) => {
        if (sheet.Name !== setting.Sheet || !sheet.Data) return sheet;
        let maxCellCount = 0;
        let rows: Row[] = [];
        for (let i = 0; i < sheet.Data.length; i++) {
            const columns = sheet.Data[i]?.Data ?? [];
            let cols = columns.filter((_, idx) => !setting.Cols?.includes(idx));
            rows.push({ Columns: cols.length, Data: cols });
            if (maxCellCount < cols.length) maxCellCount = cols.length;
        }
        return { ...sheet, Data: rows, Columns: maxCellCount };
    });
}

export function deleteRows(sheets: Sheet[], setting: Setting): Sheet[] {
    return sheets.map((sheet) => {
        if (sheet.Name !== setting.Sheet || !sheet.Data) return sheet;
        let rows = sheet.Data.filter((_, idx) => !setting.Rows?.includes(idx));
        return { ...sheet, Data: rows, Rows: rows.length };
    });
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

export function setCurrentWorkbookId(id: string): void {
    localStorage.setItem("currentWorkbookId", id);
}

export function getCurrentWorkbookId(): string | null {
    return localStorage.getItem("currentWorkbookId");
}

export function clearCurrentWorkbookId(): void {
    localStorage.removeItem("currentWorkbookId");
}
