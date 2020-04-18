let a = 10

fn printMe(msg) {
    msg = "** " + msg + " **"
    println(msg)
    demo(8)
}

let demo = fn(a) {
    printf("a=%d\n", a)
    if (a > 5) {
        println("a > 5")
        return
    } else {
        println("a<=5")
    }
}

for i in 1..5 {
    printf("i=%d\n", i)
}

demo(4)

if (a > 12) {
    printMe("a>12")
} else {
    printMe("a < 12")
}

println(len("Hello World"))

let x = 5 let y = "good"
let add = fn(x, y) {
    return x + y
}

sum = add(10, 5)
printf("x=%d, len(%s)=%d, 10+5=%d\n", x, y, strings.len(y), sum)

println("Program end.")
