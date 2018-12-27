# Array

 In Magpie, you could use \[\] to initialize an empty array:

```swift
emptyArr = []
emptyArr[3] = 3 //will auto expand the array
println(emptyArr)
```

 You can create an array with the given size\(or length\) using below two ways:

```swift
//create an array with 10 elements initialized to nil.
//Note: this only support integer literal.
let arr = []10
println(arr)

//use the builtin 'newArray' method.
let anotherArr = newArray(len(arr))
println(anotherArr) //result: [nil, nil, nil, nil, nil, nil, nil, nil, nil, nil]

let arr1 = ["1","a5","5", "5b","4","cc", "7", "dd", "9"]
let arr2 = newArray(6, arr1, 10, 11, 12) //first parameter is the size
println(arr2) //result: ["1", "a5", "5", "5b", "4", "cc", "7", "dd", "9", 10, 11, 12]

let arr3 = newArray(20, arr1, 10, 11, 12)
println(arr3) //result : ["1", "a5", "5", "5b", "4", "cc", "7", "dd", "9", 10, 11, 12, nil, nil, nil, nil, nil, nil, nil, nil]
```

 Array could contain any number of different data types. Note: the last comma before the closing '\]' is optional.

```swift
mixedArr = [1, 2.5, "Hello", ["Another", "Array"], {"Name": "HHF", "SEX": "Male"}]
```

 You could use index to access array element.

```swift
println('mixedArr[2]={mixedArr[2]}')
println(["a", "b", "c", "d"][2])
```

 Because array is an object, so you could use the object's method to operate on it.

```swift
if ([].empty()) {
    println("array is empty")
}

emptyArr.push("Hello")
println(emptyArr)

//you could also use 'addition' to add an item to an array
emptyArr += 2
println(emptyArr)

//You could also use `<<(insertion operator)` to add item(s) to an array, the insertion operator support chain operation.
emptyArr << 2 << 3
println(emptyArr)
```

 Array could be iterated using `for` loop

```swift
numArr = [1,3,5,2,4,6,7,8,9]
for item in numArr where item % 2 == 0 {
    println(item)
}

let strArr = ["1","a5","5", "5b","4","cc", "7", "dd", "9"]
for item in strArr where /^\d+/.match(item) {
    println(item)
}

for item in ["a", "b", "c", "d"] where $_ % 2 == 0 {  //$_ is the index
    printf("idx=%d, v=%s\n", $_, item)
}
```

 You could also use the builtin `reverse` function to reverse array element:

```swift
let arr = [1,3,5,2,4,6,7,8,9]
println("Source Array =", arr)

revArr = reverse(arr)
println("Reverse Array =", revArr)
```

 Array also support `array multiplication operator`\(\*\):

```swift
let arr = [3,4] * 3
println(arr) // result: [3,4,3,4,3,4]
```

