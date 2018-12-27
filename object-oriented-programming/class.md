# Class

 Magpie has limited support for the oop concept, below is a list of features:

* inheritance and polymorphism
* operator overloading
* property\(with getter or setter or both\)
* static member/method/property
* indexer
* class category
* class annotations\(limited support\)
* constructor method and normal methods support default value and variadic parameters

 The Magpie parser could parse `public`, `private`, `protected`, but it has no effect in the evaluation phase. That means Magpie do not support access modifiers at present.

 You use `class` keyword to declare a class and use `new class(xxx)` to create an instance of a `class`.

```swift
class Animal {
    let name = ""
    fn init(name) {    //'init' is the constructor
        //do somthing
    }
}
```

 In magpie, all class is inherited from the root class `object`. `object` class include some common method like `toString()`, `instanceOf()`, `is_a()`, `classOf()`, `hashCode`.

 Above code is same as:

```swift
class Animal : object {
    let name = ""
    fn init(name) {    //'init' is the constructor
        //do somthing
    }
}
```

