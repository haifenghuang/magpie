#!/usr/bin/env bash
export GOPATH=$(pwd)

# for newer go version's build error: "package XXX is not in GOROOT"
export GO111MODULE=off

# format each go file
echo "Formatting go file..."
for file in `find ./src/magpie -name "*.go"`; do
	echo "    `basename $file`"
	go fmt $file > /dev/null
done

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

    echo "Building ${interpreter_name}...       ($output_name)"
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o $output_name main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo "Building mdoc...         (mdoc)"
go build -ldflags "-s -w" -o mdoc mdoc.go

# run: ./fmt demo.mp
echo "Building Formatter...    (fmt)"
go build -ldflags "-s -w" -o fmt fmt.go

# run:    ./highlight demo.mp               (generate: demo.mp.html)
#     or  ./fmt demo.mp | ./highlight   (generate: output.html)
echo "Building Highlighter...  (highlight)"
go build -ldflags "-s -w" -o highlight highlight.go
