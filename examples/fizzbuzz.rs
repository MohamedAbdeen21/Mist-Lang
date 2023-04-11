fn main() {
    fn fizzbuzz(x:Int) String {
        ("fizz" * (x%3 == 0) + "buzz" * (x%5 == 0)).otherwise(string(x))
    }
    range(1,100)
        .map(fizzbuzz)
        .reverse()
        .map(println);
    return;
}
