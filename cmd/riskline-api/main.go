package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/new-world-coder/riskline/internal/httpserver"
	"github.com/new-world-coder/riskline/pkg/engine"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	eng, err := engine.Default()
	if err != nil {
		log.Fatalf("engine: %v", err)
	}

	srv := httpserver.New(eng)
	log.Printf("riskline-api listening on %s (POST /v1/classify)", *addr)
	if err := http.ListenAndServe(*addr, srv.Handler()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
