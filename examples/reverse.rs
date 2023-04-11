fn reverse(x: List) List {
    if (x.len() == 1) {return x;}
    return [x[x.len() - 1]] + reverse(x.slice(0,x.len() - 2));
}

fn main() {
    let printer: Func = fn(x: Int) {print(string(x) + " ")};
    let init: List = [1,1,1,2,3,3,6];
    
    init.map(printer);
    println();
    reverse(init).map(printer);;
    println();
    return;
}
