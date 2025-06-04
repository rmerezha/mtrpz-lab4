package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func handleContainer(args []string) {
	if len(args) < 1 {
		fmt.Println("expected subcommand: stop/kill/restart/rm")
		os.Exit(1)
	}
	cmd := args[0]
	flags := parseFlags(args[1:], []string{"-f", "-c", "--url", "--token"})
	host, ok := flags["-f"]
	if !ok {
		fmt.Println("-f flag is required")
		os.Exit(3)
	}
	container, ok := flags["-c"]
	if !ok {
		fmt.Println("-c flag is required")
		os.Exit(3)
	}
	url, ok := flags["--url"]
	if !ok {
		fmt.Println("-url flag is required")
		os.Exit(3)
	}
	token, ok := flags["--token"]
	if !ok {
		fmt.Println("-token flag is required")
		os.Exit(3)
	}

	body, _ := json.Marshal(map[string]string{
		"action":    cmd,
		"host":      host,
		"container": container,
	})

	req, _ := http.NewRequest("POST", url+"/api/v1/container/action", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := doRequest(req)
	fmt.Println("Action", cmd, "status:", resp.Status)
}
