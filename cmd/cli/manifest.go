package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func handleManifest(args []string) {
	if len(args) < 1 {
		fmt.Println("expected subcommand: up/down/ps")
		os.Exit(1)
	}
	cmd := args[0]
	flags := parseFlags(args[1:], []string{"-f", "--url", "--token"})
	file, ok := flags["-f"]
	if !ok {
		fmt.Println("-f flag is required")
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

	switch cmd {
	case "up":
		manifestData, err := os.ReadFile(file)
		checkErr(err)

		req, _ := http.NewRequest("POST", url+"/api/v1/manifest/up", bytes.NewReader(manifestData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/x-yaml")
		resp := doRequest(req)
		fmt.Println("Manifest uploaded", resp.Status)

	case "down":
		body, _ := json.Marshal(map[string]string{"manifest": file})
		req, _ := http.NewRequest("POST", url+"/api/v1/manifest/down", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp := doRequest(req)
		fmt.Println("Manifest down", resp.Status)

	case "ps":
		body, _ := json.Marshal(map[string]string{"manifest": file})
		req, _ := http.NewRequest("POST", url+"/api/v1/manifest/ps", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp := doRequest(req)
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))

	default:
		fmt.Println("unknown manifest subcommand")
	}
}
