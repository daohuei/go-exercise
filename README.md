# go-exercise

## Usage

```
Usage of ./bin/go-practice:
  -option int
        0: create empty db file
        1: delete db file
        2: run the main function (default 2)
```

## Build

```bash
go build -o ./bin/go-practice ./cmd
```

## Unit Test

```bash
go test -v  ./pkg/testdb # or ./...
```
