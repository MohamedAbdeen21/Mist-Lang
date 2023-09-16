fn multiply(multiplier: Int) Func {
    return fn(x:Int) Int {return x*multiplier;};
}

fn main() {
    let printer: Func = fn(x: Int) {print(string(x) + " ")};
    print(multiply(4)(2));
    
    range(0,25).map(multiply(4)).map(printer);
    println();
}
