# Tuple

 In Magpie, `tuple` is just like array, but it could not be changed once it has been created.

 Tuples are constructed using parenthesized list notation:

```swift
//Create an empty tuple
let t1 = tuple()

//Same as above.
let t2 = ()

// Create a one element tuple.
// Note: the trailing comma is necessary to distinguish it from the
//       parenthesized expression (1).
// 1-tuples are seldom used.
let t3 = (1,)

//Create a two elements tuple
let t4 = (2,3)
```

 Any object may be converted to a tuple by using the built-in `tuple` function.

```swift
let t = tuple("hello")
println(t)  // result: ("hello")
```

 Like arrays, tuples are indexed sequences, so they may be indexed and sliced. The index expression tuple\[i\] returns the tuple element at index i, and the slice expression tuple\[i:j\] returns a subsequence of a tuple.

```swift
let t = (1,2,3)[2]
print(t) // result:3
```

 Tuples are iterable sequences, so they may be used as the operand of a for-loop, a list comprehension, or various built-in functions.

```swift
//for-loop
for i in (1,2,3) {
    println(i)
}

//tuple comprehension
let t1 =  [x+1 for x in (2,4,6)]
println(t1) //result: [3, 5, 7]. Note: Result is array, not tuple
```

 Unlike arrays, tuples cannot be modified. However, the mutable elements of a tuple may be modified.

```swift
arr1 = [1,2,3]
t = (0, arr1, 5, 6)
println(t)    // result: (0, [1, 2, 3], 5, 6)
arr1.push(4)
println(t)    //result:  (0, [1, 2, 3, 4], 5, 6)
```

 Tuples are hashable \(assuming their elements are hashable\), so they may be used as keys of a hash.

```swift
key1=(1,2,3)
key2=(2,3,4)
let ht = {key1 : 10, key2 : 20}
println(ht[key1]) // result: 10
println(ht[key2]) // result: 20

//Below is not supported(will issue a syntax error):
let ht = {(1,2,3) : 10, (2,3,4) : 20} //error!
println(ht[(1,2,3)])  //error!
println(ht[(2,3,4)])  //error!
```

 Tuples may be concatenated using the + operator, it will create a new tuple.

```swift
let t = (1, 2) + (3, 4)
println(t) // result: (1, 2, 3, 4)
```

 A tuple used in a Boolean context is considered true if it is non-empty.

```swift
let t = (1,)
if t {
    println("t is not empty!")
} else {
    println("t is empty!")
}

//result : "t is not empty!"
```

 Tuple's json marshaling and unmarshaling will be treated as array:

```swift
let tupleJson = ("key1","key2")
let tupleStr = json.marshal(tupleJson)
//Result:[
//        "key1"，
//        "key2"，
//       ]
println(json.indent(tupleStr, "  "))

let tupleJson1 = json.unmarshal(tupleStr)
println(tupleJson1) //result: ["key1", "key2"]
```

 Tuple plus an array will return an new array, not a tuple

```swift
t2 = (1,2,3) + [4,5,6]
println(t2) // result: [(1, 2, 3), 4, 5, 6]
```

 You could also use the builtin `reverse` function to reverse tuples's elements:

```swift
let tp = (1,3,5,2,4,6,7,8,9)
println(tp) //result: (1, 3, 5, 2, 4, 6, 7, 8, 9)

revTuple = reverse(tp)
println(revTuple) //result: (9, 8, 7, 6, 4, 2, 5, 3, 1)
```

