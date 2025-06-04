package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rmerezha/mtrpz-lab4/api"
	"github.com/rmerezha/mtrpz-lab4/auth"
	"github.com/rmerezha/mtrpz-lab4/planner"
)

var (
	port      = flag.String("port", "8080", "Port to run the master server on")
	ip        = flag.String("ip", "0.0.0.0", "IP address to bind the server to")
	tokenFile = flag.String("token-file", "tokens.txt", "Path to the token file")
	tokenPass = flag.String("token-pass", "", "Password required to generate new tokens")
)

func main() {
	flag.Parse()

	if *tokenPass == "" {
		fmt.Fprintln(os.Stderr, "--token-pass must be specified")
		os.Exit(1)
	}

	authManager, err := auth.NewManager(*tokenFile)
	if err != nil {
		log.Fatalf("failed to initialize auth manager: %v", err)
	}
	defer authManager.Close()

	if err := authManager.LoadFromFile(); err != nil {
		log.Fatalf("failed to load tokens from %s: %v", *tokenFile, err)
	}

	pl := planner.NewPlanner()

	mux := http.NewServeMux()
	server := &api.Server{
		Planner:  pl,
		Auth:     authManager,
		Password: *tokenPass,
	}
	server.RegisterRoutes(mux)

	addr := fmt.Sprintf("%s:%s", *ip, *port)
	log.Printf("Master server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
