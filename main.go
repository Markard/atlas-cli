package main

import (
	"atlas-cli/cmd"
	"atlas-cli/pkg/logger"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log/slog"
	"os"
)

func main() {
	errDotEnvLoad := godotenv.Load()
	if errDotEnvLoad != nil {
		fmt.Println(errDotEnvLoad)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	slog.SetDefault(setupPrettySlog(logFile))

	errExec := cmd.Execute()
	if errExec != nil {
		os.Exit(1)
	}
}

func setupPrettySlog(logFile *os.File) *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	handler := opts.NewPrettyHandler(mw)

	return slog.New(handler)
}
