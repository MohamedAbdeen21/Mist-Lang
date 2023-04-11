fn main() {
    let printer: Func = fn(x: Int) {print(string(x) + " ")};
    print("Evens: ");
    println(range(0,100).filter(fn(x:Int) Bool {x%2 == 0}).len());
    print("Odds: ");
    println(range(0,100).filter(fn(x:Int) Bool {x%2 != 0}).len());
    return;
}
