# LogX - ZAP Adapter

Adapter to wrap loggers from Uber ZAP log package (https://github.com/uber-go/zap)

## Install

```shell
go get -u github.com/logx-go/zap-adapter
```

## Usage

```golang
package main

import (
	"github.com/logx-go/contract/pkg/logx"
	"github.com/logx-go/zap-adapter/pkg/zapadapter"
	"go.uber.org/zap"
)

func main() {
	z, _ := zap.NewDevelopment(
		zap.WithCaller(false), // Caller info will be handled by the zapadapter
	)

	defer z.Sync() // flushes buffer, if any

	logger := zapadapter.New(z)

	logSomething(logger)
}

func logSomething(logger logx.Logger) {
	logger.Info("Hello World")
}
```

## Development

### Requirement
- Golang >=1.20
- golangci-lint (https://golangci-lint.run/)

### Tests

```shell
go test ./... -race
```

### Lint

```shell
golangci-lint run
```

## License

MIT License (see [LICENSE](LICENSE) file)
