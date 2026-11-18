package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/source"
	"github.com/Guiribei/tests/internal/store"
	webapp "github.com/Guiribei/tests/internal/web"
)

const defaultRepository = "https://github.com/Guiribei/tests.git"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	databasePath := env("DATABASE_PATH", "./data/provas.db")
	examsDir := env("EXAMS_DIR", filepath.Join(filepath.Dir(databasePath), "exams"))
	repositoryURL := env("EXAMS_REPOSITORY_URL", defaultRepository)
	interval, err := time.ParseDuration(env("SYNC_INTERVAL", "1m"))
	if err != nil || interval < 0 {
		logger.Error("SYNC_INTERVAL deve ser uma duração positiva ou zero")
		os.Exit(1)
	}

	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var catalog exam.Catalog
	refresh := func(syncGit bool) {
		if syncGit {
			ctx, cancel := context.WithTimeout(runCtx, 45*time.Second)
			err := source.Sync(ctx, examsDir, repositoryURL)
			cancel()
			if err != nil {
				logger.Warn("falha ao sincronizar Git, consultando a cópia local", "directory", examsDir, "error", err)
			}
		}
		if err := catalog.Reload(examsDir); err != nil {
			logger.Warn("falha ao ler arquivos de prova", "directory", examsDir, "error", err)
		}
	}
	refresh(envBool("SYNC_ON_START", true))

	if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
		logger.Error("falha ao criar diretório do banco", "error", err)
		os.Exit(1)
	}
	db, err := store.Open(databasePath)
	if err != nil {
		logger.Error("falha ao abrir banco de dados", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	handler, err := webapp.New(webapp.Config{
		Catalog:       &catalog,
		Store:         db,
		Logger:        logger,
		RepositoryURL: repositoryURL,
		SecureCookies: envBool("SECURE_COOKIES", false),
	})
	if err != nil {
		logger.Error("falha ao preparar aplicação", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              env("ADDR", ":8080"),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("servidor iniciado", "address", server.Addr, "exams", len(catalog.All()))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("servidor encerrado inesperadamente", "error", err)
			os.Exit(1)
		}
	}()

	if interval > 0 {
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-runCtx.Done():
					return
				case <-ticker.C:
					refresh(true)
				}
			}
		}()
	}
	<-runCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("falha no encerramento gracioso", "error", err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
