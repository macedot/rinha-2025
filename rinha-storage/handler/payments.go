package handler

import (
	"encoding/json"
	"math"
	"rinha-storage/api"
	"rinha-storage/store"
	"time"

	"github.com/valyala/fasthttp"
)

func PaymentHandler(memoryDB *store.PaymentDB) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var request api.PaymentRequest
		if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
			ctx.Error("Invalid request", fasthttp.StatusBadRequest)
			return
		}

		requestedAt, err := time.Parse(time.RFC3339, request.RequestedAt)
		if err != nil {
			ctx.Error("Invalid requestedAt format. Use RFC3339 format", fasthttp.StatusBadRequest)
			return
		}

		if request.ServerType != "default" && request.ServerType != "fallback" {
			ctx.Error("Invalid server type. Must be 'default' or 'fallback'", fasthttp.StatusBadRequest)
			return
		}

		record := &store.PaymentRecord{
			RequestedAt: requestedAt,
			Amount:      int64(math.Round(request.Amount * 100)),
			ServerType:  request.ServerType,
		}

		memoryDB.AddRecord(record)
		// go fileDB.SaveRecord(record)

		ctx.SetStatusCode(fasthttp.StatusCreated)
	}
}
