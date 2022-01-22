<p align="center">
    <img alt="magpie language logo" src="https://github.com/haifenghuang/magpie/blob/master/magpie.png?raw=true" width="310">
</p>

# Magpie Programming Language

Chinese version: [中文](README_cn.md)

## Summary

Magpie is a toy language interpreter, written in Go. It has C-style syntax, and is largely
inspired by Ruby, Python, Perl and c#.

It support the normal control flow, functional programming and object oriented programming.
and also can import golang's module.

It has a built-in documentation generator(mdoc) for generating html document from magpie source.

It has a simple debugger which you can experience with it.

It also has a REPL with realtime syntax highlighter.

I also made a simple programming language written using `magpie`.

You can even run most of the `magpie` script in a web browser.

## Documention

Complete language tutorial can be found in [docs](docs)

## Features

* Class with support for property, indexer & operator overloading
* await/async for asynchronous programming
* Builtin support for linq
* Builtin support for datetime literal
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
* Simple debugger
* Simple Macro processing

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

    /* Code from https://docs.microsoft.com/en-us/dotnet/csharp/language-reference/keywords/let-clause */
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

```csharp
s1 = "hello, 黄"       // strings are UTF-8 encoded
三 = 3                 // UTF-8 identifier
i = 20_000_000         // int
u = 10u                // uint
f = 123_456.789_012    // float
b = true               // bool
a = [1, "2"]           // array
h = {"a": 1, "b": 2}   // hash
t = (1,2,3)            // tuple
dt = dt/2018-01-01 12:01:00/  //datetime literal
n = nil
```

### Const

```csharp
const PI = 3.14159
PI = 3.14 //error

const (
    INT,    //default to 0
    DOUBLE,
    STRING
)
let i = INT
```

### Enum

