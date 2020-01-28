<p align="center">
    <img alt="abs language logo" src="https://github.com/haifenghuang/magpie/blob/master/magpie.png?raw=true" width="310">
</p>

# Magpie Programming Language

Chinese version: [中文](README_cn.md)

## Summary

Magpie is a toy language interpreter, written in Go. It has C-style syntax, and is largely
inspired by Ruby, Python, Perl and c#.

It support the normal control flow, functional programming and object oriented programming.
and also can import golang's module.

It has a built-in documentation generator(mdoc) for generating html document from magpie source.

It also has a REPL with realtime syntax highlighter.

## Documention

Complete language tutorial can be found in [docs](docs)

## Features

* Class with support for property, indexer & operator overloading
* await/async for asynchronous programming
* Builtin support for linq
* First class function
* function with Variadic parameters and default values
* function with multiple return values
* int, uint, float, bool, array, tuple, hash(all support json marshal & unmarshal, all can be extended)
* try-catch-finally exception handling
* Optional Type support(Java 8 like)
* using statment(C# like)
* Elixir like pipe operator
* Using method of Go Package(RegisterFunctions and RegisterVars)
* Syntax-highlight REPL
* Doc-generation tool `mdoc`
* Integrated services processing

## Example1(Linq)

```csharp
// async/await
async fn add(a, b) { a + b }

result = await add(1, 2)
println(result)

// linq example
class Linq {
    static fn TestSimpleLinq() {
        //Prepare Data Source
        let ingredients = [
            {Name: "Sugar",  Calories: 500},
            {Name: "Egg",    Calories: 100},
            {Name: "Milk",   Calories: 150},
            {Name: "Flour",  Calories: 50},
            {Name: "Butter", Calories: 200},
        ]

        //Query Data Source
        ingredient = from i in ingredients where i.Calories >= 150 orderby i.Name select i

        //Display
        for item in ingredient => println(item)
    }

    static fn TestFileLinq() {
        //Read Data Source from file.
        file = newFile("./examples/linqSample.csv", "r")

        //Query Data Source
        result = from field in file where int(field[1]) > 300000 select field[0]

        //Display
        for item in result => printf("item = %s\n", item)

        file.close()
    }

    static fn TestComplexLinq() {
        //Data Source
        stringList = [
            "A penny saved is a penny earned.",
            "The early bird catches the worm.",
            "The pen is mightier than the sword."
        ]

        //Query Data Source
        earlyBirdQuery =
            from sentence in stringList
            let words = sentence.split(" ")
            from word in words
            let w = word.lower()
            where w[0] == "a" || w[0] == "e" ||
                  w[0] == "i" || w[0] == "o" ||
                  w[0] == "u"
            select word

        //Display
        for v in earlyBirdQuery => printf("'%s' starts with a vowel\n", v)
    }
}

Linq.TestSimpleLinq()
println("======================================")
Linq.TestFileLinq()
println("======================================")
Linq.TestComplexLinq()
```

## Example2(Rest Service)

```csharp
//service Hello on "0.0.0.0:8090" {
service Hello on "0.0.0.0:8090:debug" { //':debug': for debugging request
  //In '@route', you could use 'url(must), methods, host, schemes, headers, queries'
  @route(url="/authentication/login", methods=["POST"])
  fn login(writer, request) {
    //writer.writeJson({ sessionId: "3d5bd2cA15ef047689" })
    //writer.writeJson({ sessionId: "3d5bd2cA15ef047689" }), 200 # same as above
    //return { sessionId: "3d5bd2cA15ef047689" }, 200 # same as above
    return { sessionId: "3d5bd2cA15ef047689" } // same as above
  }

  @route(url="/authentication/logout", methods=["POST"])
  fn logout(writer, request) {
    // writer.writeHeader(http.STATUS_CREATED) # return http status code 201
    return http.STATUS_CREATED // same as above
  }

  @route(url="/meters/setting-result/{acceptNo}", methods=["GET"])
  fn load_survey_result(writer, request) {
    //using 'vars' dictionary to access the url parameters
    //writer.writeJson({ acceptNo: vars["acceptNo"], resultCode: "1"})
    return { acceptNo: vars["acceptNo"], resultCode: "1"} // same as above
  }

  @route(url="/articles/{category}/{id:[0-9]+}", methods=["GET"])
  fn getArticle(writer, request) {
    //using 'vars' dictionary to access the url parameters
    //writer.writeJson({ category: vars["category"], id: vars["id"]})
    return { category: vars["category"], id: vars["id"]} // same as above
  }
}
```

## Getting started

Below demonstrates some features of the Magpie language:

### Basic

```swift
s1 = "hello, 黄"       // strings are UTF-8 encoded
三 = 3                 // UTF-8 identifier
i = 20_000_000         // int
u = 10u                // uint
f = 123_456.789_012    // float
b = true               // bool
a = [1, "2"]           // array
h = {"a": 1, "b": 2}   // hash
t = (1,2,3)            // tuple
n = nil
```

### Const

```swift
const PI = 3.14159
PI = 3.14 //error
```

### Enum

```swift
LogOption = enum {
    Ldate         = 1 << 0,
    Ltime         = 1 << 1,
    Lmicroseconds = 1 << 2,
    Llongfile     = 1 << 3,
    Lshortfile    = 1 << 4,
    LUTC          = 1 << 5,
    LstdFlags     = 1 << 4 | 1 << 5
}

opt = LogOption.LstdFlags
println(opt)
```

### Control Flow

* if
* for
* while
* do
* case-in

#### if

```swift
//if
let a, b = 10, 5
if (a > b) {
    println("a > b")
}
elseif a == b {
    println("a = b")
}
else {
    println("a < b")
}

if 10.isEven() {
    println("10 is even")
}

if 9.isOdd() {
    println("9 is odd")
}
```

#### for

```swift
i = 9
for { // forever loop
    i = i + 2
    if (i > 20) { break }
    println('i = {i}')
}

i = 0
for (i = 0; i < 5; i++) {  // c-like for, '()' is a must
    if (i > 4) { break }
    if (i == 2) { continue }
    println('i is {i}')
}

# for x in arr <where expr> {}
let a = [1,2,3,4]
for i in a where i % 2 != 0 {
    println(i)
}
```

#### while

```swift
i = 10
while (i>3) {
    i--
    println('i={i}')
}
```

#### do

```swift
i = 10
do {
    i--
    if (i==3) { break }
}
```

#### case-in

```swift
let i = [{"a": 1, "b": 2}, 10]
let x = [{"a": 1, "b": 2},10]
case i in {
    1, 2 { println("i matched 1, 2") }
    3    { println("i matched 3") }
    x    { println("i matched x") }
    else { println("i not matched anything")}
}
```

### Array

```swift
a = [1,2,3,4]
for i in a where i % 2 != 0 {
    println(i)
}

if ([].empty()) {
    println("array is empty")
}

a.push(5)

revArr = reverse(a)
println("Reversed Array = ", revArr)
```

### Hash

```swift
hashObj = {
    12     : "twelve",
    true   : 1,
    "Name" : "HHF"
}
println(hashObj)

hashObj += {"key1" : "value1"}
hashObj -= "key1"
hashObj.push(15, "fifteen") //first parameter is the key, second is the value

hs = {"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7}
for k, v in hs where v % 2 == 0 {
    println('{k} : {v}')
}

doc = {
    "one": {
        "two":  { "three": [1, 2, 3], "six":(1,2,3)},
        "four": { "five":  [11, 22, 33]},
    },
}
doc["one"]["two"]["three"][2] = 44
printf("doc[one][two][three][2]=%v\n", doc["one"]["two"]["three"][2])

doc.one.four.five = 4
printf("doc.one.four.five=%v\n", doc.one.four.five)
```

### Tuple

```swift
t = () //same as 't = tuple()'

for i in (1,2,3) {
    println(i)
}
```

### Function

* Default value
* Variadic parameters
* Mutiple return values

```swift
//Function with default values and variadic parameters
add = fn(x, y=5, z=7, args...) {
    w = x + y + z
    for i in args {
        w += i
    }
    return w
}
w = add(2,3,4,5,6,7)
println(w)

let z = (x,y) => x * y + 5
println(z(3,4)) //result :17

# multiple returns
fn testReturn(a, b, c, d=40) {
    return a, b, c, d
}
let (x, y, c, d) = testReturn(10, 20, 30) // d is nil
```

### Command Execution

You could use backtick for command execution(like Perl).

```swift
if (RUNTIME_OS == "linux") {
    var = "~"
    out = `ls -la $var`
    println(out)
}
elseif (RUNTIME_OS == "windows") {
    out = `dir`
    println(out)

    println("")
    println("")
    //test command not exists
    out = `dirs`
    if (!out.ok) {
        printf("Error: %s\n", out)
    }
}
```

### async/await processing
Magpie support `async/await`.

```csharp
let add = async fn(a, b) { a + b }

result = await add(3, 4)
println(result)
```

### Class

* Simple
* Inheritance
* Operator overloading
* Property(like c#)
* Indexer

#### Simple

```swift
class Animal {
    let name = ""
    fn init(name) {    //'init' is the constructor
        //do somthing
    }
}
```

#### Inheritance

```swift
class Dog : Animal { //Dog inherits from Animal
}
```

#### Operator overloading

```swift
class Vector {
    let x = 0;
    let y = 0;

    fn init (a, b) {
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
v1 = new Vector(1,2);
v2 = new Vector(4,5);

v3 = v1 + v2
println(v3.String());

v4 = v1 + 10
println(v4.String());
```

#### Property(like c#)

```swift
class Date {
    let month = 7;  // Backing store
    property Month
    {
        get { return month }
        set {
            if ((value > 0) && (value < 13))
            {
                month = value
            } else {
               println("BAD, month is invalid")
            }
        }
    }

    property Year; // same as 'property Year { get; set;}'
    property Day { get; }

    fn init(year, month, day) {
        this.Year = year
        this.Month = month
        this.Day = day
    }

    fn getDateInfo() {
        printf("Year:%v, Month:%v, Day:%v\n", this.Year, this.Month, this.Day)
    }
}

dateObj = new Date(2000, 5, 11)
dateObj.getDateInfo()
```

#### Indexer

```swift
class IndexedNames
{
    let namelist = []
    let size = 10
    fn init()
    {
        let i = 0
        for (i = 0; i < size; i++)
        {
            namelist[i] = "N. A."
        }
    }

    fn getNameList() {
        println(namelist)
    }

    property this[index] //index must be property
    {
        get
        {
            let tmp;
            if ( index >= 0 && index <= size - 1 )
            {
               tmp = namelist[index]
            }
            else
            {
               tmp = ""
            }
            return tmp
         }
         set
         {
             if ( index >= 0 && index <= size-1 )
             {
                 namelist[index] = value
             }
         }
    }
}
namesObj = new IndexedNames()

//Below code will call Indexer's setter function
namesObj[0] = "Zara"
namesObj[1] = "Riz"
namesObj[2] = "Nuha"
namesObj[3] = "Asif"
namesObj[4] = "Davinder"
namesObj[5] = "Sunil"
namesObj[6] = "Rubic"

namesObj.getNameList()
for (i = 0; i < namesObj.size; i++)
{
    println(namesObj[i]) //Calling Indexer's getter function
}
```

### Standard input/output/error

There are three predefined object for representing standard input, standard output, standard error.
They are `stdin`, `stdout`, `stderr`.

```swift
stdout.writeLine("Hello world")
//same as above
fmt.fprintf(stdout, "Hello world\n")

print("Please type your name:")
name = stdin.read(1024)  //read up to 1024 bytes from stdin
println("Your name is " + name)

//You can also using Insertion operator (`<<`) and Extraction operator(`>>`)
//just like c++ to operate stdin/stdout/stderr.
stdout << "hello " << "world!" << " How are you?" << endl;
```

### Exception Handling(try-catch-finally)

```swift
// Note: Only support throw string type
let exceptStr = "SUMERROR"
try {
    let th = 1 + 2
    if (th == 3) { throw exceptStr }
}
catch "OTHERERROR" {
    println("Catched OTHERERROR")
}
catch exceptStr {
    println("Catched is SUMERROR")
}
catch {
    println("Catched ALL")
}
finally {
    println("finally running")
}
```

### Optional Type(Java 8 like)

```swift
fn safeDivision?(a, b) {
    if (b == 0){
        return optional.empty();
    } else {
        return optional.of(a/b);
    }
}

op1 = safeDivision?(10, 0)
if (!op1.isPresent()) {
    println(op1)
}
```

### Regular expression

```swift
//literal: /pattern/.match(str)
let regex = /\d+\t/.match("abc 123 mnj")
if (regex) {
    println("regex matched using regular expression literal")
}

//Use '=~'(str =~ pattern)
if "abc 123	mnj" =~ ``\d+\t`` {
    println("regex matched using '=~'")
}else {
    println("regex not matched using '=~'")
}
```

### Pipe Operator

```swift
// Test pipe operator(|>)
x = ["hello", "world"] |> strings.join(" ") |> strings.upper() |> strings.lower() |> strings.title()
printf("x=<%s>\n", x)
```

### linq

```swift
let mm = [1,2,3,4,5,6,7,8,9,10]
result = linq.from(mm).where(x => x % 2 == 0).select(x => x + 2).toSlice()

println('before mm={mm}')
println('after result={result}')
```

### json module(for json marshal & unmarshal)

```swift
let hsJson = {"key1" : 10,
              "key2" : "Hello Json %s %s Module",
              "key3" : 15.8912,
              "key4" : [1,2,3.5, "Hello"],
              "key5" : true,
              "key6" : {"subkey1": 12, "subkey2": "Json"},
              "key7" : fn(x,y){x+y}(1,2)
}
let hashStr = json.marshal(hsJson) //same as 'json.toJson(hsJson)'
println(json.indent(hashStr, "  "))

let hsJson1 = json.unmarshal(hashStr)
println(hsJson1)


let arrJson = [1,2.3,"HHF",[],{ "key" : 10, "key1" : 11}]
let arrStr = json.marshal(arrJson)
println(json.indent(arrStr))
let arr1Json = json.unmarshal(arrStr)  //same as 'json.fromJson(arrStr)'
println(arr1Json)
```

### User Defined Operator

```swift
//infix operator '=@' which accept two parameters.
fn =@(x, y) {
    return x + y * y
}
let pp = 10 =@ 5 // Use the '=@' user defined infix operator
printf("pp=%d\n", pp) // result: pp=35


//prefix operator '=^' which accept only one parameter.
fn =^(x) {
    return -x
}
let hh = =^10 // Use the '=^' prefix operator
printf("hh=%d\n", hh) // result: hh=-10
```

### using statement(C# like)

```swift
// No need for calling infile.close().
using (infile = newFile("./file.demo", "r")) {
    if (infile == nil) {
        println("opening 'file.demo' for reading failed, error:", infile.message())
        os.exit(1)
    }

    let line;
    let num = 0
    //Read file by using extraction operator(">>")
    while (infile>>line != nil) {
        num++
        printf("%d %s\n", num, line)
    }
}
```

## Contributing

Contributing is very welcomed. If you make any changes to the language, please let me know,
so i could put you in the `Credits` sections.

## Credits

* mayoms:
    This project is based on mayoms's [monkey](https://github.com/mayoms/monkey) interpreter.

* ahmetb：
    Linq module is base on ahmetb's [linq](https://github.com/ahmetb/go-linq)

* shopspring：
   Decimal module is based on shopspring's [decimal](https://github.com/shopspring/decimal)

* gorilla:
   Service module is based on gorilla's [mux](https://github.com/gorilla/mux)

## Installation

Just download the repository and run `./run.sh`

## License

MIT
