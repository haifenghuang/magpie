# Operator overloading

Magpie also support operator overloading:

```swift
class Vector {
    let x = 0;
    let y = 0;

    // constructor
    fn init (a, b, c) {
        if (!a) { a = 0;}
        if (!b) {b = 0;}
        x = a; y = b
    }

    fn +(v) { //overloading '+'
        if (type(v) == "INTEGER" {
            return new Vector(x + v, y + v);
        } elseif v.is_a(Vector) {
            return new Vector(x + v.x, y + v.y);
        }
        return nil;
    }

    fn String() {
        return fmt.sprintf("(%v),(%v)", this.x, this.y);
    }
}

fn Vectormain() {
    v1 = new Vector(1,2);
    v2 = new Vector(4,5);
    
    // call + function in the vector object
    v3 = v1 + v2 //same as 'v3 = v1.+(v2)'
    // returns string "(5),(7)"
    println(v3.String());
    
    v4 = v1 + 10 //same as v4 = v1.+(10);
    //returns string "(11),(12)"
    println(v4.String());
}

Vectormain()
```

