#!/usr/bin/env bash
export GOPATH=$(pwd)

# format each go file
#echo "Formatting go file..."
#for file in `find ./src/magpie -name "*.go"`; do
#	echo "    `basename $file`"
#	go fmt $file > /dev/null
#done

interpreter_name=magpie

# cross-compiling
platforms=("windows/amd64" "linux/amd64" "darwin/amd64")
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$interpreter_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo "Building mdoc...(mdoc)"
go build -o mdoc mdoc.go

# run: ./fmt demo.my
echo "Building Formatter...(fmt)"
go build -o fmt fmt.go

# run:    ./highlight demo.my               (generate: demo.my.html)
#     or  ./fmt demo.my | ./highlight   (generate: output.html)
echo "Building Highlighter...(highlight)"
go build -o highlight highlight.go
