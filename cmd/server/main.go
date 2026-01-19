package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"golang-test-task/internal/config"
	"golang-test-task/internal/httpapi"
	"golang-test-task/internal/storage"
)

func main() {
	// тянем конфиг из env, чтобы сервис можно было крутить без пересборки
	cfg := config.Load()

	// поднимаем пул к постгре, дальше будем жить на нем
	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// быстрый пинг, чтобы не стартовать в пустоту
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// инициализируем стор и схему, чтобы сразу было куда писать
	store := storage.NewPostgres(db)
	if err := store.EnsureSchema(ctx); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	// поднимаем http сервер с роутером
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.NewRouter(store),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// отдельный канал под ошибки сервера, чтобы не терять фейлы
	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	log.Printf("listening on %s", cfg.HTTPAddr)

	// слушаем сигналы, чтобы нормально выключиться
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalCh:
		log.Printf("shutdown on %s", sig)
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	}

	// даем себе окно на аккуратный shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
