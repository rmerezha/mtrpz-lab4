package main

import (
	"encoding/json"
	"fmt"
	"github.com/rmerezha/mtrpz-lab4/config"
	"net/http"
	"os"
	"strings"
)

func parseFlags(args []string, keys []string) map[string]string {
	flags := make(map[string]string)
	for i := 0; i < len(args); i++ {
		for _, key := range keys {
			if args[i] == key && i+1 < len(args) {
				flags[key] = args[i+1]
				i++
				break
			}
		}
	}
	return flags
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func doRequest(req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	checkErr(err)
	if resp.StatusCode >= 300 {
		fmt.Printf("HTTP error: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	return resp
}

func printContainerListJSON(body []byte) {
	var containers []config.ContainerStatus
	err := json.Unmarshal(body, &containers)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		fmt.Println(string(body))
		return
	}

	fmt.Printf("%-3s  %-10s  %-10s  %-6s  %-15s  %-12s  %-10s\n", "#", "Manifest", "Name", "Host", "Image", "Ports", "State")
	fmt.Println(strings.Repeat("-", 75))

	for i, c := range containers {
		ports := "-"
		if len(c.Config.Ports) > 0 {
			ports = strings.Join(c.Config.Ports, ",")
		}
		fmt.Printf("%-3d  %-10s  %-10s  %-6s  %-15s  %-12s  %-10s\n",
			i+1,
			c.ManifestName,
			c.Config.Name,
			c.Config.Host,
			shorten(c.Config.Image, 15),
			shorten(ports, 12),
			c.State,
		)
	}
}

func shorten(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
