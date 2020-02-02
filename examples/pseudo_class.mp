// A pseudo-class using function with closure
fn Person(name, age) {
    self = {}
    self.name = name
    self.age = age

    //You could also use 'self.getName = fn() { return self.name }'
    self.getName = () => self.name
    self.getAge  = () => self.age
    self.message = () => self.name + ", aged " + str(self.age)

    self.sets = fn(newName, newAge) {
        self.name = newName
        self.age = newAge
    }

    self
}

p = Person("Bob", 40)
printf("1 - info = %v\n", p.message())
printf("1 - name = %v\n", p.getName())
printf("1 - age  = %v\n", p.getAge())

printf("\n=========================\n\n")

p.sets("HHF", 42)
printf("2 - info = %v\n", p.message())
printf("2 - name = %v\n", p.getName())
printf("2 - age  = %v\n", p.getAge())