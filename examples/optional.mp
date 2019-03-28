fn safeDivision?(a, b) {
    if (b == 0){
        return optional.empty();
    } else {
        return optional.of(a/b);
    }
}

op1 = safeDivision?(10, 0)
if (!op1.isPresent()) { //same as 'if (!op1)'
    println(op1)

}

op2 = safeDivision?(10, 2)
if (op2.isPresent()) { // same as 'if (op2)'
    println(op2)

    let val = op2.get()
    printf("safeDivision?(10, 2)=%d\n", int(val))
}

println("-------------------------------\n")

class Address {
    let block;
    let city;
    let state;
    let country;

    fn init(block, city, state, country) {
        this.block = block
        this.city = city
        this.state = state
        this.country = country
    }

    property Block   { get; }
    property City    { get; }
    property State   { get; }
    property Country { get; }

    fn str() {
        return "Address: " + this.block + ", " + this.city + ", " + this.state + ", " + this.country
    }
}

class Person {
    let name;
    let address;

    fn init(name, address) {
        this.name = name
        this.address = address
    }

    property Name    { get; }

    // Returns an Optional.
    // If a function returns an optional, Suggestion is that you use the '?' to denote it,
    // though it is not mandatory.
    fn getAddress?() {
//        if (this.address == nil) {
//            return optional.empty()
//        }
//        return optional.of(this.address)

        //same as above
        return optional.ofNullable(this.address)
    }
}

person1 = new Person("hhf", new Address("block", "city", "state", "country"))
option1 = person1.getAddress?();
option1.ifPresent(fn(addr) { // if it has address, then print the address
  println(addr.str())
})


let EMPTY_ADDRESS = new Address("", "", "", "");
person2 = new Person("hhf", nil)
addr = person2.getAddress?().orElse(EMPTY_ADDRESS); //if person2's address is nil, then use EMPTY_ADDRESS
println(addr.str())

option2 = person2.getAddress?();
option2.ifPresentOrElse(
    (addr) => { println(addr.str()) },
    ()     => { println("Address not present")}
)

option3 = optional.of(20)
filter = option3.filter(
  (a) => { a < 10 }
)
if (filter == optional.empty()) {
    println("a is not less than 10")
}

optionalOfMap = option3.filter( (a) => { a < 10 } )
                       .map(    (a) => optional.of(a + 1));
if (optionalOfMap == optional.empty()) {
    println("a is not less than 10")
}

currentValue = optional.of(20);
result = currentValue.filter((a) => { a < 30} )
                     .map(   (a) => { a+1 });
orElse = result.orElse(555);
if (orElse == 21) {
    println("orElse is equal to 21")
}
