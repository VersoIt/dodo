package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.Info("analytics service started!")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
