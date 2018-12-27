# Overview

 This project is based on mayoms's project [monkey](https://github.com/mayoms/monkey) with some bug fixes and a lot of new features including:

* Added simple class support\(Indexer, operator overloading, property, static method/property/field and class annotation\)
* Modified string module\(which can correctly handle utf8 character encoding\)
* Added `file` module\(with some new methods\).
* Added `math`，`time`, `sort`, `os`, `log`, `net`, `http`, `filepath`, `fmt`, `sync`, `list`, `csv`, `regexp`, `template`, etc...
* `sql`\(db\) module\(which can correctly handing null values\)
* `flag` module\(for handling command line options\)
* `json` module\(for json marshaling and unmarshaling\)
* `linq` module\(Code come from [linq](https://github.com/ahmetb/go-linq) with some modifications\)
* `decimal` module\(Code come from [decimal](https://github.com/shopspring/decimal) with some minor modifications\)
* Regular expression literal support\(partially like perls\)
* channel support\(like golang's channel\)
* more operator support\(&&, \|\|, &, \|, ^, +=, -=, ?:, ??, etc.\)
* utf-8 support\(e.g. you could use utf8 character as variable name\)
* more flow control support\(e.g. try/catch/finally, for-in, case, c-like for loop\)
* `defer` support
* \`spawn\` support\(goroutine\)
* `enum` support
* `using` support\(like C\#'s `using`\)
* pipe operator support\(see demo for help\)
* function with default value and variadic parameters
* list comprehension and hash comprehension support
* user defined operator support
* Extending basic type with something like 'int.xxx\(params\)'
* Optional object support
* Using method of Go Package\(`RegisterFunctions` and `RegisterVars`\)

 There are a number of tasks to complete, as well as a number of bugs. The purpose of this project was to dive deeper into Go, as well as get a better understanding of how programming languages work.

 It has been successful in those goals. There may or may not be continued work - I do plan on untangling a few messy spots, and there are a few features I'd like to see implemented. This will happen as time and interest allows.

