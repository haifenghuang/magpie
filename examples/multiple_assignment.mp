a,b,c = [1,3,5,2,4,6,7,8,9], 5, "hello"
let d = []

println('a={a}')
println('b={b}')
println('c={c}')
println('d={d}')

a, c = c, a
println('a={a}')
println('c={c}')

m = fn() { return 10, "hello" }
x, y, z = m(), fn(x, y) { x + y }(1,2)
println('x={x}')
println('y={y}')
println('z={z}')

/*
    Below examples demonstrate Destructuring assignment.
    Note, however that you must use 'let' keyword and the
    left part of the assignment must be enclosed in parentheses.
*/
let (d,e,f) = [1,5,8]
printf("d=%v, e=%v, f=%v\n", d, e, f)
//d=1, e=5, f=8

//right-hand side is a tuple
let (g, h, i) = (10, 20, "hhf")
printf("g=%v, h=%v, i=%v\n", g, h, i)
//g=10, h=20, i=hhf

//righ-hand side is a hash
let (j, k, l) = {"j": 50, "l": "good"}
printf("j=%v, k=%v, l=%v\n", j, k, l)
//j=50, k=nil, l=good


//////////////////////////////////////////////////
// now, please see the difference
//////////////////////////////////////////////////
//using Destructuring assignment
let (d,e,f) = [1,5,8]
printf("d=%v, e=%v, f=%v\n", d, e, f) //d=1, e=5, f=8

//not using Destructuring assignment
let d,e,f = [1,5,8]
printf("d=%v, e=%v, f=%v\n", d, e, f) //d=[1, 5, 8], e=nil, f=nil