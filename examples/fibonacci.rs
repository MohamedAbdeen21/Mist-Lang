fn fib(x: Int) Int {
    if (x < 2) {x} else {fib(x-2) + fib(x-1)}
}

fn main() {
    range(0,25).map(fib).map(println);
    return;
}