```csharp
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

```csharp
//if
let a, b = 10, 5
if (a > b) {
    println("a > b")
}
elif a == b {
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

```csharp
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

# read line by line
using (f = open("./file.log", "r")) {
    for line in <$f> where line =~ ``magpie`` {
        println(line) //print only lines which match 'magpie'
    }
}
```

#### while

```csharp
i = 10
while (i>3) {
    i--
    println('i={i}')
}

# read line by line
using (f = open("./file.log", "r")) {
    while <$f> {
        println($_) //$_: line read from file
    }
}
```

#### do

```csharp
i = 10
do {
    i--
    if (i==3) { break }
}
```

#### case-in

```csharp
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

```csharp
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

```csharp
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

// same as below
//doc[one][two][three][2] = 44
doc["one"]["two"]["three"][2] = 44
printf("doc[one][two][three][2]=%v\n", doc["one"]["two"]["three"][2])

doc.one.four.five = 4
printf("doc.one.four.five=%v\n", doc.one.four.five)
```

### Tuple

```csharp
t = () //same as 't = tuple()'

for i in (1,2,3) {
    println(i)
}
```

### datetime literal

```csharp
let month = "01"
let dt0 = dt/2018-{month}-01 12:01:00/
println(dt0)

let dt1 = dt/2018-01-01 12:01:00/.addDate(1, 2, 3).add(time.SECOND * 10) //add 1 year, two months, three days and 10 seconds
printf("dt1 = %v\n", dt1)

/* 'datetime literal' + string:
     string support 'YMDhms' where
       Y:Year    M:Month    D:Day
       h:hour    m:minute   s:second

*/
//same result as 'dt1'
let dt2 = dt/2018-01-01 12:01:00/ + "1Y2M3D10s" //add 1 year, two months, three days and 10 seconds
printf("dt2 = %v\n", dt2)
//same resutl as above
//printf("dt2 = %s\n", dt2.toStr()) //use 'toStr()' method to convert datetime to string.

let dt3 = dt/2019-01-01 12:01:00/
//you could also use strftiem() to convert time object to string. below code converts time object to 'yyyy/mm/dd hh:mm:ss'
format = dt3.strftime("%Y/%m/%d %T")
println(dt3.toStr(format))

////////////////////////////////
// time object to timestamp
////////////////////////////////
println(dt3.unix()) //to timestamp(UTC)
println(dt3.unixNano()) //to timestamp(UTC)
println(dt3.unixLocal()) //to timestamp(LOCAL)
println(dt3.unixLocalNano()) //to timestamp(LOCAL)

////////////////////////////////
// timestamp to time object
////////////////////////////////
timestampUTC = dt3.unix()      //to timestamp(UTC)
println(unixTime(timestampUTC)) //timestamp to time

timestampLocal = dt3.unixLocal() //to timestamp(LOCAL)
println(unixTime(timestampLocal)) //timestamp to time

////////////////////////////////
// datetime comparation
////////////////////////////////
//two datetime literals could be compared using '>', '>=', '<', '<=' and '=='
let dt4 = dt/2018-01-01 12:01:00/
let dt5 = dt/2019-01-01 12:01:00/

println(dt4 <= dt5) //returns true
```

### Regular expression

In magpie, regard to regular expression, you could use:

* Regular expression literal
* 'regexp' module
* =&#126; and !&#126; operators(like perl's)

```swift
//Use regular expression literal( /pattern/.match(str) )
let regex = /\d+\t/.match("abc 123	mnj")
if (regex) { println("regex matched using regular expression literal") }

//Use 'regexp' module
if regexp.compile(``\d+\t``).match("abc 123	mnj") {
    println("regex matched using 'regexp' module")
}

//Use '=~'(str =~ pattern)
if "abc 123	mnj" =~ ``\d+\t`` {
    println("regex matched using '=~'")
}else {
    println("regex not matched using '=~'")
}
```

### Conversion

```csharp
// convert to string using str() function
x = str(10) // convert 10 to string  

// convert to int using int() function
x1 = int("10")   // x1 = 10
x2 = +"10" // same as above

y1 = int("0x10") // y1 = 16
y2 = +"0x10" // same as above

// convert to float using float() funciton
x = float("10.2")

// convert to array using array() funciton
x = array("10") // x = ["10"]
y = array((1, 2, 3)) // convert tuple to array

// convert to tuple using tuple() funciton
x = tuple("10") // x = ("10",)
y = tuple([1, 2, 3]) // convert array to tuple

// convert to hash using hash() function
x = hash(["name", "jack", "age", 20]) // array->hash: x = {"name" : "jack", "age" : 20}
y = hash(("name", "jack", "age", 20)) // tuple->hash: x = {"name" : "jack", "age" : 20}

// if the above conversion functions have no arguments, they simply return
// new corresponding types
i = int()   // i = 0
f = float() // f = 0.0
s = str()   // s = ""
h = hash()  // h = {}
a = array() // a = []
t = tuple() // t = ()
```

### Simple Macro Processing

```csharp
#define DEBUG

// only support two below formats:
//    1. #ifdef xxx { body }
//    2. #ifdef xxx { body } #else { body }, here only one '#else' is supported'.
#ifdef DEBUG2
{
    add = fn(x, y) { x + y }
    printf("add = %d\n", add(1, 2))
}
#else
{
    sub = fn(x, y) { x - y }
    printf("sub = %d\n", sub(3, 1))
}

#define TESTING
#ifdef TESTING
{
    add = fn(x, y) { x + y }
    printf("add = %d\n", add(1, 2))
}
```

### Function

* Default value
* Variadic parameters
* Mutiple return values

```csharp
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
let (x, y, c, d) = testReturn(10, 20, 30) // d is 40

//same as above 'let' statement
//x, y, c, d = testReturn(10, 20, 30) // d is 40
```

### Command Execution

You could use backtick for command execution(like Perl).

```csharp
if (RUNTIME_OS == "linux") {
    var = "~"
    out = `ls -la $var`
    println(out)
}
elif (RUNTIME_OS == "windows") {
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

```csharp
class Animal {
    let name = ""
    fn init(name) {    //'init' is the constructor
        //do somthing
    }
}
```

#### Inheritance

```csharp
class Dog : Animal { //Dog inherits from Animal
}
```

#### Operator overloading

```csharp
class Vector {
    let x = 0;
    let y = 0;

    fn init (a, b) {
        x = a; y = b
    }

    fn +(v) { //overloading '+'
        if (type(v) == "INTEGER" {
            return new Vector(x + v, y + v);
        } elif v.is_a(Vector) {
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

```csharp
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

```csharp
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

```csharp
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

```csharp
// Note: Only support throw string type
let exceptStr = "SUMERROR"
try {
    let th = 1 + 2
    if (th == 3) { throw exceptStr }
}
catch e {
    printf("Catched: %s\n", e)
}
finally {
    println("finally running")
}
```

### Optional Type(Java 8 like)

```csharp
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

```csharp
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

```csharp
// Test pipe operator(|>)
x = ["hello", "world"] |> strings.join(" ") |> strings.upper() |> strings.lower() |> strings.title()

//same as above
//x = ["hello", "world"] |> strings.join(" ") |> strings.upper |> strings.lower |> strings.title
printf("x=<%s>\n", x)
```

### linq

```csharp
let mm = [1,2,3,4,5,6,7,8,9,10]
result = linq.from(mm).where(x => x % 2 == 0).select(x => x + 2).toSlice()

println('before mm={mm}')
println('after result={result}')
```

### json module(for json marshal & unmarshal)

```csharp
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

```csharp
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

```csharp
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
