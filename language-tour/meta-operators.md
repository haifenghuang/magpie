# Meta-Operators

 Magpie has some `meta-operators` borrowed from perl6. There are strict rules for meta-operators:

* Meta-operators can only operator on arrays.
* Each array's element must be number type\(uint, int, float\) or string type.
* If the meat-operators serve as an infix operator, and if the left and right are all arrays, they must have the same number of elements.

```swift
let arr1 = [1,2,3] ~* [4,5,6]
let arr2 = [1,2,3] ~* 4
let arr3 = [1,2,"HELLO"] ~* 2
let value1 = ~*[10,2,2]
let value2 = ~+[2,"HELLO",2]

println(arr1)   //result: [4, 10, 18]
println(arr2)   //result: [4, 8, 12]
println(arr3)   //result: [2,4,"HELLOHELLO"]
println(value1) //result: 40
println(value2) //result: 2HELLO2
```

 At the moment, Magpie has six meta-operators：

*  ~+
*  ~-
*  ~\*
*  ~/
*  ~%
*  ~^

 The six meta-operators could be served as either infix expression or prefix expression.

 The meta-operator for infix expression will return an array. The meta-operator for prefix expression will return a value\(uint, int, float, string\).

 Below talbe give an example of meta-operator and their meanings:\(only `~+` is showed\):

|  Meta-Operator |  Expression |  Example |  Result |
| :--- | :--- | :--- | :--- |
|  ~+ |  Infix Expression |  \[x1, y1, z1\] ~+ \[x2, y2, z2\] |  \[x1+x2, y1+y2, z1+z2\] \(Array\) |
|  ~+ |  Infix Expression |  \[x1, y1, z1\] ~+ 4 |  \[x1+4, y1+4, z1+4\] \(Array\) |
| ~+ |  Prefix Expression |  ~+\[x1, y1, z1\] |  x1+y1+z1 \(Note: a value, not an array\) |



