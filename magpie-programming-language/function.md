# Function

 Function in magpie is a first-class object. This means the language supports passing functions as arguments to other functions, returning them as the values from other functions, and assigning them to variables or storing them in data structures.

 Function also could have default parameters and variadic parameters.

```swift
//define a function
let add = fn() { [5,6] }
let n = [1, 2] + [3, 4] + add()
println(n)


let complex = {
   "add" : fn(x, y) { return fn(z) {x + y + z } }, //function with closure
   "sub" : fn(x, y) { x - y },
   "other" : [1,2,3,4]
}
println(complex["add"](1, 2)(3))
println(complex["sub"](10, 2))
println(complex["other"][2])


let warr = [1+1, 3, fn(x) { x + 1}(2),"abc","def"]
println(warr)


println("\nfor i in 5..1 where i > 2 :")
for i in fn(x){ x+1 }(4)..fn(x){ x+1 }(0) where i > 2 {
  if (i == 3) { continue }
  println('i={i}')
}


// default parameter and variadic parameters
add = fn (x, y=5, z=7, args...) {
    w = x + y + z
    for i in args {
        w += i
    }
    return w
}

w = add(2,3,4,5,6,7)
println(w)
```

 You could also declare named function like below:

```swift
fn sub(x,y=2) {
    return x - y
}
println(sub(10)) //output : 8
```

 You could also create a function using the `fat arrow` syntax:

```swift
let x = () => 5 + 5
println(x())  //result: 10

let y = (x) => x * 5
println(y(2)) //result: 10

let z = (x,y) => x * y + 5
println(z(3,4)) //result :17


let add = fn (x, factor) {
  x + factor(x)
}
result = add(5, (x) => x * 2)
println(result)  //result : 15
```

 If the function has no parameter, then you could omit the parentheses. e.g.

```swift
println("hhf".upper)  //result: "HHF"
//Same as above
println("hhf".upper())
```

 Before ver5.0, Magpie do not support multiple return values, But there are many ways to do it.

 Below suggest a way of doing it:

```swift
fn div(x, y) {
    if y == 0 {
        return [nil, "y could not be zero"]
    }
    return [x/y, ""]
}

ret = div(10,5)
if ret[1] != "" {
    println(ret[1])
} else {
    println(ret[0])
}
```

 Starting from ver5.0, Magpie support multiple return values using 'let'. The returned values are wrapped as a tuple.

```swift
fn testReturn(a, b, c, d=40) {
    return a, b, c, d
}

let (x, y, c, d) = testReturn(10, 20, 30)
// let x, y, c, d = testReturn(10, 20, 30)  same as above

printf("x=%v, y=%v, c=%v, d=%v\n", x, y, c, d)
//Result: x=10, y=20, c=30, d=40
```

{% hint style="info" %}
 You must use `let` to support multiple return values, below statement will issue a compile error.
{% endhint %}

```swift
(x, y, c, d) = testReturn(10, 20, 30) // no 'let', compile error
x, y, c, d = testReturn(10, 20, 30)   // no 'let', compile error
```

