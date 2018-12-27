# Class Category

 Magpie also support class Category like objective-c（C\# is called 'extension methods'）.

```swift
class Animal {
    fn Walk() {
        println("Animal Walk!")
    }
}

//Class category like objective-c
class Animal (Run) { //Create an 'Run' category of Animal class.
    fn Run() {
        println("Animal Run!")
        this.Walk() //can call Walk() method of Animal class.
    }
}

animal = new Animal()
animal.Walk()

println()
animal.Run()
```

