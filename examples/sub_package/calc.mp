//Identifiers and functions which have their first 
//letter 'Uppercased' will be exported.

fn _add(x, y) {
    x + y
}
//exported
fn Add(x, y) {
    return _add(x, y)
}

//exported
fn Minus(x, y) {
    x - y
}

//not exported, the first letter is not Uppercase.
fn multiply(x, y) {
    x * y
}

const (
    SQUARE_OF_TWO = 4,   //exported
    sQUARE_OF_FOUR = 16, //not exported
)