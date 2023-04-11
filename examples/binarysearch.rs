fn BSearch(lst: List, target: Int) Int {
    let index: Int = lst.len() / 2;
  
    if (lst.len() == 1 && lst[0] != target) {
        return -1;
    } else if (lst[index] == target) {
        return lst[index];
    } else if (lst[index] < target) {
        return BSearch(lst.slice(index, lst.len()-1), target);
    } else {
        return BSearch(lst.slice(0, index-1), target);
    }
}

fn main() {
    println(BSearch([1], 1));
    println(BSearch([1, 2, 3, 4, 5, 6], 10));
    println(BSearch([1, 2, 3, 4, 5, 6], 5));
    return;
}
