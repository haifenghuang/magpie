# linq module

 In magpie, the `linq` module support seven types of object:

* File object \(create using newFile builtin function\)
* Csv reader object \(created using newCsvReader builtin function\)
* String object
* Array object
* Tuple object
* Hash object
* Channel object \(created using chan builtin function\)

```swift
let mm = [1,2,3,4,5,6,7,8,9,10]
println('before mm={mm}')

result = linq.from(mm).where(fn(x) {
    x % 2 == 0
}).select(fn(x) {
    x = x + 2
}).toSlice()
println('after result={result}')

result = linq.from(mm).where(fn(x) {
    x % 2 == 0
}).select(fn(x) {
    x = x + 2
}).last()
println('after result={result}')

let sortArr = [1,2,3,4,5,6,7,8,9,10]
result = linq.from(sortArr).sort(fn(x,y){
    return x > y
})
println('[1,2,3,4,5,6,7,8,9,10] sort(x>y)={result}')

result = linq.from(sortArr).sort(fn(x,y){
    return x < y
})
println('[1,2,3,4,5,6,7,8,9,10] sort(x<y)={result}')

thenByDescendingArr = [
    {"Owner" : "Google",    "Name" : "Chrome"},
    {"Owner" : "Microsoft", "Name" : "Windows"},
    {"Owner" : "Google",    "Name" : "GMail"},
    {"Owner" : "Microsoft", "Name" : "VisualStudio"},
    {"Owner" : "Google",    "Name" : "GMail"},
    {"Owner" : "Microsoft", "Name" : "XBox"},
    {"Owner" : "Google",    "Name" : "GMail"},
    {"Owner" : "Google",    "Name" : "AppEngine"},
    {"Owner" : "Intel",     "Name" : "ParallelStudio"},
    {"Owner" : "Intel",     "Name" : "VTune"},
    {"Owner" : "Microsoft", "Name" : "Office"},
    {"Owner" : "Intel",     "Name" : "Edison"},
    {"Owner" : "Google",    "Name" : "GMail"},
    {"Owner" : "Microsoft", "Name" : "PowerShell"},
    {"Owner" : "Google",    "Name" : "GMail"},
    {"Owner" : "Google",    "Name" : "GDrive"}
]

result = linq.from(thenByDescendingArr).orderBy(fn(x) {
    return x["Owner"]
}).thenByDescending(fn(x){
    return x["Name"]
}).toOrderedSlice()    //Note: You need to use toOrderedSlice

//use json.indent() for formatting the output
let thenByDescendingArrStr = json.marshal(result)
println(json.indent(thenByDescendingArrStr, "  "))

//test 'selectManyByIndexed'
println()
let selectManyByIndexedArr1 = [[1, 2, 3], [4, 5, 6, 7]]
result = linq.from(selectManyByIndexedArr1).selectManyByIndexed(
fn(idx, x){
    if idx == 0 { return linq.from([10, 20, 30]) }
    return linq.from(x)
}, fn(x,y){
    return x + 1
})
println('[[1, 2, 3], [4, 5, 6, 7]] selectManyByIndexed() = {result}')

let selectManyByIndexedArr2 = ["st", "ng"]
result = linq.from(selectManyByIndexedArr2).selectManyByIndexed(
fn(idx,x){
    if idx == 0 { return linq.from(x + "r") }
    return linq.from("i" + x)
},fn(x,y){
    return x + "_"
})
println('["st", "ng"] selectManyByIndexed() = {result}')
```

## **Linq for file**

 Now, magpie has a powerful `linq for file` support. it can be used to operate files a little bit like awk. See below for example:

```swift
//test: linq for "file"
file = newFile("./examples/linqSample.csv", "r") //open linqSample.csv file for reading
result = linq.from(file,",",fn(line){ //the second parameter is field separator, the third is a selector function
    if line.trim().hasPrefix("#") { //if line start '#'
        return true // return 'true' means we ignore this line
    } else {
        return false
    }
}).where(fn(fields) {
    //The 'fields' is an array of hashes, like below:
    //  fields = [
    //      {"line": LineNo1, "nf": line1's number of fields, 0: line1, 1: field1, 2: field2, ...},
    //      {"line": LineNo2, "nf": line2's number of fields, 0: line2, 1: field1, 2: field2, ...}
    //  ]

    int(fields[1]) > 300000 //only 1st Field's Value > 300000
}).sort(fn(field1,field2){
    return int(field1[1]) > int(field2[1]) //sort with first field(descending)
}).select(fn(fields) {
    fields[5]  //only output the fifth field
})
println(result)
file.close() //do not forget to close the file

//another test: linq for "file"
file = newFile("./examples/linqSample.csv", "r") //open linqSample.csv file for reading
result = linq.from(file,",",fn(line){ //the second parameter is field separator, the third is a selector function
    if line.trim().hasPrefix("#") { //if line start '#'
        return true //return 'true' means we ignore this line
    } else {
        return false
    }
}).where(fn(fields) {
    int(fields[1]) > 300000 //only 1st Field's Value > 300000
}).sort(fn(field1,field2){
    return int(field1[1]) > int(field2[1]) //sort with first field(descending)
}).selectMany(fn(fields) {
    row = [[fields[0]]] //fields[0] is the whole line, we need two "[]"s, otherwise selectMany() will flatten the output.
    linq.from(row)  //output the whole records
})
println(result)
file.close() //do not forget to close the file


//test: linq for "csv"
r = newCsvReader("./examples/test.csv") //open test.csv file for reading
r.setOptions({"Comma": ";", "Comment": "#"})
result = linq.from(r).where(fn(x) {
    //The 'x' is an array of hashes, like below:
    //  x = [
    //      {"nf" : line1's number of fields, 1: field1, 2: field2, ...},
    //      {"nf" : line2's number of fields, 1: field1, 2: field2, ...}
    //  ]
    x[2] == "Pike"//only 2nd Field = "Pike"
}).sort(fn(x,y){
    return len(x[1]) > len(y[1]) //sort with length of first field
})
println(result)
r.close() //do not forget to close the reader
```

