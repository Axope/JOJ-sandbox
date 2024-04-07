GOOS=linux GOARCH=amd64 go build -o ./container/sandbox sandbox.go

output="./container/output"
if [ ! -d $output ]; then
    mkdir $output
fi

data="./container/data"
if [ ! -d $data ]; then
    mkdir $data
fi