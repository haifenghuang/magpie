<p align="center">
    <img alt="abs language logo" src="https://github.com/haifenghuang/magpie/blob/master/magpie.png?raw=true" width="310">
</p>

# Magpie程序语言

English version: [English](README.md)

## 概述

Magpie是一个用go语言写的解析器. 语法借鉴了C, Ruby, Python, Perl和C#.
支持常用的控制流程，函数式编程和面向对象编程，也能够导入go语言的模块。
它自带一个文档生成工具(mdoc)，用来从magpie源码中生成文档(html).
同时它还包括一个实时语法高亮的REPL。

## 文档

完整的语言教程：[docs](docs)

## 特征

* 类:支持属性，索引器和操作符重载
* 支持await/async的异步编程
* 内置linq支持
* 一级函数(First class function)
* 多参函数及缺省参数值函数
* 函数可以有多个返回值
* int, uint, float, bool, array, tuple, hash(所有均支持json序列化/反序列化, 所有类型均可以扩展)
* try-catch-finally异常处理
* 可选类型支持(类似Java 8的Optional类)
* using语句(类似C#的using)
* Elixir的管道操作符(pipe operator)
* 使用Go语言的方法(RegisterFunctions与RegisterVars)
* 语法高亮REPL
* 文档自动生成工具`mdoc`
* 集成服务(service)处理

## 举例1(Linq)

```csharp
// async/await
async fn add(a, b) { a + b }

result = await add(1, 2)
println(result)

// linq example
class Linq {
    static fn TestSimpleLinq() {
        //数据源
        let ingredients = [
            {Name: "Sugar",  Calories: 500},
            {Name: "Egg",    Calories: 100},
            {Name: "Milk",   Calories: 150},
            {Name: "Flour",  Calories: 50},
            {Name: "Butter", Calories: 200},
        ]

        //检索数据源
        ingredient = from i in ingredients where i.Calories >= 150 orderby i.Name select i

        //显示
        for item in ingredient => println(item)
    }

    static fn TestFileLinq() {
        //从文件读取数据源
        file = newFile("./examples/linqSample.csv", "r")

        //检索数据源
        result = from field in file where int(field[1]) > 300000 select field[0]

        //显示
        for item in result => printf("item = %s\n", item)

        file.close()
    }

    static fn TestComplexLinq() {
        //数据源
        stringList = [
            "A penny saved is a penny earned.",
            "The early bird catches the worm.",
            "The pen is mightier than the sword."
        ]

        //检索数据源
        earlyBirdQuery =
            from sentence in stringList
            let words = sentence.split(" ")
            from word in words
            let w = word.lower()
            where w[0] == "a" || w[0] == "e" ||
                  w[0] == "i" || w[0] == "o" ||
                  w[0] == "u"
            select word

        //显示
        for v in earlyBirdQuery => printf("'%s' starts with a vowel\n", v)
    }
}

Linq.TestSimpleLinq()
println("======================================")
Linq.TestFileLinq()
println("======================================")
Linq.TestComplexLinq()
```

## 举例2(Rest Service)

```csharp
//service Hello on "0.0.0.0:8090" {
service Hello on "0.0.0.0:8090:debug" { //':debug': for debugging request
  //'@route'中，你可以使用'url(必须), methods, host, schemes, headers, queries'
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


## 入门

下面演示了Magpie语言的一些功能:

### 基本

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
n = nil
```

### 常量

```csharp
const PI = 3.14159
PI = 3.14 //错误
```

### 枚举

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

### Case-in语句

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

### 控制流程

* if
* for
* while
* do
* case-in

#### if

```csharp
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

```swfit
i = 9
for { // 无限循环
    i = i + 2
    if (i > 20) { break }
    println('i = {i}')
}

i = 0
for (i = 0; i < 5; i++) {  // 类似C语音的for循环，这里括号'()'是必须的
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

```csharp
i = 10
while (i>3) {
    i--
    println('i={i}')
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

### 数组

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

### 哈希

```csharp
hashObj = {
    12     : "twelve",
    true   : 1,
    "Name" : "HHF"
}
println(hashObj)

hashObj += {"key1" : "value1"}
hashObj -= "key1"
hashObj.push(15, "fifteen") //第一个参数是键, 第二个参数是值

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

### 元祖

```csharp
t = () //等价于't = tuple()'

for i in (1,2,3) {
    println(i)
}
```

### 转换

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

### 函数

```csharp
//带缺省值和多参数
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
println(z(3,4)) //结果 :17

# multiple returns
fn testReturn(a, b, c, d=40) {
    return a, b, c, d
}
let (x, y, c, d) = testReturn(10, 20, 30) // d为nil
```

### 命令执行

你可以使用反引号来执行命令(类似Perl)

```csharp
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
    //下面的代码测试执行失败的情况
    out = `dirs`
    if (!out.ok) {
        printf("Error: %s\n", out)
    }
}
```

### async/await异步处理
Magpie支持`async/await`。

```csharp
let add = async fn(a, b) { a + b }

result = await add(3, 4)
println(result)
```

### 类

* 简单
* 继承
* 操作符重载
* 属性(类似c#)
* 索引器

#### 简单

```csharp
class Animal {
    let name = ""
    fn init(name) {    //'init'是构造方法
        //do somthing
    }
}
```

#### 继承

```csharp
class Dog : Animal { //Dog继承于Animal
}
```

#### 操作符重载

```csharp
class Vector {
    let x = 0;
    let y = 0;

    fn init (a, b) {
        x = a; y = b
    }

    fn +(v) { //重载'+'运算符
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

#### 属性(类似c#)

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

    property Year; // 等价于'property Year { get; set;}'
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

#### 索引器

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

    property this[index] //索引器必须是属性
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

//调用索引器的setter方法
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
    println(namesObj[i]) //调用索引器的getter方法
}
```

### 标准输入/输出/错误

Magpie中预定义了下面三个对象: `stdin`, `stdout`, `stderr`。分别代表标准输入，标准输出，标准错误。

```csharp
stdout.writeLine("Hello world")
//和上面效果一样
fmt.fprintf(stdout, "Hello world\n")

print("Please type your name:")
name = stdin.read(1024)  //从标准输入读最多1024字节
println("Your name is " + name)

//你还可以使用类似C++的插入操作符(`<<`)和提取操作符(`>>`)来操作标准输入和输出。
stdout << "hello " << "world!" << " How are you?" << endl;
```

### 异常处理(try-catch-finally)

```csharp
// 注: 仅支持抛出字符串类型的异常
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

### 可选类型(类似Java 8)

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

### 正则表达式

```csharp
//字面量: /pattern/.match(str)
let regex = /\d+\t/.match("abc 123 mnj")
if (regex) {
    println("regex matched using regular expression literal")
}

//使用 '=~'(str =~ pattern)
if "abc 123	mnj" =~ ``\d+\t`` {
    println("regex matched using '=~'")
}else {
    println("regex not matched using '=~'")
}
```

### 管道操作符

```csharp
// 管道操作符(|>)
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

### json模块(序列化与反序列化)

```csharp
let hsJson = {"key1" : 10,
              "key2" : "Hello Json %s %s Module",
              "key3" : 15.8912,
              "key4" : [1,2,3.5, "Hello"],
              "key5" : true,
              "key6" : {"subkey1": 12, "subkey2": "Json"},
              "key7" : fn(x,y){x+y}(1,2)
}
let hashStr = json.marshal(hsJson) //等价于'json.toJson(hsJson)'
println(json.indent(hashStr, "  "))

let hsJson1 = json.unmarshal(hashStr)
println(hsJson1)


let arrJson = [1,2.3,"HHF",[],{ "key" : 10, "key1" : 11}]
let arrStr = json.marshal(arrJson)
println(json.indent(arrStr))
let arr1Json = json.unmarshal(arrStr)  //等价于'json.fromJson(arrStr)'
println(arr1Json)
```

### 用户自定义操作符

```csharp
//中缀运算符 '=@'接受两个参数
fn =@(x, y) {
    return x + y * y
}
let pp = 10 =@ 5 // 使用刚才定义的'=@'中缀运算符
printf("pp=%d\n", pp) // 结果: pp=35


//前缀运算符'=^'仅接受一个参数
fn =^(x) {
    return -x
}
let hh = =^10 // 使用'=^' 前缀运算符
printf("hh=%d\n", hh) // 结果: hh=-10
```

### using语句(类似C#)

```csharp
// 这里不需要调用infile.close()
using (infile = newFile("./file.demo", "r")) {
    if (infile == nil) {
        println("opening 'file.demo' for reading failed, error:", infile.message())
        os.exit(1)
    }

    let line;
    let num = 0
    //使用提取运算符(">>")读取文件
    while (infile>>line != nil) {
        num++
        printf("%d %s\n", num, line)
    }
}
```

## 贡献

非常欢迎贡献代码。如果您对该语言进行任何更改，请通知我，我会将你放在`Credits`部分。

## 感谢

* mayoms：
    本项目基于mayoms的[monkey](https://github.com/mayoms/monkey)解析器。

* ahmetb：
    Linq模块基于ahmetb的[linq](https://github.com/ahmetb/go-linq)。

* shopspring：
   Decimal模块基于shopspring的[decimal](https://github.com/shopspring/decimal)。

* gorilla:
   Service模块基于gorilla的[mux](https://github.com/gorilla/mux)。

## 安装

下载此仓库并运行`./run.sh`

## 许可证

MIT
