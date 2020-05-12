# imc
InMemoryCache

## Build

```bash
    go build imc.go
```

## Test

```bash
    go test
    go test --race
```

### Lint

```bash
    docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.26.0 golangci-lint run -v --enable-all
```
