# Document generator

 Included also has a tool\(`mdoc`\) for generating documentation in markdown format or html format.

 The tool only support below statement for document generator:

* let statement
* enum statement
* function statement
* class statement
  * let statement
  * function statement
  * property statement

```text
//generate markdown file, the generated file is named 'doc.md'
./mdoc examples/doc.my

//generate html file, the generated file is named 'doc.html'
./mdoc -html examples/doc.my

//generate html file, also generate source code of classes and functions. the generated file is named 'doc.html'
./mdoc -html -showsource examples/doc.my

//Use the some builtin css types for styling the generated html
//    0 - GitHub
//    1 - Zenburn
//    2 - Lake
//    3 - Sea Side
//    4 - Kimbie Light
//    5 - Light Blue
//    6 - Atom Dark
//    7 - Forgotten Light

./mdoc -html -showsource -css 1 examples/doc.my

//Using external css file for styling the generated html file.
//The '-cssfile' option has higher priority than the '-css' option.
//If the supplied css file does not exists, then the '-css' option will be used.
./mdoc -html -showsource -css 1 -cssfile ./examples/github-markdown.css examples/doc.my

//processing all the '.my' files in examples directory, generate html.
./mdoc -html examples
```

 The generating of HTML document is base on github REST API，so you must have network connection to make it work. You may also need to set proxy if you behind a firewall\(Environ variable:HTTP\_PROXY\).

