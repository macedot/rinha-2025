package main

import (
	"log"
	"math/rand"
	"rinha-2025/config"
	"rinha-2025/database"
	"rinha-2025/models"
	"rinha-2025/services"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/ohler55/ojg/oj"
	"github.com/valyala/fasthttp"
)

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
			if err := oj.Unmarshal(body, &payment); err != nil {
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
			body, err := oj.Marshal(response)
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
