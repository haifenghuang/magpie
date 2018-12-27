# Using go language modules

 Magpie has experimental support for working with `go` modules.

If you need to use go's language package function, you first need to use `RegisterFunctions` or `RegisterVars` to register `go` language functions or types into Magpie language.

 Below is an example of `main.go`\(extracted\):

```swift
// Because in magpie we already have built in module `fmt`, 
// here we use `gfmt` for package name.
eval.RegisterFunctions("gfmt", []interface{}{
    fmt.Errorf,
    fmt.Println, fmt.Print, fmt.Printf,
    fmt.Fprint, fmt.Fprint, fmt.Fprintln, fmt.Fscan, fmt.Fscanf, fmt.Fscanln,
    fmt.Scan, fmt.Scanf, fmt.Scanln,
    fmt.Sscan, fmt.Sscanf, fmt.Sscanln,
    fmt.Sprint, fmt.Sprintf, fmt.Sprintln,
})

eval.RegisterFunctions("io/ioutil", []interface{}{
    ioutil.WriteFile, ioutil.ReadFile, ioutil.TempDir, ioutil.TempFile,
    ioutil.ReadAll, ioutil.ReadDir, ioutil.NopCloser,
})

eval.Eval(program, scope)
```

 Now, in your magpie file, you could use it like below:

```swift
gfmt.Printf("Hello %s!\n", "go function");

//Note Here: we use 'io_ioutil', not 'io/ioutil'.
let files, err = io_ioutil.ReadDir(".")
if err != nil {
    gfmt.Println(err)
}
for file in files {
    if file.Name() == ".git" {
        continue
    }
    gfmt.Printf("Name=%s, Size=%d\n", file.Name(), file.Size())
}
```

 For more detailed examples, please see `goObj.my`.

