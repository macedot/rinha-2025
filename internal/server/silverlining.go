package server

import (
	"log"
	"net"

	"rinha-2025-go/internal/config"
	"rinha-2025-go/internal/store"
	"rinha-2025-go/pkg/util"

	"github.com/go-www/silverlining"
)

func ServeListener(ln net.Listener, handler silverlining.Handler) error {
	srv := &silverlining.Server{
		Listener: ln,
		Handler:  handler,
	}
	return srv.Serve(ln)
}

func RunSilverlining(cfg *config.Config, db *store.PaymentDB) error {
	handler := func(c *silverlining.Context) {
		path, method := string(c.Path()), c.Method()
		switch path {
		case "/payments":
			if method != silverlining.MethodPOST {
				c.WriteFullBodyString(405, "Method not allowed")
				return
			}
			// var payment models.Payment
			// if err := c.ReadJSON(&payment); err != nil {
			// 	c.WriteFullBodyString(400, "Invalid JSON payload")
			// 	return
			// }
			// services.EnqueuePayment(&payment, queue)
			c.WriteFullBodyString(200, "OK")
			return

		case "/payments-summary":
			if method != silverlining.MethodGET {
				c.WriteFullBodyString(405, "Method not allowed")
				return
			}
			// var request models.SummaryRequest
			// if err := c.BindQuery(&request); err != nil {
			// 	c.WriteJSONIndent(400, map[string]string{"error": err.Error()}, "", "  ")
			// 	return
			// }
			// response, err := services.GetSummary(&request)
			// if err != nil {
			// 	c.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
			// 	return
			// }
			// c.WriteJSON(200, response)
			return

		case "/purge-payments":
			if method != silverlining.MethodPOST {
				c.WriteFullBodyString(405, "Method not allowed")
				return
			}
			// if err := services.PurgePayments(); err != nil {
			// 	c.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
			// }
			c.WriteFullBodyString(200, "OK")
			return

		default:
			c.WriteFullBodyString(404, "Not found")
		}
	}

	log.Printf("Listening on %s", cfg.ServerSocket)
	listener := util.NewListenUnix(cfg.ServerSocket)
	return ServeListener(listener, handler)
}
