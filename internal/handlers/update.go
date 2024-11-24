package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

type UpdateHandler struct{}

type MemStorage struct {
	metrics `json:"metrics"`
}

type metricTypes struct {
	Gauge   float64 `json:"gauge"`
	Counter int64   `json:"counter"`
}

type metrics map[string]metricTypes

var cacheMetrics *MemStorage

func NewUpdateHandler() *UpdateHandler {
	return &UpdateHandler{}
}

func (h UpdateHandler) getMemStorage() *MemStorage {
	if cacheMetrics == nil {
		m := make(metrics)
		cacheMetrics = &MemStorage{m}
	}
	return cacheMetrics
}

// POST counter/someMetric/527
func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 3 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	metricsType, metricsName := path[0], path[1]
	metricsValue, err := strconv.Atoi(path[2])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	memStorage := h.getMemStorage()
	switch metricsType {
	case "counter":
		if v, ok := memStorage.metrics[metricsName]; ok {
			v.Counter += int64(metricsValue)
			memStorage.metrics[metricsName] = v
		} else {
			memStorage.metrics[metricsName] = metricTypes{Counter: int64(metricsValue), Gauge: 0}
		}
	case "gauge":
		if v, ok := memStorage.metrics[metricsName]; ok {
			v.Gauge = float64(metricsValue)
			memStorage.metrics[metricsName] = v
		} else {
			memStorage.metrics[metricsName] = metricTypes{Counter: 0, Gauge: float64(metricsValue)}
		}
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
