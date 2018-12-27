# grep and map

 The `grep` and `map` operators are just like perl's `grep` and `map`.

 The grep operator takes a list of values and a "testing expression." For each item in the list of values, the item is placed temporarily into the $\_ variable, and the testing expression is evaluated. If the expression results in a true value, the item is considered selected.

 The map operator has a very similar syntax to the grep operator and shares a lot of the same operational steps. For example, items from a list of values are temporarily placed into $\_ one at a time. However, the testing expression becomes a mapping expression.

```swift
let sourceArr = [2,4,6,8,10,12]

let m = grep  $_ > 5, sourceArr
println('m is {m}')

let cp = map $_ * 2 , sourceArr
println('cp is {cp}')

//a little bit more complex example
let fields = {
                "animal"   : "dog",
                "building" : "house",
                "colour"   : "red",
                "fruit"    : "apple"
             }
let pattern = `animal|fruit`
// =~(match), !~(unmatch)
let values = map { fields[$_] } grep { $_ =~ pattern } fields.keys()
println(values)
```

