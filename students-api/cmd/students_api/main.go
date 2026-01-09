package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/config"
	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/http/handlers/student"
	"github.com/ahmad-mujtaba1996/api_crud_operations_goLang/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to setup database ", err.Error())
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudentsList(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("Server started at ", slog.String("address", cfg.Address))

	//--------------------------------------------------- Graceful shutdown
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT) // Need to study when to use which signal.

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server ", err.Error())
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 seconds to shutdown server gracefully. Because we need to make sure all ongoing requests are completed and server shuts down without holding any ports.
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server ", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
