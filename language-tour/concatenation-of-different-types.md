# Concatenation of different types

 In Magpie, you could concatenate of different types. See below for examples:

```swift
// Number plus assignment
num = 10
num += 10 + 15.6
num += 20
println(num)

// String plus assignment
str = "Hello "
str += "world! "
str += [1, 2, 3]
println(str)

// Array plus assignment
arr = []
arr += 1
arr += 10.5
arr += [1, 2, 3]
arr += {"key": "value"}
println(arr)

// Array compare
arr1 = [1, 10.5, [1, 2, 3], {"key" : "value"}]
println(arr1)
if arr == arr1 { //support ARRAY compare
    println("arr1 = arr")
} else {
    println("arr1 != arr")
}

// Hash assignment("+=", "-=")
hash = {}
hash += {"key1" : "value1"}
hash += {"key2" : "value2"}
hash += {5 : "five"}
println(hash)
hash -= "key2"
hash -= 5
println(hash)
```

