# Optional type

Magpie has support for Optional type like java8.

```swift
fn safeDivision?(a, b) {
    if (b == 0){
        return optional.empty();
    } else {
        return optional.of(a/b);
    }
}

op1 = safeDivision?(10, 0)
if (!op1.isPresent()) {
    println(op1)

}

op2 = safeDivision?(10, 2)
if (op2.isPresent()) {
    println(op2)

    let val = op2.get()
    printf("safeDivision?(10, 2)=%d\n", int(val))
}
```

{% hint style="info" %}
 It is recommended that you use '?' as the last character of method to denote that it is an option.
{% endhint %}



