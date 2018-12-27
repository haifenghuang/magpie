# About defer keyword

 A defer statement defers the execution of a function until the surrounding function returns.

 The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns.

```swift
let add  =  fn(x,y){
    defer println("I'm defer1")
    println("I'm in add")
    defer println("I'm defer2")
    return x + y
}
println(add(2,2))
```

 The result is as below:

```text
I'm in add
I'm defer2
I'm defer1
4
```

```swift
file = newFile(filename, "r")
if (file == nil) {
    println("opening ", filename, "for reading failed, error:", file.message())
    return false
}
defer file.close()
//do other file related stuff, and not need to worry about the file close.
//when any file operation error occurs, it will close the file before it returns.
```

