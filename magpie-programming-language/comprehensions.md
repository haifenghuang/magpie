# Comprehensions

 Magpie support list\(array,string, range, tuple\) comprehensions. list comprehension will return an array. please see following examples:

```swift
//array comprehension
x = [[word.upper(), word.lower(), word.title()] for word in ["hello", "world", "good", "morning"]]
println(x) //result: [["HELLO", "hello", "Hello"], ["WORLD", "world", "World"], ["GOOD", "good", "Good"], ["MORNING", "morning", "Morning"]]

//string comprehension (here string is treated like an array)
y = [ c.upper() for c in "huanghaifeng" where $_ % 2 != 0] //$_ is the index
println(y) //result: ["U", "N", "H", "I", "E", "G"]

//range comprehension
w = [i + 1 for i in 2..10]
println(w) //result: [2, 3, 4, 5, 6, 7, 8, 9, 10, 11]

//tuple comprehension
v = [x+1 for x in (12,34,56)]
println(v) //result: [13, 35, 57]

//hash comprehension
z = [v * 10 for k,v in {"key1": 10, "key2": 20, "key3": 30}]
println(z) //result: [100, 200, 300]
```

 Magpie also support hash comprehension. hash comprehension will return a hash. please see following examples:

```swift
//hash comprehension (from hash)
z1 = { v:k for k,v in {"key1": 10, "key2": 20, "key3": 30}} //reverse key-value pair
println(z1) // result: {10 : "key1", 20 : "key2", 30 : "key3"}. Order may differ

//hash comprehension (from array)
z2 = {x:x**2 for x in [1,2,3]}
println(z2) // result: {1 : 1, 2 : 4, 3 : 9}. Order may differ

//hash comprehension (from .. range)
z3 = {x:x**2 for x in 5..7}
println(z3) // result: {5 : 25, 6 : 36, 7 : 49}. Order may differ

//hash comprehension (from string)
z4 = {x:x.upper() for x in "hi"}
println(z4) // result: {"h" : "H", "i" : "I"}. Order may differ

//hash comprehension (from tuple)
z5 = {x+1:x+2 for x in (1,2,3)}
println(z5) // result: {4 : 5, 2 : 3, 3 : 4}. Order may differ
```

