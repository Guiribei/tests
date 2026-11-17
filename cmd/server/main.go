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
	examsDir := env("EXAMS_DIR", ".")
	repositoryURL := env("EXAMS_REPOSITORY_URL", defaultRepository)

	if envBool("SYNC_ON_START", false) {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		if err := source.Sync(ctx, examsDir, repositoryURL); err != nil {
			logger.Warn("não foi possível sincronizar o repositório; usando os arquivos locais", "error", err)
		}
		cancel()
	}

	catalog, err := exam.LoadDir(examsDir)
	if err != nil {
		logger.Error("falha ao carregar provas", "directory", examsDir, "error", err)
		os.Exit(1)
	}

	databasePath := env("DATABASE_PATH", "./data/provas.db")
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
		Catalog:       catalog,
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

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
