let a1 = 1.. 9
println(a1)

let rangefn = fn(x, y) { return x..y }
println(rangefn(1,9))

let start = fn(x) { x }
let end = fn(y) { y }
let a2 = start(1) .. end(9)
println(a2)

let arr = "a" .. "h"
println(arr)

let arr2 = [1..9, 10, 11]
println(arr2)
