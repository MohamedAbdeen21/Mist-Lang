fn multiply(multiplier: Int) Func {
    return fn(x:Int) Int {return x*multiplier;};
}

fn main() {
    let printer: Func = fn(x: Int) {print(string(x) + " ")};
    
    range(0,25).map(multiply(4)).map(printer);
    println();
}
