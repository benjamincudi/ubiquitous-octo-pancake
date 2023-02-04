package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger())
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{Addr: ":80", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Println("serving traffic on :80")

	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second) // simulate request draining
	log.Println("graceful shutdown complete")
}
