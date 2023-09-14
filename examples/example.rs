fn main() {
    let list: List = [1, 2, 1];
    let list: List = list.map(fn(x:Int)Int{if (x == 1) {10} else {x}});
    println(list);
}
