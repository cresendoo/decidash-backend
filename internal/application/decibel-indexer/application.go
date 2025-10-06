package decibelindexer

import (
	"context"
	"io"
	"log/slog"
	"sync"

	"github.com/cresendoo/decidash-backend/internal/application/decibel-indexer/models"
	"github.com/cresendoo/decidash-backend/pkg/fullnode"
	"github.com/mediocregopher/radix/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	ctx    context.Context
	pool   *radix.Pool
	logger *slog.Logger

	db      *gorm.DB
	fetcher *fullnode.FullnodeFetcher
	stream  *fullnode.FullnodeRpcStream

	wg sync.WaitGroup
}

func NewApplication(ctx context.Context, logger *slog.Logger, cfg *Config) (*Application, error) {
	fetcher, err := fullnode.NewFullnodeRpcClient(
		"https://api.netna.staging.aptoslabs.com/v1",
		"",
	)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.DB), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.IndexerState{}, &models.PerpPosition{})
	if err != nil {
		return nil, err
	}

	return &Application{ctx: ctx, logger: logger, fetcher: fetcher, db: db}, nil
}

func (a *Application) Start() error {
	state, err := models.GetIndexerState(a.db, "decibel-indexer")
	if err != nil {
		return err
	}

	a.stream, err = a.fetcher.NewStream(state.LastProcessedVersion, 100)
	if err != nil {
		return err
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		slog.Info("started streaming transactions", "start", state.LastProcessedVersion)
		for {
			txs, err := a.stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				slog.Error("failed to receive transactions", "error", err)
				return
			}
			if err := a.Process(txs); err != nil {
				slog.Error("failed to process transactions", "error", err)
				return
			}
		}
	}()

	return nil
}

func (a *Application) Close() error {
	if a.stream != nil {
		a.stream.Close()
	}
	a.wg.Wait()
	return nil
}
