package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"rinha-2025/config"
	"rinha-2025/database"
	"rinha-2025/models"
	"rinha-2025/services"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gogearbox/gearbox"
)

func NewListenUnix(socketPath string) net.Listener {
	if socketPath == "" {
		return nil
	}
	socketDir := filepath.Dir(socketPath)
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		log.Fatalf("Failed to create socket directory: %v", err)
	}
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to remove existing socket: %v", err)
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on Unix socket: %v", err)
	}
	if err := os.Chmod(socketPath, 0666); err != nil {
		log.Fatalf("Failed to set socket permissions: %v", err)
	}
	return listener
}

func main() {
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(-1)
	debug.SetMemoryLimit(90 * 1024 * 1024)

	cfg := config.ConfigInstance()
	cfg.Init()

	client := services.HttpClientInstance()
	client.Init()

	redis := database.RedisInstance()
	redis.Connect(cfg)

	//listener := NewListenUnix(cfg.ServerSocket)

	go func() {
		services.ResetHealthTimeout()
		sleep := time.Duration(rand.Intn(3))
		log.Printf("Sleep for %d seconds...", sleep)
		time.Sleep(sleep * time.Second)
		for {
			now := time.Now().UTC().UnixNano()
			expiration := services.GetHealthTimeout()
			if expiration < now {
				services.SetHealthTimeout(time.Hour)
				services.RefreshServiceStatus(cfg)
				services.SetHealthTimeout(cfg.ServiceRefreshInterval)
			}
			time.Sleep(cfg.ServiceRefreshInterval)
		}
	}()

	queue := services.NewQueue()
	go func() {
		for {
			payment := queue.Dequeue()
			if err := services.ProcessPayment(&payment); err != nil {
				log.Println("ProcessPayment:", err.Error())
				queue.Enqueue(&payment)
			}
		}
	}()

	gb := gearbox.New()

	gb.Post("/payments", func(ctx gearbox.Context) {
		var payment models.Payment
		if err := ctx.ParseBody(&payment); err != nil {
			ctx.Status(500).SendString(err.Error())
			return
		}
		services.EnqueuePayment(&payment, queue)
		ctx.SendString("")
	})

	gb.Get("/payments-summary", func(ctx gearbox.Context) {
		var request models.SummaryRequest
		request.StartTime = ctx.Query("from")
		request.EndTime = ctx.Query("to")
		response, err := services.GetSummary(&request)
		if err != nil {
			ctx.Status(500).SendString(err.Error())
			return
		}
		ctx.SendJSON(response)
	})

	gb.Post("/purge-payments", func(ctx gearbox.Context) {
		if err := services.PurgePayments(); err != nil {
			ctx.Status(500).SendString(err.Error())
			return
		}
		runtime.GC()
		ctx.SendString("")
	})

	log.Fatal(gb.Start(cfg.ServerURL))
}
