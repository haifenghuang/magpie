# Variables

 Variables in Magpie could start with the keyword `let`, or nothing with the form `variable=value`.

```swift
let a, b, c = 1, "hello world", [1,2,3]
d = 4
e = 5
姓 = "黄"
```

 You can also use `Destructuring assignment`. Note, the left-hand side must be included using the '\(\)'.

```text
//righ-hand side is an array
let (d,e,f) = [1,5,8]
//d=1, e=5, f=8

//right-hand side is a tuple
let (g, h, i) = (10, 20, "hhf")
//g=10, h=20, i=hhf

//righ-hand side is a hash
let (j, k, l) = {"j": 50, "l": "good"}
//j=50, k=nil, l=good
```

 Note however, if you do not use the keyword `let`, you could not do multiple variable assignments. Below code is not correct：

```swift
//Error, multiple variable assignments must be use `let` keyword
a, b, c = 1, "hello world", [1,2,3]
```

{% hint style="info" %}
 Starting from Magpie 5.0，when the decalared variable already exists, it's value will be overwritten:
{% endhint %}

```swift
let x, y = 10, 20;
let x, y = y, x //Swap the value of x and y
printf("x=%v, y=%v\n", x, y)  //result: x=20, y=10
```

 `let` also support the placeholder\(\_\), when assigned a value, it will just ignore it.

```swift
let x, _, y = 10, 20, 30
printf("x=%d, y=%d\n", x, y) //result: x=10, y=30
```

