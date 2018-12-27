# Error Handling of standard library

 When a standard library function returns `nil` or `false`, you can use the return value's message\(\) function for the error message:

```swift
file = newFile(filename, "r")
if (file == nil) {
    println("opening ", filename, "for reading failed, error:", file.message())
}
//do something with the file

//close the file
file.close()


let ret = http.listenAndServe("127.0.0.1:9090")
if (ret == false) {
    println("listenAndServe failed, error:", ret.message())
}
```

 Maybe you are curious about why `nil` or `false` have message\(\) function? Because in magpie, `nil` and `false` both are objects, so they have method to operate on it.

{% hint style="info" %}
In your Magpie code, It is recommended that use 'optional'.
{% endhint %}

