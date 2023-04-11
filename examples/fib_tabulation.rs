fn fib(table: List, x: Int) List {
    if (x < 2) {
        return fib(table, x+1);
    } else if (x > table.len() - 1) {
        // filled the entire table, return
        return table;
    } else {
        return fib(table.update(x, table[x-1] + table[x-2]), x + 1);
    }
}

fn main() {
    let table: List = range(0,50);
    fib(table, 0).map(println);
    return;
}
