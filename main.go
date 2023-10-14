package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.bing.com/?mkt=zh-CN")
	if err != nil {
		fmt.Printf("Fetch url error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %v\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body error: %v\n", err)
		return
	}

	fmt.Printf("%s\n", body)
}
