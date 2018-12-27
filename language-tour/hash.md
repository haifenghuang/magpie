# Hash

 In magpie, the builtin hash will keep the order of keys when they are added to the hash, just like python's orderedDict.

 You could use `{}` to initialize an empty hash:

```swift
emptyHash = {}
emptyHash["key1"] = "value1"
println(emptyHash)
```

 Hash's key could be string, int, boolean:

```swift
hashObj = {
    12     : "twelve",
    true   : 1,
    "Name" : "HHF"
}
println(hashObj)
```

{% hint style="info" %}
 The last comma before the closing '}' is optional.
{% endhint %}

 You could use `+` or `-` to add or remove an item from a hash:

```swift
hashObj += {"key1" : "value1"}
hashObj += {"key2" : "value2"}
hashObj += {5 : "five"}
hashObj -= "key2"
hashObj -= 5
println(hash)
```

 In Magpie, Hash is also an object, so you could use them to operate on hash object:

```swift
hashObj.push(15, "fifteen") //first parameter is the key, second is the value
hashObj.pop(15)

keys = hashObj.keys()
println(keys)

values = hashObj.values()
println(values)
```

 You could also use the builtin `reverse` function to reverse hash's key and value:

```swift
let hs = {"key1": 12, "key2": "HHF", "key3": false}
println("Source Hash =", hs)
revHash = reverse(hs)
println("Reverse Hash =", revHash)
```

