package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.Info("orders service started!")

	// Keep the service alive
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("orders service shutting down...")
}
