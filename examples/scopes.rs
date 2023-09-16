fn main() {
    let x: Int = 10;

    let y: Int = if (true) {
        let x: Int = 3;
        print(x);
        println();
        x; // an implicit return. unlike rust, not suppressed by ';'.
    }; // can chain more else-if or else

    print(y);
    println();
    print(x)
    println();
    return;
}
