# Regular expression

 In Magpie, regard to regular expression, you could use:

* Regular expression literal
* 'regexp' module
* =~ and !~ operators\(like perl's\)

```swift
//Use regular expression literal( /pattern/.match(str) )
let regex = /\d+\t/.match("abc 123	mnj")
if (regex) { println("regex matched using regular expression literal") }

//Use 'regexp' module
if regexp.compile(`\d+\t`).match("abc 123	mnj") {
    println("regex matched using 'regexp' module")
}

//Use '=~'(str =~ pattern)
if "abc 123	mnj" =~ `\d+\t` {
    println("regex matched using '=~'")
}else {
    println("regex not matched using '=~'")
}
```

{% hint style="info" %}
For detailed explanation of 'Regular Expression' pattern matching, you could see golang's regexp module for reference.
{% endhint %}

