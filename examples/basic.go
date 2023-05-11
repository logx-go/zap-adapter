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
