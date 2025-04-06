package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"htmx-upload/internal/server"

	CL "github.com/xlab/closer"
)

// *http.Server -> ()
func gracefulShutdown(svr *http.Server) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown with error: %v", err)
		}

		log.Println("Server exiting...")
	}
}

func main() {
	svr := server.NewServer()
	CL.Bind(gracefulShutdown(svr))

	log.Println("Starting server...")

	err := svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	CL.Hold()

	log.Println("Graceful shutdown complete.")
}
