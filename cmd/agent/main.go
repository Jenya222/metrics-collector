package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second  // Частота обновления метрик
	reportInterval = 10 * time.Second // Частота отправки метрик
)

func main() {
	lastPoll := time.Now()
	lastReport := time.Now()
	metricsGauge := make(map[string]uint64)
	metricsCounter := make(map[string]uint64)

	for {
		now := time.Now()

		if now.Sub(lastPoll) >= pollInterval {
			updateMetrics(metricsGauge, metricsCounter)
			lastPoll = now
		}

		if now.Sub(lastReport) >= reportInterval {
			sendMetrics("counter", metricsCounter)
			sendMetrics("gauge", metricsGauge)
			lastReport = now
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func updateMetrics(metricsGauge map[string]uint64, metricsCounter map[string]uint64) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	v := reflect.ValueOf(m).Elem()
	t := reflect.TypeOf(*m)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue, ok := v.Field(i).Interface().(uint64)
		if ok {
			metricsGauge[fieldName] = fieldValue
		}
	}
	metricsCounter["PollCount"]++
	metricsGauge["RandomValue"] = rand.Uint64()
}

func sendMetrics(metricsType string, metrics map[string]uint64) {
	for name, value := range metrics {
		requestURL := fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", metricsType, name, value)
		req, err := http.NewRequest(http.MethodPost, requestURL, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			continue
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			continue
		}
		res.Body.Close()
		fmt.Printf("Sent metric: %s=%d (status: %d)\n", name, value, res.StatusCode)
	}
}
