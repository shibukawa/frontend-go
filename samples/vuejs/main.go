package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shibukawa/frontend-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, "hello world")
	})
	http.Handle("/", frontend.MustNewSPAHandler(ctx))

	server := &http.Server{
		Addr: ":8888",
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()
	fmt.Println("start receiving at :8888")
	log.Fatal(server.ListenAndServe())
}
