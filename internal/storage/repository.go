package storage

type MemStorage struct {
	metrics `json:"metrics"`
}

type metricTypes struct {
	Gauge   float64 `json:"gauge"`
	Counter int64   `json:"counter"`
}

type metrics map[string]metricTypes

var cacheMetrics *MemStorage

func GetMemStorage() *MemStorage {
	if cacheMetrics == nil {
		mt := make(metrics)
		cacheMetrics = &MemStorage{mt}
	}
	return cacheMetrics
}

func (m MemStorage) StoreCounter(metricsName string, value int64) {
	if v, ok := m.metrics[metricsName]; ok {
		v.Counter += value
		m.metrics[metricsName] = v
	} else {
		m.metrics[metricsName] = metricTypes{Counter: value, Gauge: 0}
	}
}

func (m MemStorage) StoreGauge(metricsName string, value float64) {
	if v, ok := m.metrics[metricsName]; ok {
		v.Gauge = value
		m.metrics[metricsName] = v
	} else {
		m.metrics[metricsName] = metricTypes{Counter: 0, Gauge: value}
	}
}
