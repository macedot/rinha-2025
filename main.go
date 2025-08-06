package main

import (
	"log"
	"math/rand"
	"rinha-2025/config"
	"rinha-2025/database"
	"rinha-2025/models"
	"rinha-2025/services"
	"runtime"
	"time"

	"github.com/go-www/silverlining"
)

func main() {
	runtime.GOMAXPROCS(1)

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

	handler := func(c *silverlining.Context) {
		switch string(c.Path()) {
		case "/payments":
			if c.Method() != silverlining.MethodPOST {
				c.WriteFullBodyString(405, "Method not allowed")
				return
			}
			var payment models.Payment
			if err := c.ReadJSON(&payment); err != nil {
				c.WriteFullBodyString(400, "Invalid JSON payload")
				return
			}
			services.EnqueuePayment(&payment, queue)
			c.WriteFullBodyString(200, "OK")
			return

		case "/payments-summary":
			if c.Method() != silverlining.MethodGET {
				c.WriteFullBodyString(405, "Method not allowed")
				return
			}
			var request models.SummaryRequest
			if err := c.BindQuery(&request); err != nil {
				c.WriteJSONIndent(400, map[string]string{"error": err.Error()}, "", "  ")
				return
			}
			log.Println("request:", request)
			response, err := services.GetSummary(&request)
			if err != nil {
				c.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
				return
			}
			c.WriteJSON(200, response)
			return

		case "/purge-payments":
			if err := services.PurgePayments(); err != nil {
				c.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
			}
			c.WriteFullBodyString(200, "OK")
			return

		default:
			c.WriteFullBodyString(404, "Not found")
		}
	}

	log.Fatal(silverlining.ListenAndServe(cfg.ServerURL, handler))
}
