package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Jenya222/metrics-collector/internal/storage"
)

type UpdateHandler struct {
	storage storage.MemStorage
}

func NewUpdateHandler() *UpdateHandler {
	return &UpdateHandler{}
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
	memStorage := storage.GetMemStorage()
	switch metricsType {
	case "counter":
		memStorage.StoreCounter(metricsName, int64(metricsValue))
	case "gauge":
		memStorage.StoreGauge(metricsName, float64(metricsValue))
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
