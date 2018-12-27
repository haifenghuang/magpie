# String

 In Magpie, there are three types of `string`:

* Raw string
* Double quoted string\(Could not contains newline\)
* Single quoted string\(Interpolated String\)

 Raw string literals are character sequences between back quotes, as in `foo`. Within the quotes, any character may appear except back quote.

 See below for some examples:

```swift
normalStr = "Hello " + "world!"
println(normalStr)

println("123456"[2])

rawStr = `Welcome to
visit us!`
println(rawStr)

//when you use single quoted string, and want variable to be interpolated,
//you just put the variable into '{}'. see below:
str = "Hello world"
println('str={str}') //output: "Hello world"
str[6]="W"
println('str={str}') //output: "Hello World"
```

 In Magpie, strings are utf8-encoded, you could use utf-8 encoded name as a variable name.

```swift
三 = 3
五 = 5
println(三 + 五) //output : 8
```

 strings are also object, so you could use some of the methods provided by `strings` module.

```swift
upperStr = "hello world".upper()
println(upperStr) //output : HELLO WORLD
```

 string could also be iterated:

```swift
for idx, v in "abcd" {
    printf("idx=%d, v=%s\n", idx, v)
}

for v in "Hello World" {
    printf("idx=%d, v=%s\n", $_, v) //$_ is the index
}
```

 You could concatenate an object to a string:

```swift
joinedStr = "Hello " + "World"
joinedStr += "!"
println(joinedStr)
```

 You could also can use the builtin `reverse` function to reverse character of the string:

```swift
let str = "Hello world!"
println("Source Str =", str)
revStr = reverse(str)
println("Reverse str =", revStr)
```

 If you hava a string, you want to convert it to number, you could add a "+" prefix before the string.

```swift
a = +"121314" // a is an int
println(a) // result: 121314

// Integer also support "0x"(hex), "0b"(binary), "0o"(octal) prefix
a = +"0x10" // a is an int
println(a) // result: 16

a = +"121314.6789" // a is a float
println(a) // result: 121314.6789
```



