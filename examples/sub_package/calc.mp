//exported
fn Add(x, y) {
    x + y
}

//exported
fn Minus(x, y) {
    x - y
}

//not exported, the first letter is not Uppercase.
fn multiply(x, y) {
    x * y
}

let SQUARE_OF_TWO = 4
let sQUARE_OF_FOUR = 16 //not exported