package main

import (
	"context"
	"fmt"
	"github.com/a-klimenko/go-otus-final-project/internal/app"
	"github.com/a-klimenko/go-otus-final-project/internal/logger"
	internalgrpc "github.com/a-klimenko/go-otus-final-project/internal/server/grpc"
	sqlstorage "github.com/a-klimenko/go-otus-final-project/internal/storage/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const RotatorLogFile = "/opt/rotator/logs/rotator.log"
const RotatorLogLevel = "info"

func main() {
	logFile, err := os.OpenFile(RotatorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			log.Fatalf("can not close log file: %v", err)
		}
	}()
	logg := logger.New(RotatorLogLevel, logFile)

	storage := sqlstorage.New()

	err = storage.Connect()
	if err != nil {
		logg.Error(fmt.Sprintf("can not connect to storage: %s", err))
	}

	rotator := app.New(logg, storage)

	grpcServer := internalgrpc.NewServer(logg, rotator)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	doneCh := make(chan struct{})
	go func() {
		<-ctx.Done()

		if err := grpcServer.Stop(); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		doneCh <- struct{}{}
	}()

	logg.Info("rotator is running...")
	go func() {
		logg.Info("starting grpc server...")
		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	<-doneCh
}
