#!/bin/sh
export GOPATH=$(pwd)

# format each go file
#echo "Formatting go file..."
#for file in `find ./src/magpie -name "*.go"`; do
#	echo "    `basename $file`"
#	go fmt $file > /dev/null
#done

echo ""

# run: ./magpie demo.my or ./magpie
echo "Building REPL...(magpie)"
go build -o magpie main.go

echo "Building REPL(win)...(magpie.exe)"
GOOS=windows go build -o magpie.exe main.go

echo "Building mdoc...(mdoc)"
go build -o mdoc mdoc.go

# run: ./fmt demo.my
echo "Building Formatter...(fmt)"
go build -o fmt fmt.go

# run:    ./highlight demo.my               (generate: demo.my.html)
#     or  ./fmt demo.my | ./highlight   (generate: output.html)
echo "Building Highlighter...(highlight)"
go build -o highlight highlight.go
