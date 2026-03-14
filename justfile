default:
    go build -o moni .

build:
    go build -o moni .

test:
    go test ./...

lint:
    go vet ./...

run *args:
    go run . {{args}}
