package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	env := LoadEnv()

	url := env.URL
	numRequests := env.NUM_REQUESTS

	var wg sync.WaitGroup
	var mu sync.Mutex
	var responseTimes []time.Duration
	var successCount, clientErrorCount, serverErrorCount int

	client := &http.Client{
		Timeout: time.Duration(env.TIMEOUT) * time.Second,
	}

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			elapsedTime, statusCode := sendRequest(client, url, i)
			mu.Lock()
			responseTimes = append(responseTimes, elapsedTime)
			if statusCode >= 200 && statusCode < 400 {
				successCount++
			} else if statusCode >= 400 && statusCode < 500 {
				clientErrorCount++
			} else if statusCode >= 500 {
				serverErrorCount++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	fmt.Printf("Успешные запросы (200-399): %d\n", successCount)
	fmt.Printf("Ошибки клиента (400-499): %d\n", clientErrorCount)
	fmt.Printf("Ошибки сервера (500 и выше): %d\n", serverErrorCount)

	var minTime, maxTime time.Duration
	var totalTime time.Duration

	if len(responseTimes) > 0 {
		minTime = responseTimes[0]
		maxTime = responseTimes[0]
		for _, t := range responseTimes {
			totalTime += t
			if t < minTime {
				minTime = t
			}
			if t > maxTime {
				maxTime = t
			}
		}
		avgTime := totalTime / time.Duration(len(responseTimes))
		fmt.Printf("Минимальное время ответа: %s\n", minTime)
		fmt.Printf("Максимальное время ответа: %s\n", maxTime)
		fmt.Printf("Среднее время ответа: %s\n", avgTime)
	}
}

func sendRequest(client *http.Client, url string, requestID int) (time.Duration, int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Ошибка при создании запроса %d: %s\n", requestID, err)
		return 0, 0
	}

	req.Header.Add("Authorization", "Basic YWRtaW46S2VmNDVvbGRA")

	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при отправке запроса %d: %s\n", requestID, err)
		return 0, 0
	}
	defer resp.Body.Close()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Запрос %d завершен с кодом ответа: %d, время ответа: %s\n", requestID, resp.StatusCode, elapsedTime)
	return elapsedTime, resp.StatusCode
}
