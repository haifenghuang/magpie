fn safeDivision(a, b) {
    if (b == 0){
        return optional.empty();
    } else {
        return optional.of(a/b);
    }
}

op1 = safeDivision(10, 0)
if (!op1.isPresent()) {
    println(op1)

}

op2 = safeDivision(10, 2)
if (op2.isPresent()) {
    println(op2)

    let val = op2.get()
    printf("safeDivision(10, 2)=%d\n", int(val))
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
    fn getAddress() {
        if (this.address == nil) {
            return optional.empty()
        }
        return optional.of(this.address)
   }
}

person1 = new Person("hhf", new Address("block", "city", "state", "country"))
option1 = person1.getAddress();
option1.ifPresent(fn(addr) { // if has address, then, print the address
  println(addr.str())
})


let EMPTY_ADDRESS = new Address("", "", "", "");
person2 = new Person("hhf", nil)
addr = person2.getAddress().orElse(EMPTY_ADDRESS); //if person2's address is nil, then use EMPTY_ADDRESS
printf(addr.str())

