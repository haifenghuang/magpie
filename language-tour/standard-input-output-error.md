# Standard input/output/error

 There are three predefined object for representing standard input, standard output, standard error. They are `stdin`, `stdout`, `stderr`.

```swift
stdout.writeLine("Hello world")
//same as above
fmt.fprintf(stdout, "Hello world\n")

print("Please type your name:")
name = stdin.read(1024)  //read up to 1024 bytes from stdin
println("Your name is " + name)
```

 You can also using Insertion operator \(`<<`\) and Extraction operator\(`>>`\) just like c++ to operate stdin/stdout/stderr.

```swift
// Output to stdout by using insertion operator("<<")
// 'endl' is a predefined object, which is "\n".
stdout << "hello " << "world!" << " How are you?" << endl;

// Read from stdin by using extraction operator(">>")
let name;
stdout << "Your name please: ";
stdin >> name;
printf("Welcome, name=%v\n", name)
```

 Insertion operator \(`<<`\) and Extraction operator\(`>>`\) can also be used for operating file object.

```swift
//Read file by using extraction operator(">>")
infile = newFile("./file.demo", "r")
if (infile == nil) {
    println("opening 'file.demo' for reading failed, error:", infile.message())
    os.exit(1)
}
let line;
let num = 0
while ( infile>>line != nil) {
    num++
    printf("%d	%s\n", num, line)
}
infile.close()


//Writing to file by using inserttion operator("<<")
outfile = newFile("./outfile.demo", "w")
if (outfile == nil) {
    println("opening 'outfile.demo' for writing failed, error:", outfile.message())
    os.exit(1)
}
outfile << "Hello" << endl
outfile << "world" << endl
outfile.close()
```

