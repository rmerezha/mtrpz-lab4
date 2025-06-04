package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func handleToken(args []string) {
	flags := parseFlags(args, []string{"--url"})
	url, ok := flags["--url"]
	if !ok {
		fmt.Println("-url flag is required")
		os.Exit(3)
	}
	fmt.Print("Enter password: ")
	var pass string
	fmt.Scanln(&pass)

	body, _ := json.Marshal(map[string]string{"password": pass})
	req, _ := http.NewRequest("POST", url+"/api/v1/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := doRequest(req)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
}
