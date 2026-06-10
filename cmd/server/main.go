package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-api/config"
	db "user-api/db/sqlc"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// ── Logger ───────────────────────────────────────────────────────────────
	logger.Init()
	defer logger.Sync()

	// ── Config ───────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// ── Database ─────────────────────────────────────────────────────────────
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("database ping failed", zap.Error(err))
	}
	logger.Info("database connection established")

	// ── Wire dependencies ─────────────────────────────────────────────────────
	queries := db.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// ── Fiber app ─────────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "go-users-api",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Return structured JSON on Fiber errors.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if ok := false; !ok {
				_ = e
			}
			if ferr, ok := err.(*fiber.Error); ok {
				code = ferr.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Register(app, userHandler)

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := ":" + cfg.ServerPort
		logger.Info("starting server", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down server …")
	if err := app.Shutdown(); err != nil {
		logger.Error("error during shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}
