package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
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
	manifestData, err := os.ReadFile(file)
	checkErr(err)

	switch cmd {
	case "up":
		req, _ := http.NewRequest("POST", url+"/api/v1/manifest/up", bytes.NewReader(manifestData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/x-yaml")
		resp := doRequest(req)
		fmt.Println("Manifest uploaded", resp.Status)

	case "down":
		var parsed struct {
			Name string `yaml:"name"`
		}
		if err := yaml.Unmarshal(manifestData, &parsed); err != nil {
			log.Fatalf("failed to parse YAML: %v", err)
		}
		body, _ := json.Marshal(map[string]string{"manifest": parsed.Name})
		req, _ := http.NewRequest("POST", url+"/api/v1/manifest/down", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp := doRequest(req)
		fmt.Println("Manifest down", resp.Status)

	case "ps":
		var parsed struct {
			Name string `yaml:"name"`
		}
		if err := yaml.Unmarshal(manifestData, &parsed); err != nil {
			log.Fatalf("failed to parse YAML: %v", err)
		}
		body, _ := json.Marshal(map[string]string{"manifest": parsed.Name})
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
