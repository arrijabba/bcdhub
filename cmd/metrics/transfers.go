package main

import (
	"encoding/json"
	"fmt"

	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/logger"
	"github.com/baking-bad/bcdhub/internal/metrics"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/streadway/amqp"
)

func getTransfer(data amqp.Delivery) error {
	var transferID string
	if err := json.Unmarshal(data.Body, &transferID); err != nil {
		return fmt.Errorf("[getTransfer] Unmarshal message body error: %s", err)
	}

	transfer := models.Transfer{ID: transferID}
	if err := ctx.ES.GetByID(&transfer); err != nil {
		return fmt.Errorf("[getTransfer] Find transfer error: %s", err)
	}

	if err := parseTransfer(transfer); err != nil {
		return fmt.Errorf("[getTransfer] Compute error message: %s", err)
	}

	return nil
}

func parseTransfer(transfer models.Transfer) error {
	h := metrics.New(ctx.ES, ctx.DB)
	if h.SetTransferAliases(ctx.Aliases, &transfer) {
		if _, err := ctx.ES.UpdateDoc(elastic.DocTransfers, transfer.ID, transfer); err != nil {
			return err
		}
	}

	logger.Info("Transfer %s processed", transfer.ID)
	return nil
}
