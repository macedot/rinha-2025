package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"rinha-2025-go/internal/config"
	"rinha-2025-go/internal/server"
	"rinha-2025-go/internal/store"
)

func main() {
	runtime.GOMAXPROCS(1)
	appCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg := config.GetConfig()
	paymentDB := store.NewPaymentDB(cfg, appCtx)
	defer paymentDB.Close()
	log.Fatalln(server.RunSilverlining(cfg, paymentDB).Error())
}
