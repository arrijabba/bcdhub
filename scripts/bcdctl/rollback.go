package main

import (
	"github.com/baking-bad/bcdhub/internal/logger"
	"github.com/baking-bad/bcdhub/internal/rollback"
)

type rollbackCommand struct {
	Level   int64  `short:"l" long:"level" description:"Level to rollback"`
	Network string `short:"n" long:"network" description:"Network"`
}

var rollbackCmd rollbackCommand

// Execute
func (x *rollbackCommand) Execute(_ []string) error {
	state, err := ctx.Blocks.Last(x.Network)
	if err != nil {
		panic(err)
	}

	logger.Warning("Do you want to rollback '%s' from %d to %d? (yes - continue. no - cancel)", state.Network, state.Level, x.Level)
	if !yes() {
		logger.Info("Cancelled")
		return nil
	}

	manager := rollback.NewManager(ctx.Searcher, ctx.Storage, ctx.BigMapDiffs, ctx.Transfers)
	if err = manager.Rollback(ctx.StorageDB.DB, state, x.Level); err != nil {
		return err
	}
	logger.Info("Done")

	return nil
}
