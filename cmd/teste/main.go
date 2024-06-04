package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	targetURL   = "http://localhost:8080/"
	requestRate = 5
	apiKey      = "ABC"
)

func main() {
	// First scenario: block after 5 requests per second, no API-KEY
	firstScenario()

	// Wait for the rate limiter to reset
	fmt.Println("Waiting for the rate limiter to reset...")
	time.Sleep(10 * time.Second)

	// Second scenario: 10 requests per second, with RANDOM API-KEY
	secondScenario()

	// Wait for the rate limiter to reset
	fmt.Println("Waiting for the rate limiter to reset...")
	time.Sleep(5 * time.Second)

	// Third scenario: 20 requests per second, with API-KEY "ABC"
	thirdScenario()
}

func firstScenario() {
	fmt.Println("First scenario: block after 5 requests per second, no API-KEY")
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < requestRate; i++ {
		wg.Add(1)
		go sendRequests(&wg, done, "")
	}

	time.Sleep(2 * time.Second)
	close(done)
	wg.Wait()
	fmt.Println("=== First scenario done ===")
}

func secondScenario() {
	fmt.Println("Second scenario: block after 10 requests per second, random API-KEY")
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < requestRate; i++ {
		wg.Add(1)
		go sendRequests(&wg, done, "XYZ")
	}

	time.Sleep(2 * time.Second)
	close(done)
	wg.Wait()
	fmt.Println("=== Second scenario done ===")
}

func thirdScenario() {
	fmt.Println("Third scenario: block after 20 requests per second, with API-KEY \"ABC\"")
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < requestRate; i++ {
		wg.Add(1)
		go sendRequests(&wg, done, apiKey)
	}

	time.Sleep(2 * time.Second)
	close(done)
	wg.Wait()
	fmt.Println("=== Third scenario done ===")
}

func sendRequests(wg *sync.WaitGroup, done <-chan struct{}, apiKey string) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second / time.Duration(requestRate))
	defer ticker.Stop()

	client := http.Client{}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			req, err := http.NewRequest("GET", targetURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}

			// Add the optional API_KEY header if provided
			if apiKey != "" {
				req.Header.Set("API_KEY", apiKey)
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Status: %s\n", resp.Status)
				resp.Body.Close()
			}
		}
	}
}
