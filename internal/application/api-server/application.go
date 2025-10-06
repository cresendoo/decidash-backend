package apiserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/cresendoo/decidash-backend/internal/xaptos"
	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/cresendoo/decidash-backend/pkg/xredis"
	"github.com/mediocregopher/radix/v3"
)

type Application struct {
	ctx context.Context

	pool *radix.Pool

	httpServer *http.Server
	logger     *slog.Logger

	aptos   *aptos.Client
	sponsor *aptos.Account
}

func NewApplication(ctx context.Context, logger *slog.Logger, cfg *Config) (*Application, error) {
	app := Application{ctx: ctx, logger: logger}
	var err error

	app.pool, err = xredis.NewRedisPool(cfg.Redis.Addr, cfg.Redis.Pool, cfg.Redis.DB, "")
	if err != nil {
		return nil, err
	}
	app.sponsor, err = xaptos.AccountFromEd25519PrivateKey(cfg.AptosAccounts.FeePayer)
	if err != nil {
		return nil, err
	}
	app.aptos, err = aptos.NewClient(aptos.NetworkConfig{
		Name:    "devnet",
		ChainId: 4,
		NodeUrl: "https://api.netna.staging.aptoslabs.com/v1",
	})
	if err != nil {
		return nil, err
	}

	app.httpServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: app.setRouter(),
	}

	return &app, nil
}

func (a *Application) Start() error {
	go func() {
		slog.Info("listen http server", "port", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); !errorx.Is(err, http.ErrServerClosed) {
			slog.Error("init http serve", "error", err)
			return
		}
		slog.Info("stopped serving new connection")
	}()
	return nil
}

func (a *Application) Close() error {
	slog.Info("start shutdown")
	defer slog.Info("finish shutdown")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(10)*time.Second,
	)
	defer cancel()
	if err := a.httpServer.Shutdown(ctx); err != nil {
		if err != context.DeadlineExceeded {
			return errorx.Wrap(err)
		} else {
			if err := a.httpServer.Close(); err != nil {
				return errorx.Wrap(err)
			}
		}
	}
	return nil
}
