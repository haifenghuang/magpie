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



println("\n\n=========================")
println("Double backtick is raw string")

let rawString =``abc
def
lmn
``
println(rawString)

