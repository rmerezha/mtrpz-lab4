package main

import (
	"fmt"
	"net/http"
	"os"
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
