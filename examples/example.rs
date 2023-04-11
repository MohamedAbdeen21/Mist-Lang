fn main() {
    let list: List = [1, 2];
    let list: List = list.map(fn(x:String)Int{if (x == 1) {10} else {x}});
    println(list[0]);
}
