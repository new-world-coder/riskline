package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/new-world-coder/riskline/internal/httpserver"
	"github.com/new-world-coder/riskline/pkg/config"
	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/ruleset"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	regimesFlag := flag.String("regimes", "", "comma-separated default regime packs (overridable per request)")
	flag.Parse()

	loader, err := ruleset.DefaultLoader()
	if err != nil {
		log.Fatalf("loader: %v", err)
	}

	regs, err := config.ResolveRegimes(config.ParseRegimesFlag(*regimesFlag), ".")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	eng, err := engine.NewWithLoader(loader, regs)
	if err != nil {
		log.Fatalf("engine: %v", err)
	}

	srv := httpserver.New(eng)
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("riskline-api listening on %s regimes=%v (POST /v1/classify)", *addr, regs)
	if err := httpSrv.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
