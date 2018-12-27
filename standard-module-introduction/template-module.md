# template module

 The `template` module contains 'text' and 'html' template handling.

 Use `newText(...)` or `parseTextFiles(...)` to create a new 'text' template.

 Use `newHtml(...)` or `parseHtmlFiles(...)` to create a new 'html' template.

```swift
arr = [
    { "key" : "key1", "value" : "value1" },
    { "key" : "key2", "value" : "value2" },
    { "key" : "key3", "value" : "value3" }
]

//use parseTextFiles(), write to a string
template.parseTextFiles("./examples/looping.tmpl").execute(resultValue, arr)
println('{resultValue}')

//use parseTextFiles(), write to a file
file = newFile("./examples/outTemplate.log", "a+")
template.parseTextFiles("./examples/looping.tmpl").execute(file, arr)
file.close() //do not to forget to close the file

//use parse()
//Note here: we need to use "{{-" and "-}}" to remove the newline from the output
template.newText("array").parse(`Looping
{{- range . }}
        key={{ .key }}, value={{ .value -}}
{{- end }}
`).execute(resultValue, arr)
println('{resultValue}')
```

