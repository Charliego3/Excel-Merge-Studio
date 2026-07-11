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
