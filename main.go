package main

import (
	"encoding/json"
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

	"github.com/valyala/fasthttp"
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

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		switch string(path) {
		case "/payments":
			var payment models.Payment
			body := ctx.Request.Body()
			if err := json.Unmarshal(body, &payment); err != nil {
				ctx.Error(err.Error(), fasthttp.StatusBadRequest)
				return
			}
			services.EnqueuePayment(&payment, queue)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return

		case "/payments-summary":
			var request models.SummaryRequest
			request.StartTime = string(ctx.QueryArgs().Peek("from"))
			request.EndTime = string(ctx.QueryArgs().Peek("to"))
			response, err := services.GetSummary(&request)
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
				return
			}
			body, err := json.Marshal(response)
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
				return
			}
			ctx.SetBody(body)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return

		case "/purge-payments":
			if err := services.PurgePayments(); err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
	}

	log.Fatal(fasthttp.ListenAndServe(cfg.ServerURL, requestHandler))
}
