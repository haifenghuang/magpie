//#!/usr/bin/magpie
const a = 3.14159
//a = 3.14 //error
printf("PI=%f\n", a)

fn Test() {
  let a = 15 //good
  printf("In Test(), a=%d\n", a)
}

Test()

const b = fn(x) {
  println(x)
}

//b = a //error
b("Hello")