package main

import (
	"net/http"
	"sync"
	"time"
)

type HealthCheckResult struct {
	URL          string
	Healthy      bool
	StatusCode   int
	ResponseTime time.Duration
	Error        string
}

// func CheckEndpoint(url string, timeout time.Duration) HealthCheckResult {
// 	client := http.Client{
// 		Timeout: timeout,
// 	}

// 	start := time.Now()
// 	resp, err := client.Get(url)
// 	duration := time.Since(start)

// 	result := HealthCheckResult{
// 		URL:          url,
// 		ResponseTime: duration,
// 	}
// 	if err != nil {
// 		result.Healthy = false
// 		result.Error = err.Error()
// 		return result
// 	}
// 	defer resp.Body.Close()

// 	result.StatusCode = resp.StatusCode
// 	result.Healthy = resp.StatusCode >= 200 && resp.StatusCode < 400

// 	return result
// }


func CheckEndpoint(url string, timeout time.Duration) HealthCheckResult {
	client := http.Client{
		Timeout: timeout,
	}

	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	result := HealthCheckResult{
		URL:          url,
		ResponseTime: duration,
	}
	if err != nil {
		result.Healthy = false
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Healthy = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}


func ConcurrentCheckEndpoints(configs []HealthCheckConfig) [] HealthCheckResult{
	var wg sync.WaitGroup
	results := make([]HealthCheckResult, len(configs))
	ch:= make(chan struct{
		Index int 
		Result HealthCheckResult
	})

	for i,cfg := range configs {
		wg.Add(1)
		go func(index int , config HealthCheckConfig){
			defer wg.Done()
			res:= CheckEndpoint(config)
			ch<- struct{
				Index int; 
				Result HealthCheckResult
			}{Index:index , Result:res}
		}(i,cfg)
			
	}

	go func(){
		wg.Wait()
		close(ch)
	}

	for item := range ch {
        results[item.Index] = item.Result
    }

    return results
}