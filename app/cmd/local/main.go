package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kohge4/dynamic-img-gen-cdk/app/internal/router"
)

func main() {
	v1Router := router.New()
	innerRouter := router.NewInnerRouter()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: v1Router,
	}
	innerSrv := &http.Server{
		Addr:    ":8081",
		Handler: innerRouter,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go func() {
		if err := innerSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	if err := innerSrv.Shutdown(ctx); err != nil {
		log.Fatal("InnerServer forced to shutdown: ", err)
	}
}
