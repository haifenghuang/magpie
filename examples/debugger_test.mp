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

println("Program end.")

