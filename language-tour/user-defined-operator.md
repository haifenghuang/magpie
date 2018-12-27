# User Defined Operator

 In magpie, you are free to define some operators, but you cannot overwrite predefined operators.

{% hint style="info" %}
 Not all operators could be user defined.
{% endhint %}

 Below is an example for showing how to write User Defined Operators:

```swift
//infix operator '=@' which accept two parameters.
fn =@(x, y) {
    return x + y * y
}

//prefix operator '=^' which accept only one parameter.
fn =^(x) {
    return -x
}

let pp = 10 =@ 5 // Use the '=@' user defined infix operator
printf("pp=%d\n", pp) // result: pp=35

let hh = =^10 // Use the '=^' prefix operator
printf("hh=%d\n", hh) // result: hh=-10
```

```swift
fn .^(x, y) {
    arr = []
    while x <= y {
        arr += x
        x += 2
    }
    return arr
}

let pp = 10.^20
printf("pp=%v\n", pp) // result: pp=[10, 12, 14, 16, 18, 20]
```

 Below is a list of predefined operators and user defined operators:

|  Predefined Operators |  User Defined Operators |
| :--- | :--- |
|  == =~ =&gt; |  =X |
|  ++ += |  +X |
|  -- -= -&gt; |  -X |
|  &gt;= &lt;&gt; |  &gt;X |
|  &lt;= &lt;&lt; |  &lt;X |
|  != !~ |  !X |
|  \*= \*\* |  \*X |
|  .. .. | .X |
|  &= && |  &X |
|  \|= \|\| |  \|X |
|  ^= |  ^X |

{% hint style="info" %}
 In the table above, `X` could be `.=+-*/%&,|^~<,>},!?@#$`
{% endhint %}



