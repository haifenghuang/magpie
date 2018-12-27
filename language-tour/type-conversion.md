# Type conversion

 You can use the builtin `int()`, `uint()`, `float()`, `str()`, `array()`, `tuple()`, `hash`, `decimal` functions for type conversion.

```swift
let i = 0xa
let u = uint(i)                 // result: 10
let s = str(i)                  // result: "10"
let f = float(i)                // result: 10
let a = array(i)                // result: [10]
let t = tuple(i)                // result:(10,)
let h = hash(("key", "value"))  // result: {"key": "value}
let d = decimal("123.45634567") // result: 123.45634567
```

 You could create a tuple from an array:

```text
let t = tuple([10, 20])   //result:(10,20)
```

 Similarly, you could also create an array from a tuple:

```text
let arr = array((10,20))  //result:[10,20]
```

 You could only create a hash from an array or a tuple:

```swift
//create an empty hash
let h1 = hash()  //same as h1 = {}

//create a hash from an array
let h1 = hash([10, 20])     //result: {10 : 20}
let h2 = hash([10,20,30])   //result: {10 : 20, 30 : nil}

//create a hash from a tuple
let h3 = hash((10, 20))     //result: {10 : 20}
let h4 = hash((10,20,30))   //result: {10 : 20, 30 : nil}
```

